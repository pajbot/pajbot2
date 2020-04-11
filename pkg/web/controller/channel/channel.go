package channel

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/router"
	"github.com/pajbot/pajbot2/pkg/web/views"
)

func Load(a pkg.Application, cfg *config.Config) {
	m := router.Subrouter("/c/{channel:[a-zA-Z0-9]+}")

	router.RGet(m, "/dashboard", handleDashboard(a))
}

func handleDashboard(a pkg.Application) func(w http.ResponseWriter, r *http.Request) {
	type BotInfo struct {
		Name      string
		Connected bool
	}

	type ChannelInfo struct {
		ID   string
		Name string
	}

	type Extra struct {
		Bots    []BotInfo
		Channel ChannelInfo
	}

	return func(w http.ResponseWriter, r *http.Request) {
		extra := &Extra{}
		vars := mux.Vars(r)
		channel, ok := vars["channel"]
		if !ok {
			return
		}

		for it := a.TwitchBots().Iterate(); it.Next(); {
			bot := it.Value()

			bc := bot.GetBotChannel(channel)
			if bc == nil {
				continue
			}

			extra.Channel.ID = bc.Channel().GetID()
			extra.Channel.Name = bc.Channel().GetName()

			// fmt.Fprintf(w, "Bot: %s\n", bot.TwitchAccount().ID())
			bi := BotInfo{
				Name:      bot.TwitchAccount().Name(),
				Connected: bot.Connected(),
			}
			extra.Bots = append(extra.Bots, bi)
		}

		extraBytes, _ := json.Marshal(extra)

		err := views.RenderExtra("dashboard", w, r, extraBytes)
		if err != nil {
			log.Println("Error rendering dashboard view:", err)
			return
		}
	}
}
