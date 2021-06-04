package web

// TODO: Move these to webhook/github.go or something

import (
	"encoding/json"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pajbot/pajbot2/pkg"
)

type followResponse struct {
	Data []struct {
		FromID     string `json:"from_id"`
		ToID       string `json:"to_id"`
		FollowedAt string `json:"followed_at"`
	} `json:"data"`
}

func apiCallbacksFollow(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	if challenge != "" {
		fmt.Fprint(w, html.EscapeString(challenge))
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var response followResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	for _, follow := range response.Data {
		fmt.Printf("User with id %s followed %s at %s\n", follow.FromID, follow.ToID, follow.FollowedAt)
	}
}

type streamsResponse struct {
	Data []struct {
		ID           string        `json:"id"`
		UserID       string        `json:"user_id"`
		GameID       string        `json:"game_id"`
		CommunityIds []interface{} `json:"community_ids"`
		Type         string        `json:"type"`
		Title        string        `json:"title"`
		ViewerCount  int           `json:"viewer_count"`
		StartedAt    time.Time     `json:"started_at"`
		Language     string        `json:"language"`
		ThumbnailURL string        `json:"thumbnail_url"`
	} `json:"data"`
}

func apiCallbacksStreams(w http.ResponseWriter, r *http.Request) {
	challenge := r.URL.Query().Get("hub.challenge")
	if challenge != "" {
		fmt.Fprint(w, html.EscapeString(challenge))
		fmt.Println("Responding to streams")
		return
	}

	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")
	fmt.Printf("Streams response xd \n")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	var response streamsResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		panic(err)
	}

	if len(response.Data) > 0 {
		fmt.Printf("%#v\n", response.Data)
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
		fmt.Printf("Online!\n")
	} else {
		fmt.Printf("%#v\n", response.Data)
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
		fmt.Printf("Offline!\n")
	}
}

func handlePush(b pkg.Sender, body []byte, p *customPayload) {
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

func writeCommit(b pkg.Sender, commit Commit, repository RepositoryData) {
	// msg := fmt.Sprintf("%s (%s) committed to %s (%s): %s %s", commit.Author.Name, commit.Author.Username, repository.Name, commit.Timestamp, commit.Message, commit.URL)
	// XXX: Missing channel
	// b.Say(msg)
}

func handleStatus(b pkg.Sender, body []byte, p *customPayload) {
	var data StatusHookResponse

	err := json.Unmarshal(body, &data)
	if err != nil {
		p.Add("error", "Json Unmarshal error: "+err.Error())
		return
	}

	switch data.State {
	case "pending":
		// TODO: Re-implement
		// b.SaySafef("Build for %s just started", data.Repository.Name)

	case "success":
		// TODO: Re-implement
		// b.SaySafef("Build for %s succeeded! FeelsGoodMan", data.Repository.Name)

	case "error":
		fallthrough

	case "failure":
		// TODO: Re-implement
		// b.SaySafef("Build for %s failed: %s FeelsBadMan", data.Repository.Name, data.TargetURL)
	}

	p.Add("success", true)
}
