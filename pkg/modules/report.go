package modules

import (
	"fmt"
	"strings"
	"time"

	"github.com/pajbot/pajbot2/pkg"
	mbase "github.com/pajbot/pajbot2/pkg/modules/base"
	"github.com/pajbot/pajbot2/pkg/report"
	"github.com/pajbot/pajbot2/pkg/twitchactions"
	"github.com/pajbot/utils"
)

func init() {
	Register("report", func() pkg.ModuleSpec {
		return &Spec{
			id:    "report",
			name:  "Report",
			maker: newReport,
		}
	})
}

type Report struct {
	mbase.Base

	reportHolder *report.Holder
}

var _ pkg.Module = &Report{}

func newReport(b *mbase.Base) pkg.Module {
	m := &Report{
		Base: *b,

		reportHolder: _server.reportHolder,
	}

	return m
}

func (m *Report) ProcessReport(user pkg.User, parts []string) pkg.Actions {
	duration := 600

	if parts[0] == "!report" {
	} else if parts[0] == "!longreport" {
		duration = 28800
	} else {
		return nil
	}

	if !user.HasPermission(m.BotChannel().Channel(), pkg.PermissionReport) {
		return twitchactions.DoWhisper(user, "You don't have permissions to use the !report command")
	}

	var reportedUsername string
	var reason string

	reportedUsername = strings.ToLower(utils.FilterUsername(parts[1]))

	if reportedUsername == user.GetName() {
		return twitchactions.DoWhisper(user, "You can't report yourself")
	}

	if len(parts) >= 3 {
		reason = strings.Join(parts[2:], " ")
	}

	return m.report(m.BotChannel().Bot(), user, m.BotChannel().Channel(), reportedUsername, reason, duration)
}

func (m *Report) report(bot pkg.Sender, reporter pkg.User, targetChannel pkg.Channel, targetUsername string, reason string, duration int) pkg.Actions {
	// s := fmt.Sprintf("%s reported %s in #%s (%s) - https://api.gempir.com/channel/forsen/user/%s", reporter.GetName(), targetUsername, targetChannel.GetName(), reason, targetUsername)

	r := report.Report{
		Channel: report.ReportUser{
			ID:   targetChannel.GetID(),
			Name: targetChannel.GetName(),
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
	r.Logs = bot.GetUserContext().GetContext(r.Channel.ID, r.Target.ID)

	oldReport, inserted, _ := m.reportHolder.Register(r)

	if !inserted {
		// Report for this user in this channel already exists

		reportAge := time.Since(oldReport.Time)
		if reportAge < time.Minute*10 {
			// User was reported less than 10 minutes ago, don't let this user be timed out again
			fmt.Printf("Skipping timeout because user was timed out too shortly ago: %s\n", reportAge)
			return twitchactions.DoWhisper(reporter, "User successfully reported, but the last report was less than ten minutes ago so the timeout is skipped")
		}

		fmt.Println("Update report")
		r.ID = oldReport.ID
		m.reportHolder.Update(r)
	}

	actions := &twitchactions.Actions{}

	actions.Whisper(reporter, fmt.Sprintf("Successfully reported user %s", targetUsername))
	actions.Timeout(bot.MakeUser(targetUsername), time.Duration(duration)*time.Second)

	return actions
}

func (m *Report) OnWhisper(event pkg.MessageEvent) pkg.Actions {
	const usageString = `Usage: #channel !report username (reason) i.e. #forsen !report Karl_Kons spamming stuff`

	user := event.User
	message := event.Message

	parts := strings.Split(message.GetText(), " ")
	if len(parts) < 2 {
		return nil
	}

	return m.ProcessReport(user, parts)
}

func (m *Report) OnMessage(event pkg.MessageEvent) pkg.Actions {
	user := event.User
	message := event.Message

	parts := strings.Split(message.GetText(), " ")
	if len(parts) < 2 {
		return nil
	}

	return m.ProcessReport(user, parts)
}
