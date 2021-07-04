package webhook

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/router"
)

func apiHook(w http.ResponseWriter, r *http.Request) {
	/*
		p := customPayload{}
		v := mux.Vars(r)
		hookType := r.Header.Get("x-github-event")
		hookSignature := r.Header.Get("x-hub-signature")
		channel := v["channel"]

		// Get hook from config according to channel
		channelHook, ok := hooks[channel]
		if !ok {
			// No hook for this channel found
			p.Add("error", "No hook found for given channel")
			utils.WebWrite(w, p.data)
			return
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			p.Add("error", "Internal error")
			utils.WebWrite(w, p.data)
			return
		}

		verified := verifySignature(channelHook.Secret, hookSignature, body)

		if !verified {
			p.Add("error", "Invalid secret")
			utils.WebWrite(w, p.data)
			return
		}

		b, _ := twitchBots[channel]

		if b == nil {
			// no bot found for channel
			p.Add("error", "No bot found for channel "+channel)
			utils.WebWrite(w, p.data)
			return
		}

		switch hookType {
		case "push":
			//handlePush(b, body, &p)
		case "status":
			//handleStatus(b, body, &p)
		}

		utils.WebWrite(w, p.data)
	*/
}

func Load(parent *mux.Router, cfg *config.TwitchWebhookConfig) {
	m := parent.PathPrefix("/webhook").Subrouter()

	// NEW AND FRESH AND COOL
	router.RPost(m, `/eventsub`, apiEventsub(cfg))

	// DEPRECATED
	router.RGet(m, `/{channel:\w+}`, apiHook)
}
