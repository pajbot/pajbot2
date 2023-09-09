package webhook

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/pajbot/pajbot2/pkg"
	"github.com/pajbot/pajbot2/pkg/common/config"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"github.com/pajbot/utils"
)

// Commit github json to go
type Commit struct {
	ID        string    `json:"id"`
	TreeID    string    `json:"tree_id"`
	Distinct  bool      `json:"distinct"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	URL       string    `json:"url"`
	Author    struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"author"`
	Committer struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"username"`
	} `json:"committer"`
}

// RepositoryData xD
type RepositoryData struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	FullName string `json:"full_name"`
	Owner    struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"owner"`
	HTMLURL string `json:"html_url"`
	URL     string `json:"url"`
}

type sender struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	URL   string `json:"url"`
}

type pusher struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// PushHookResponse github json to go
type PushHookResponse struct {
	Commits    []Commit       `json:"commits"`
	HeadCommit Commit         `json:"head_commit"`
	Repository RepositoryData `json:"repository"`
	Ref        string         `json:"ref"`
	BaseRef    string         `json:"base_ref"`
	Pusher     pusher         `json:"pusher"`
	Sender     sender         `json:"sender"`
}

// StatusHookResponse json to struct
type StatusHookResponse struct {
	ID          int            `json:"id"`
	Sha         string         `json:"sha"`
	Name        string         `json:"name"`
	TargetURL   string         `json:"target_url"`
	Context     string         `json:"context"`
	Description string         `json:"description"`
	State       string         `json:"state"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Repository  RepositoryData `json:"repository"`
	Sender      sender         `json:"sender"`
}

func signBody(secret, body []byte) []byte {
	computed := hmac.New(sha1.New, secret)
	computed.Write(body)
	return []byte(computed.Sum(nil))
}

func verifySignature(secretString string, signature string, body []byte) bool {
	const signaturePrefix = "sha1="
	const signatureLength = 45

	if len(signature) != signatureLength || !strings.HasPrefix(signature, signaturePrefix) {
		return false
	}

	secret := []byte(secretString)
	actual := make([]byte, 20)
	hex.Decode(actual, []byte(signature[5:]))

	return hmac.Equal(signBody(secret, body), actual)
}

func apiGithub(cfg *config.AuthGithubWebhook) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		c := state.Context(w, r)

		hookType := r.Header.Get("x-github-event")
		hookSignature := r.Header.Get("x-hub-signature")

		v := mux.Vars(r)
		channelID := v["channelID"]

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println(err)
			return
		}
		defer r.Body.Close()

		verified := verifySignature(cfg.Secret, hookSignature, body)

		if !verified {
			utils.WebWriteError(w, 400, "bad secret")
			return
		}

		var botChannel pkg.BotChannel

		if channelID != "" {
			for it := c.Application.TwitchBots().Iterate(); it.Next(); {
				bot := it.Value()
				if bot == nil {
					continue
				}

				botChannel = bot.GetBotChannelByID(channelID)
				if botChannel == nil {
					continue
				}

				break
			}
		}

		if botChannel == nil {
			log.Println("No bot channel active to handle this request", channelID)
			// No bot channel active to handle this request
			return
		}

		// TODO: handle event type
		switch hookType {
		case "push":
			var pushData PushHookResponse

			err := json.Unmarshal(body, &pushData)
			if err != nil {
				utils.WebWriteError(w, 400, "bad push data")
				return
			}

			targetBranch := strings.TrimPrefix(pushData.Ref, "refs/heads/")
			if len(targetBranch) == 0 {
				targetBranch = strings.TrimPrefix(pushData.BaseRef, "refs/heads/")
			}

			if len(targetBranch) == 0 {
				log.Println("Unable to figure out branch name for this push:", pushData)
				break
			}

			if strings.Contains(targetBranch, "/") {
				log.Println("Ignoring push for branch", targetBranch)
				// Skip any branches that contain a / - they are most likely a feature branch
				break
			}

			delay := 100
			for _, commit := range pushData.Commits {
				func(iCommit Commit) {
					time.AfterFunc(time.Millisecond*time.Duration(delay), func() {
						botChannel.Say(fmt.Sprintf("%s (%s) committed to %s@%s (%s): %s %s", iCommit.Author.Name, iCommit.Author.Username, pushData.Repository.Name, targetBranch, iCommit.Timestamp, iCommit.Message, iCommit.URL))
					})
				}(commit)
				delay += 250
			}
		default:
			log.Println("Unhandled hook type:", hookType)
		}

		w.WriteHeader(http.StatusOK)
	}
}
