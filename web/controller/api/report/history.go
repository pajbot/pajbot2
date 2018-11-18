package report

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/pkg"
	"github.com/pajlada/pajbot2/pkg/report"
	"github.com/pajlada/pajbot2/pkg/users"
	"github.com/pajlada/pajbot2/pkg/utils"
	"github.com/pajlada/pajbot2/web/router"
	"github.com/pajlada/pajbot2/web/state"
)

func Load(parent *mux.Router) {
	m := parent.PathPrefix("/report").Subrouter()

	router.RGet(m, `/history`, apiHistory)
}

func apiHistory(w http.ResponseWriter, r *http.Request) {
	const queryF = `
SELECT
id,
channel_id, channel_name, channel_type,
reporter_id, reporter_name,
target_id, target_name,
reason, logs,
time,
handler_id, handler_name,
action, action_duration,
time_handled

FROM
	ReportHistory
ORDER BY time_handled DESC
LIMIT 50;
	`

	c := state.Context(w, r)

	if c.Session == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return
	}

	user := users.NewSimpleTwitchUser(c.Session.TwitchUserID, c.Session.TwitchUserName)
	if user == nil {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint")
		return
	}

	if !user.HasGlobalPermission(pkg.PermissionReportAPI) {
		utils.WebWriteError(w, 400, "Not authorized to view this endpoint!!!")
		return
	}

	rows, err := c.SQL.Query(queryF)
	if err != nil {
		fmt.Println("error in mysql query apiUser:", err)
		utils.WebWriteError(w, 500, "Internal error")
		return
	}

	defer rows.Close()

	type xd struct {
		Reports []report.HistoricReport
	}

	var response xd

	for rows.Next() {
		var r report.HistoricReport
		var logsString string
		if err := rows.Scan(&r.ID, &r.Channel.ID, &r.Channel.Name, &r.Channel.Type, &r.Reporter.ID, &r.Reporter.Name, &r.Target.ID, &r.Target.Name, &r.Reason, &logsString, &r.Time, &r.Handler.ID, &r.Handler.Name, &r.Action, &r.ActionDuration, &r.TimeHandled); err != nil {
			fmt.Println("error when scanning row:", err)
			utils.WebWriteError(w, 500, "Internal error")
			return
		}
		r.Logs = strings.Split(logsString, "\n")

		response.Reports = append(response.Reports, r)
	}

	utils.WebWrite(w, response)
}
