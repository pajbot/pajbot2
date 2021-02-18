package report

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pajbot/pajbot2/pkg"
)

type ReportUser struct {
	ID   string
	Name string
	Type string `json:",omitempty"`
}

func (u ReportUser) GetName() string {
	return u.Name
}

func (u ReportUser) GetID() string {
	return u.ID
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

type HistoricReport struct {
	ID             uint32
	Channel        ReportUser
	Reporter       ReportUser
	Target         ReportUser
	Reason         string
	Logs           []string
	Time           time.Time
	Handler        ReportUser
	Action         uint8
	ActionDuration uint32
	TimeHandled    time.Time
}

type Holder struct {
	db           *sql.DB
	pubSub       pkg.PubSub
	userStore    pkg.UserStore
	channelStore pkg.ChannelStore
	botStore     pkg.BotStore

	reportsMutex *sync.Mutex
	reports      map[uint32]Report
}

var _ pkg.PubSubConnection = &Holder{}
var _ pkg.PubSubSource = &Holder{}
var _ pkg.PubSubSubscriptionHandler = &Holder{}

func New(app pkg.Application) (*Holder, error) {
	h := &Holder{
		db:           app.SQL(),
		pubSub:       app.PubSub(),
		userStore:    app.UserStore(),
		botStore:     app.TwitchBots(),
		channelStore: app.ChannelStore(),

		reportsMutex: &sync.Mutex{},
		reports:      make(map[uint32]Report),
	}

	err := h.Load()
	if err != nil {
		return nil, err
	}

	h.pubSub.Subscribe(h, "HandleReport", nil)
	h.pubSub.Subscribe(h, "TimeoutEvent", nil)
	h.pubSub.Subscribe(h, "BanEvent", nil)
	h.pubSub.HandleSubscribe(h, "ReportReceived")

	return h, nil
}

func (h *Holder) AuthenticatedUser() pkg.User {
	return nil
}

func (h *Holder) IsApplication() bool {
	return true
}

func (h *Holder) Connection() pkg.PubSubConnection {
	return h
}

func (h *Holder) Load() error {
	const queryF = `
SELECT
	id, channel_id, channel_name,
	channel_type, reporter_id, reporter_name,
	target_id, target_name,
	reason, logs, time
FROM
	report`

	rows, err := h.db.Query(queryF) // GOOD
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
	Action    uint8
	ChannelID string
	ReportID  uint32
	Duration  *uint32
	Handler   ReportUser
}

func (h *Holder) Register(report Report) (*Report, bool, error) {
	const queryF = `
	INSERT INTO report
		(channel_id, channel_name, channel_type,
		reporter_id, reporter_name, target_id, target_name, reason, logs, time)
	VALUES
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
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

	row := h.db.QueryRow(queryF, report.Channel.ID, report.Channel.Name, report.Channel.Type, report.Reporter.ID, report.Reporter.Name, report.Target.ID, report.Target.Name, report.Reason, strings.Join(report.Logs, "\n"), report.Time) // GOOD
	var id int64
	err := row.Scan(&id)
	if err != nil {
		fmt.Printf("Error inserting report %v into SQL: %s\n", report, err)
		return nil, false, err
	}

	report.ID = uint32(id)

	h.pubSub.Publish(h, "ReportReceived", report)

	h.reports[report.ID] = report

	return &report, true, nil
}

func (h *Holder) Update(report Report) error {
	if report.ID == 0 {
		return errors.New("missing report ID in Update")
	}

	const queryF = `UPDATE report SET time=$1, logs=$2 WHERE id=$3`
	_, err := h.db.Exec(queryF, report.Time, strings.Join(report.Logs, "\n"), report.ID) // GOOD
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
	Action   uint8
}

func (h *Holder) insertHistoricReport(report Report, action handleReportMessage) {
	const queryF = `
INSERT INTO
	report_history
(
channel_id, channel_name, channel_type,
reporter_id, reporter_name,
target_id, target_name,
reason, logs,
time,
handler_id, handler_name,
action, action_duration,
time_handled
)

VALUES (
$1,$2,$3,
$4,$5,
$6,$7,
$8,$9,
$10,
$11,$12,
$13,$14,
$15
)`

	var actionDuration uint32
	if action.Duration != nil {
		actionDuration = *action.Duration
	}
	_, err := h.db.Exec(queryF, // GOOD
		report.Channel.ID, report.Channel.Name, report.Channel.Type,
		report.Reporter.ID, report.Reporter.Name,
		report.Target.ID, report.Target.Name,
		report.Reason, strings.Join(report.Logs, "\n"),
		report.Time,
		action.Handler.ID, action.Handler.Name,
		action.Action, actionDuration,
		time.Now())
	if err != nil {
		panic(err)
	}
}

func (h *Holder) handleReport(source pkg.PubSubSource, action handleReportMessage) error {
	h.reportsMutex.Lock()
	defer h.reportsMutex.Unlock()

	user := source.AuthenticatedUser()
	if user == nil {
		fmt.Println("Missing auth in HandleReport")
		return nil
	}

	report, ok := h.reports[action.ReportID]
	if !ok {
		fmt.Printf("No report found with ID %d\n", action.ReportID)
		// No report found with this ID
		return nil
	}

	if !user.HasPermission(report.Channel, pkg.PermissionModeration) {
		fmt.Println("user does not have moderation permission")
		return nil
	}

	// Remove report from SQL and our local map
	err := h.dismissReport(report.ID)
	if err != nil {
		fmt.Println("Error dismissing report", err)
		return err
	}

	action.Handler.Name = user.GetName()
	action.Handler.ID = user.GetID()

	bot := h.botStore.GetBotFromChannel(report.Channel.ID)
	reporterInst := bot.MakeUser(report.Reporter.Name)

	// TODO: Insert into new table: HandledReport
	h.insertHistoricReport(report, action)

	msg := &reportHandled{
		ReportID: report.ID,
		Handler: ReportUser{
			ID:   user.GetID(),
			Name: h.userStore.GetName(user.GetID()),
		},
		Action: action.Action,
	}

	h.pubSub.Publish(h, "ReportHandled", msg)

	switch action.Action {
	case pkg.ReportActionBan:
		h.pubSub.Publish(h, "Ban", &pkg.PubSubBan{
			Channel: report.Channel.Name,
			Target:  report.Target.Name,
			// Reason:  report.Reason,
		})
		bot.Whisper(reporterInst, fmt.Sprintf("Thanks to your report, user %s has been permanently banned :)", report.Target.Name))

	case pkg.ReportActionTimeout:
		var duration uint32
		duration = 600
		if action.Duration != nil {
			duration = *action.Duration
		}
		h.pubSub.Publish(h, "Timeout", &pkg.PubSubTimeout{
			Channel:  report.Channel.Name,
			Target:   report.Target.Name,
			Duration: duration,
			// Reason:   report.Reason,
		})
		bot.Whisperf(reporterInst, "Thanks to your report, user %s has been timed out for %d seconds :)", report.Target.Name, duration)

	case pkg.ReportActionDismiss:
		bot.Whisper(reporterInst, fmt.Sprintf("Your report of %s has been dismissed with no further action taken :\\", report.Target.Name))
		// We don't need to do anything else here, as we've already dismissed the report prior to the ban/timeout/untimeout events being sent out

	case pkg.ReportActionUndo:
		h.pubSub.Publish(h, "Untimeout", &pkg.PubSubUntimeout{
			Channel: report.Channel.Name,
			Target:  report.Target.Name,
		})
		bot.Whisper(reporterInst, fmt.Sprintf("Your report of %s has been undone with no further action taken :\\", report.Target.Name))

	default:
		fmt.Println("Unhandled action", action.Action)
	}

	return nil
}

// dismissReport assumes that reportsMutex has already been locked
func (h *Holder) dismissReport(reportID uint32) error {
	// Delete from SQL
	const queryF = "DELETE FROM report WHERE id=$1"

	_, err := h.db.Exec(queryF, reportID) // GOOD
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

func (h *Holder) MessageReceived(source pkg.PubSubSource, topic string, data []byte) error {
	switch topic {
	case "HandleReport":
		var msg handleReportMessage
		err := json.Unmarshal(data, &msg)
		if err != nil {
			fmt.Println("Error unmarshalling:", err)
			return nil
		}

		fmt.Printf("Handle report: %#v\n", msg)

		return h.handleReport(source, msg)

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

type reportReceivedParameters struct {
	ChannelID string
}

func (h *Holder) ConnectionSubscribed(source pkg.PubSubSource, topic string, parameters json.RawMessage) (bool, error) {
	switch topic {
	case "ReportReceived":
		fmt.Println("aaaaaaaaaaaaaaaa")
		user := source.AuthenticatedUser()
		if user == nil {
			fmt.Println("no user")
			return false, nil
		}

		fmt.Println("Parameters:", string(parameters))

		var parsedParams reportReceivedParameters
		err := json.Unmarshal(parameters, &parsedParams)
		if err != nil {
			fmt.Println("Error parsing subscription parameters:", err)
			return false, nil
		}

		channel := h.channelStore.TwitchChannel(parsedParams.ChannelID)
		if channel == nil {
			fmt.Println("Channel with id", parsedParams.ChannelID, "is not being moderated by us")
			return false, nil
		}

		fmt.Println("Channel ID:", parsedParams.ChannelID)

		fmt.Println("User name:", user.GetName())

		hasPermission := user.HasPermission(channel, pkg.PermissionModeration)

		if !hasPermission {
			fmt.Println("user", user.GetName(), "does not have permission in channel", channel.GetName())
			return false, nil
		}

		fmt.Println("Send reports to new connection")

		for _, report := range h.reports {
			if report.Channel.GetID() != channel.GetID() {
				continue
			}
			bytes, err := json.Marshal(report)
			if err != nil {
				fmt.Println(err)
				return true, err
			}
			source.Connection().MessageReceived(h, topic, bytes)
		}
	}

	return true, nil
}
