package report

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/pubsub"
	"github.com/pajlada/pajbot2/pkg/users"
)

type ReportUser struct {
	ID   string
	Name string
	Type string `json:",omitempty"`
}

type Report struct {
	ID       uint32
	Channel  ReportUser
	Reporter ReportUser
	Target   ReportUser
	Reason   string `json:",omitempty"`
	Logs     []string
}

type Holder struct {
	db     *sql.DB
	pubSub *pubsub.PubSub

	reportsMutex *sync.Mutex
	reports      map[uint32]Report
}

var _ pubsub.Connection = &Holder{}
var _ pubsub.SubscriptionHandler = &Holder{}

func New(db *sql.DB, pubSub *pubsub.PubSub) (*Holder, error) {
	h := &Holder{
		db:     db,
		pubSub: pubSub,

		reportsMutex: &sync.Mutex{},
		reports:      make(map[uint32]Report),
	}

	err := h.Load()
	if err != nil {
		return nil, err
	}

	pubSub.Subscribe(h, "HandleReport", nil)
	pubSub.HandleSubscribe(h, "ReportReceived")

	return h, nil
}

func (h *Holder) Load() error {
	rows, err := h.db.Query("SELECT `id`, `channel_id`, `channel_name`, `channel_type`, `reporter_id`, `reporter_name`, `target_id`, `target_name`, `reason`, `logs` FROM `Report`")
	if err != nil {
		return err
	}
	defer rows.Close()

	h.reportsMutex.Lock()
	defer h.reportsMutex.Unlock()

	for rows.Next() {
		var report Report
		var logsString string

		if err := rows.Scan(&report.ID, &report.Channel.ID, &report.Channel.Name, &report.Channel.Type, &report.Reporter.ID, &report.Reporter.Name, &report.Target.ID, &report.Target.Name, &report.Reason, &logsString); err != nil {
			return err
		}

		report.Logs = strings.Split(logsString, "\n")

		h.reports[report.ID] = report
	}

	return nil
}

type handleReportMessage struct {
	Action    string
	ChannelID string
	ReportID  uint32
}

func (h *Holder) Register(report Report) error {
	const queryF = `
	INSERT INTO Report
		(channel_id, channel_name, channel_type,
		reporter_id, reporter_name, target_id, target_name, reason, logs)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	res, err := h.db.Exec(queryF, report.Channel.ID, report.Channel.Name, report.Channel.Type, report.Reporter.ID, report.Reporter.Name, report.Target.ID, report.Target.Name, report.Reason, strings.Join(report.Logs, "\n"))
	if err != nil {
		fmt.Printf("Error inserting report %v into SQL: %s\n", report, err)
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error getting last insert id forsenT: %s\n", err)
		return err
	}

	report.ID = uint32(id)

	h.pubSub.Publish("ReportReceived", report, pkg.PubSubAdminAuth())

	{
		h.reportsMutex.Lock()
		defer h.reportsMutex.Unlock()

		h.reports[report.ID] = report
	}

	return nil
}

type reportHandled struct {
	ReportID uint32
	Handler  ReportUser
	Action   string
}

func (h *Holder) handleReport(reportID uint32, action string) error {
	h.reportsMutex.Lock()
	defer h.reportsMutex.Unlock()

	if report, ok := h.reports[reportID]; ok {
		// Remove from SQL
		const queryF = "DELETE FROM Report WHERE `id`=?"

		res, err := h.db.Exec(queryF, report.ID)
		if err != nil {
			return err
		}

		count, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if count != 1 {
			return errors.New("Report existed in internal storage, but not in SQL")
		}

		// TODO: Insert into new table: HandledReport

		msg := reportHandled{
			ReportID: report.ID,
			Handler: ReportUser{
				ID:   "11148817",
				Name: "pajlada",
			},
			Action: action,
		}

		h.pubSub.Publish("ReportHandled", &msg, pkg.PubSubAdminAuth())

		switch action {
		case "ban":
			h.pubSub.Publish("Ban", &pkg.PubSubBan{
				Channel: report.Channel.Name,
				Target:  report.Target.Name,
				Reason:  report.Reason,
			}, pkg.PubSubAdminAuth())

		case "undo":
			h.pubSub.Publish("Untimeout", &pkg.PubSubUntimeout{
				Channel: report.Channel.Name,
				Target:  report.Target.Name,
			}, pkg.PubSubAdminAuth())
		}

		// Delete from our internal storage
		delete(h.reports, report.ID)

	} else {
		fmt.Printf("No report with the id %d found\n", reportID)
	}

	return nil
}

func (h *Holder) MessageReceived(topic string, data []byte, auth *pkg.PubSubAuthorization) error {
	var msg handleReportMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		fmt.Println("Error unmarshalling:", err)
		return nil
	}

	fmt.Printf("Handle report: %#v\n", msg)

	return h.handleReport(msg.ReportID, msg.Action)
}

func (h *Holder) ConnectionSubscribed(connection pubsub.Connection, topic string, authorization *pkg.PubSubAuthorization) (error, bool) {
	switch topic {
	case "ReportReceived":
		if authorization == nil {
			return nil, false
		}

		// Verify authorization
		const queryF = `
SELECT twitch_username FROM User
	WHERE twitch_userid=? AND twitch_nonce=? LIMIT 1;
`

		rows, err := h.db.Query(queryF, authorization.TwitchUserID, authorization.Nonce)
		if err != nil {
			return err, true
		}
		defer rows.Close()

		if !rows.Next() {
			return nil, false
		}

		hasPermission, err := users.HasGlobalPermission(authorization.TwitchUserID, pkg.PermissionModeration)
		if err != nil {
			return err, false
		}

		if !hasPermission {
			return nil, false
		}

		for _, report := range h.reports {
			bytes, err := json.Marshal(report)
			if err != nil {
				return err, true
			}
			connection.MessageReceived(topic, bytes, pkg.PubSubAdminAuth())
		}
	}

	return nil, true
}
