package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/pajlada/pajbot2/bot"
)

func apiHook(w http.ResponseWriter, r *http.Request) {
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
		write(w, p.data)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		p.Add("error", "Internal error")
		write(w, p.data)
		return
	}

	verified := verifySignature(channelHook.Secret, hookSignature, body)

	if !verified {
		p.Add("error", "Invalid secret")
		write(w, p.data)
		return
	}

	b, _ := twitchBots[channel]

	if b == nil {
		// no bot found for channel
		p.Add("error", "No bot found for channel "+channel)
		write(w, p.data)
		return
	}

	switch hookType {
	case "push":
		//handlePush(b, body, &p)
	case "status":
		//handleStatus(b, body, &p)
	}

	write(w, p.data)
}

func handlePush(b *bot.Bot, body []byte, p *customPayload) {
	var pushData PushHookResponse

	err := json.Unmarshal(body, &pushData)
	if err != nil {
		p.Add("error", "Json Unmarshal error: "+err.Error())
		return
	}

	delay := 100

	for _, commit := range pushData.Commits {
		func(iCommit Commit) {
			time.AfterFunc(time.Millisecond*time.Duration(delay), func() { writeCommit(b, iCommit, pushData.Repository) })
		}(commit)
		delay += 250
	}
	p.Add("success", true)
}

func writeCommit(b *bot.Bot, commit Commit, repository RepositoryData) {
	msg := fmt.Sprintf("%s (%s) committed to %s (%s): %s %s", commit.Author.Name, commit.Author.Username, repository.Name, commit.Timestamp, commit.Message, commit.URL)
	b.SaySafef(msg)
}

func handleStatus(b *bot.Bot, body []byte, p *customPayload) {
	var data StatusHookResponse

	err := json.Unmarshal(body, &data)
	if err != nil {
		p.Add("error", "Json Unmarshal error: "+err.Error())
		return
	}

	switch data.State {
	case "pending":
		b.SaySafef("Build for %s just started", data.Repository.Name)

	case "success":
		b.SaySafef("Build for %s succeeded! FeelsGoodMan", data.Repository.Name)

	case "error":
		fallthrough

	case "failure":
		b.SaySafef("Build for %s failed: %s FeelsBadMan", data.Repository.Name, data.TargetURL)
	}

	p.Add("success", true)
}
