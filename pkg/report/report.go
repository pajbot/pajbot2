package report

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

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
	Time     time.Time
}

type Holder struct {
	db        *sql.DB
	pubSub    *pubsub.PubSub
	userStore pkg.UserStore

	reportsMutex *sync.Mutex
	reports      map[uint32]Report
}

var _ pubsub.Connection = &Holder{}
var _ pubsub.SubscriptionHandler = &Holder{}

func New(db *sql.DB, pubSub *pubsub.PubSub, userStore pkg.UserStore) (*Holder, error) {
	h := &Holder{
		db:        db,
		pubSub:    pubSub,
		userStore: userStore,

		reportsMutex: &sync.Mutex{},
		reports:      make(map[uint32]Report),
	}

	err := h.Load()
	if err != nil {
		return nil, err
	}

	pubSub.Subscribe(h, "HandleReport", nil)
	pubSub.Subscribe(h, "TimeoutEvent", nil)
	pubSub.Subscribe(h, "BanEvent", nil)
	pubSub.HandleSubscribe(h, "ReportReceived")

	return h, nil
}

func (h *Holder) Load() error {
	rows, err := h.db.Query("SELECT `id`, `channel_id`, `channel_name`, `channel_type`, `reporter_id`, `reporter_name`, `target_id`, `target_name`, `reason`, `logs`, `time` FROM `Report`")
	if err != nil {
		return err
	}
	defer rows.Close()

	h.reportsMutex.Lock()
	defer h.reportsMutex.Unlock()

	for rows.Next() {
		var report Report
		var logsString string

		if err := rows.Scan(&report.ID, &report.Channel.ID, &report.Channel.Name, &report.Channel.Type, &report.Reporter.ID, &report.Reporter.Name, &report.Target.ID, &report.Target.Name, &report.Reason, &logsString, &report.Time); err != nil {
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
	Duration  *uint32
}

func (h *Holder) Register(report Report) (*Report, bool, error) {
	const queryF = `
	INSERT INTO Report
		(channel_id, channel_name, channel_type,
		reporter_id, reporter_name, target_id, target_name, reason, logs, time)
	VALUES
		(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	h.reportsMutex.Lock()
	defer h.reportsMutex.Unlock()

	// Don't accept reports for users that have already been reported
	for _, oldReport := range h.reports {
		if oldReport.Channel.ID == report.Channel.ID && oldReport.Target.ID == report.Target.ID {
			fmt.Println("Report already registered for this target in this channel")
			return &oldReport, false, nil
		}
	}

	res, err := h.db.Exec(queryF, report.Channel.ID, report.Channel.Name, report.Channel.Type, report.Reporter.ID, report.Reporter.Name, report.Target.ID, report.Target.Name, report.Reason, strings.Join(report.Logs, "\n"), report.Time)
	if err != nil {
		fmt.Printf("Error inserting report %v into SQL: %s\n", report, err)
		return nil, false, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		fmt.Printf("Error getting last insert id: %s\n", err)
		return nil, false, err
	}

	report.ID = uint32(id)

	h.pubSub.Publish("ReportReceived", report, pkg.PubSubAdminAuth())

	h.reports[report.ID] = report

	return &report, true, nil
}

func (h *Holder) Update(report Report) error {
	if report.ID == 0 {
		return errors.New("Missing report ID in Update")
	}

	const queryF = `UPDATE Report SET time=?, logs=? WHERE id=?`
	_, err := h.db.Exec(queryF, report.Time, strings.Join(report.Logs, "\n"), report.ID)
	if err != nil {
		return err
	}

	h.reportsMutex.Lock()
	defer h.reportsMutex.Unlock()
	h.reports[report.ID] = report

	// TODO: Send some "ReportUpdated" message

	return nil
}

type reportHandled struct {
	ReportID uint32
	Handler  ReportUser
	Action   string
}

func (h *Holder) handleReport(action handleReportMessage, auth *pkg.PubSubAuthorization) error {
	h.reportsMutex.Lock()
	defer h.reportsMutex.Unlock()

	if auth == nil {
		fmt.Println("Missing auth in HandleReport")
		return nil
	}

	report, ok := h.reports[action.ReportID]
	if !ok {
		fmt.Printf("No report found with ID %d\n", action.ReportID)
		// No report found with this ID
		return nil
	}

	// Remove report from SQL and our local map
	err := h.dismissReport(report.ID)
	if err != nil {
		fmt.Println("Error dismissing report", err)
		return err
	}

	// TODO: Insert into new table: HandledReport

	msg := reportHandled{
		ReportID: report.ID,
		Handler: ReportUser{
			ID:   auth.TwitchUserID,
			Name: h.userStore.GetName(auth.TwitchUserID),
		},
		Action: action.Action,
	}

	h.pubSub.Publish("ReportHandled", &msg, pkg.PubSubAdminAuth())

	switch action.Action {
	case "ban":
		h.pubSub.Publish("Ban", &pkg.PubSubBan{
			Channel: report.Channel.Name,
			Target:  report.Target.Name,
			Reason:  report.Reason,
		}, pkg.PubSubAdminAuth())

	case "timeout":
		var duration uint32
		duration = 600
		if action.Duration != nil {
			duration = *action.Duration
		}
		h.pubSub.Publish("Timeout", &pkg.PubSubTimeout{
			Channel:  report.Channel.Name,
			Target:   report.Target.Name,
			Duration: duration,
			Reason:   report.Reason,
		}, pkg.PubSubAdminAuth())

	case "undo":
		h.pubSub.Publish("Untimeout", &pkg.PubSubUntimeout{
			Channel: report.Channel.Name,
			Target:  report.Target.Name,
		}, pkg.PubSubAdminAuth())
	default:
		fmt.Println("Unhandled action", action.Action)
	}

	return nil
}

// dismissReport assumes that reportsMutex has already been locked
func (h *Holder) dismissReport(reportID uint32) error {
	// Delete from SQL
	const queryF = "DELETE FROM Report WHERE `id`=?"

	_, err := h.db.Exec(queryF, reportID)
	if err != nil {
		return err
	}

	// Delete from our internal storage
	delete(h.reports, reportID)

	return nil
}

func (h *Holder) handleBanEvent(banEvent pkg.PubSubBanEvent) error {
	h.reportsMutex.Lock()
	defer h.reportsMutex.Unlock()

	for reportID, report := range h.reports {
		if report.Channel.ID == banEvent.Channel.ID && report.Target.ID == banEvent.Target.ID {
			// Found matching report
			h.dismissReport(reportID)

			break
		}
	}

	return nil
}

func (h *Holder) MessageReceived(topic string, data []byte, auth *pkg.PubSubAuthorization) error {
	switch topic {
	case "HandleReport":
		var msg handleReportMessage
		err := json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling:", err)
			return nil
		}

		fmt.Printf("Handle report: %#v\n", msg)

		return h.handleReport(msg, auth)

	case "BanEvent":
		var msg pkg.PubSubBanEvent
		err := json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling:", err)
			return nil
		}

		return h.handleBanEvent(msg)
	}

	return nil
}

func (h *Holder) ConnectionSubscribed(connection pubsub.Connection, topic string, auth *pkg.PubSubAuthorization) (error, bool) {
	switch topic {
	case "ReportReceived":
		if auth == nil {
			return nil, false
		}

		// Verify authorization
		const queryF = `
SELECT twitch_username FROM User
	WHERE twitch_userid=? AND twitch_nonce=? LIMIT 1;
`

		rows, err := h.db.Query(queryF, auth.TwitchUserID, auth.Nonce)
		if err != nil {
			fmt.Println(err)
			return err, true
		}
		defer rows.Close()

		if !rows.Next() {
			return nil, false
		}

		hasPermission, err := users.HasGlobalPermission(auth.TwitchUserID, pkg.PermissionModeration)
		if err != nil {
			fmt.Println(err)
			return err, false
		}

		if !hasPermission {
			return nil, false
		}

		fmt.Println("Send reports to new connection")

		for _, report := range h.reports {
			bytes, err := json.Marshal(report)
			if err != nil {
				fmt.Println(err)
				return err, true
			}
			connection.MessageReceived(topic, bytes, pkg.PubSubAdminAuth())
		}
	}

	return nil, true
}
