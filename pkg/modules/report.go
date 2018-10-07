package modules

import (
	"fmt"
	"strings"
	"time"

	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/report"
)

type Report struct {
	server       *server
	reportHolder *report.Holder
}

var _ pkg.Module = &Report{}

func NewReportModule(reportHolder *report.Holder) *Report {
	return &Report{
		server:       &_server,
		reportHolder: reportHolder,
	}
}

func (m Report) Name() string {
	return "Report"
}

func (m *Report) Register() error {
	return nil
}

func (m *Report) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	const usageString = `Usage: !report username channel (reason) i.e. !report Karl_Kons forsen spamming stuff`

	parts := strings.Split(message.GetText(), " ")
	if len(parts) < 2 {
		return nil
	}

	duration := 600

	if parts[0] == "!report" {
	} else if parts[0] == "!longreport" {
		duration = 28800
	} else {
		return nil
	}

	var reportedUsername string
	var reportedChannel string
	var reason string

	reportedUsername = strings.ToLower(parts[1])
	if reportedUsername == source.GetName() {
		return nil
	}

	if len(parts) >= 3 {
		reportedChannel = strings.ToLower(strings.TrimPrefix(parts[2], "#"))
	} else {
		bot.Whisper(source, usageString)
		return nil
	}

	channel := bot.MakeChannel(reportedChannel)
	if !source.HasPermission(channel, pkg.PermissionReport) {
		bot.Whisper(source, "you don't have permissions to use the !report command")
		return nil
	}

	if len(parts) >= 4 {
		reason = strings.Join(parts[3:], " ")
	}

	m.report(bot, source, channel, reportedUsername, reason, duration)

	return nil
}

func (m *Report) report(bot pkg.Sender, reporter pkg.User, targetChannel pkg.Channel, targetUsername string, reason string, duration int) {
	// s := fmt.Sprintf("%s reported %s in #%s (%s) - https://api.gempir.com/channel/forsen/user/%s", reporter.GetName(), targetUsername, targetChannel.GetChannel(), reason, targetUsername)

	r := report.Report{
		Channel: report.ReportUser{
			ID:   targetChannel.GetID(),
			Name: targetChannel.GetChannel(),
			Type: "twitch",
		},
		Reporter: report.ReportUser{
			ID:   reporter.GetID(),
			Name: reporter.GetName(),
		},
		Target: report.ReportUser{
			ID:   bot.GetUserStore().GetID(targetUsername),
			Name: targetUsername,
		},
		Reason: reason,
		Time:   time.Now(),
	}

	oldReport, inserted, _ := m.reportHolder.Register(r)

	if !inserted {
		// Report for this user in this channel already exists

		if time.Now().Sub(oldReport.Time) < time.Minute*10 {
			// User was reported less than 10 minutes ago, don't let this user be timed out again
			fmt.Printf("Skipping timeout because user was timed out too shortly ago: %s\n", time.Now().Sub(oldReport.Time))
			return
		}

		fmt.Println("Update report")
		r.ID = oldReport.ID
		m.reportHolder.Update(r)
	}

	bot.Timeout(targetChannel, bot.MakeUser(targetUsername), duration, "")
}

func (m *Report) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) < 2 {
		return nil
	}

	duration := 600

	if parts[0] == "!report" {
	} else if parts[0] == "!longreport" {
		duration = 28800
	} else {
		return nil
	}

	if !user.HasPermission(source, pkg.PermissionReport) {
		return nil
	}

	var reportedUsername string
	var reason string

	reportedUsername = strings.ToLower(parts[1])

	if reportedUsername == user.GetName() {
		return nil
	}

	if len(parts) >= 3 {
		reason = strings.Join(parts[2:], " ")
	}

	m.report(bot, user, source, reportedUsername, reason, duration)

	return nil
}
