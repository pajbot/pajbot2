package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/web/router"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"github.com/pajbot/pajbot2/pkg/web/views"
)

func Load() {
	router.Get("/dashboard", Dashboard)
}

const dashboardPermissions = pkg.PermissionAdmin | pkg.PermissionReport | pkg.PermissionModeration | pkg.PermissionReportAPI

func Dashboard(w http.ResponseWriter, r *http.Request) {
	state := state.Context(w, r)

	const queryF = `
SELECT
	twitch_user_channel_permission.channel_id,
	permissions
FROM
	"user"
LEFT JOIN
	twitch_user_channel_permission ON "user".twitch_userid="twitch_user_channel_permission".twitch_user_id
WHERE twitch_userid=$1`

	type ChannelInfo struct {
		ID   string
		Name string
	}

	type Extra struct {
		Channels []ChannelInfo
	}

	var extraBytes []byte

	if state.Session != nil {
		rows, err := state.SQL.Query(queryF, state.Session.TwitchUserID)
		if err != nil {
			// TODO: render error page somehow?
			fmt.Println("ERROR 1:", err)
			return
		}
		defer rows.Close()

		extra := &Extra{}

		for rows.Next() {
			var twitchChannelID string
			var permission pkg.Permission

			err := rows.Scan(&twitchChannelID, &permission)
			if err != nil {
				// TODO: render error page somehow?
				fmt.Println("ERROR 2:", err)
				return
			}

			if (permission & dashboardPermissions) != 0 {
				extra.Channels = append(extra.Channels, ChannelInfo{
					ID:   twitchChannelID,
					Name: state.TwitchUserStore.GetName(twitchChannelID),
				})
			}
			extraBytes, _ = json.Marshal(extra)
		}
	}

	err := views.RenderExtra("dashboard", w, r, extraBytes)
	if err != nil {
		log.Println("Error rendering dashboard view:", err)
	}
}
