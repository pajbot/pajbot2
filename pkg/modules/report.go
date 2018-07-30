package modules

import (
	"log"
	"strings"

	"github.com/pajlada/pajbot2/pkg"
)

type Report struct {
	server *server
}

func NewReportModule() *Report {
	return &Report{
		server: &_server,
	}
}

func (m Report) Name() string {
	return "Report"
}

func (m *Report) Register() error {
	return nil
}

func (m *Report) OnWhisper(bot pkg.Sender, source pkg.User, message pkg.Message) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) < 2 {
		return nil
	}

	if parts[0] != "!report" {
		return nil
	}

	// XXX: Channel-specific permissions
	if !source.HasPermission(pkg.PermissionReport) {
		return nil
	}

	var reportedUsername string
	var reportedChannel string
	var reason string

	log.Printf("lol: %#v", parts)

	reportedUsername = strings.ToLower(parts[1])
	if len(parts) >= 3 {
		reportedChannel = strings.ToLower(strings.TrimPrefix(parts[2], "#"))
	} else {
		// Reply with proper whisper usage
		return nil
	}

	if len(parts) >= 4 {
		reason = strings.Join(parts[3:], " ")
	}

	log.Printf("%s reported %s in #%s (%s)", source.GetName(), reportedUsername, reportedChannel, reason)

	return nil

}

func (m *Report) OnMessage(bot pkg.Sender, source pkg.Channel, user pkg.User, message pkg.Message, action pkg.Action) error {
	parts := strings.Split(message.GetText(), " ")
	if len(parts) < 2 {
		return nil
	}

	if parts[0] != "!report" {
		return nil
	}

	if !user.HasPermission(pkg.PermissionReport) {
		return nil
	}

	var reportedUsername string
	reportedChannel := source.GetChannel()
	var reason string

	reportedUsername = strings.ToLower(parts[1])

	if len(parts) >= 3 {
		reason = strings.Join(parts[2:], " ")
	}

	log.Printf("%s reported %s in #%s (%s)", user.GetName(), reportedUsername, reportedChannel, reason)

	return nil
}
