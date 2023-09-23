package webhook

import (
	"bufio"
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

func subSlice[S ~[]E, E any](s S, upTo int) S {
	if upTo <= 0 {
		return s[:0]
	}
	l := len(s)
	if upTo > l {
		return s[:l]
	}
	return s[:upTo]
}

func GenerateTwitchMessages(pushData PushHookResponse) []string {
	targetBranch := strings.TrimPrefix(pushData.Ref, "refs/heads/")
	if len(targetBranch) == 0 {
		targetBranch = strings.TrimPrefix(pushData.BaseRef, "refs/heads/")
	}

	repositoryName := pushData.Repository.Name

	if len(targetBranch) == 0 {
		log.Println("Unable to figure out branch name for this push:", pushData)
		return nil
	}

	if strings.Contains(targetBranch, "/") {
		log.Println("Ignoring push for branch", targetBranch)
		// Skip any branches that contain a / - they are most likely a feature branch
		return nil
	}

	messages := []string{}

	for _, commit := range subSlice(pushData.Commits, 5) {
		var sb strings.Builder
		if _, err := sb.WriteString(commit.Author.Username); err != nil {
			log.Println("ERROR WRITING TO STRING:", err)
			continue
		}

		// TODO: parse other authors
		scanner := bufio.NewScanner(strings.NewReader(commit.Message))
		commitMessage := ""
		coAuthors := []string{}
		for scanner.Scan() {
			line := scanner.Text()
			if commitMessage == "" {
				commitMessage = line
			} else {
				const coAuthorPrefix = "co-authored-by: "
				const coAuthorPrefixLen = len(coAuthorPrefix)
				if strings.HasPrefix(strings.ToLower(line), coAuthorPrefix) {
					author := line[coAuthorPrefixLen:]
					emailStartIndex := strings.Index(author, " <")
					if emailStartIndex != -1 {
						author = strings.TrimSpace(author[:emailStartIndex])
						if len(author) > 1 {
							// author must be at least 1 character long
							coAuthors = append(coAuthors, author)
						}
					}
				}
			}
		}

		if len(coAuthors) > 0 {
			if _, err := sb.WriteString(" (with "); err != nil {
				log.Println("ERROR WRITING TO STRING:", err)
				continue
			}
			if _, err := sb.WriteString(strings.Join(subSlice(coAuthors, 5), ", ")); err != nil {
				log.Println("ERROR WRITING TO STRING:", err)
				continue
			}
			if _, err := sb.WriteString(")"); err != nil {
				log.Println("ERROR WRITING TO STRING:", err)
				continue
			}
		}

		// commitMessage := strings.SplitN(commit.Message, "\n", 2)[0]

		if _, err := sb.WriteString(fmt.Sprintf(" committed to %s@%s (%s): %s", repositoryName, targetBranch, commit.Timestamp, commitMessage)); err != nil {
			log.Println("ERROR WRITING TO STRING:", err)
			continue
		}

		if sb.Len() < 350 {
			// Only append the commit hash if the commit message is pretty small
			if _, err := sb.WriteString(fmt.Sprintf(" %s", commit.URL)); err != nil {
				log.Println("ERROR WRITING TO STRING:", err)
				continue
			}
		}

		messages = append(messages, sb.String())
	}

	return messages
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

			messages := GenerateTwitchMessages(pushData)

			delay := 100
			for _, message := range messages {
				func(message string) {
					time.AfterFunc(time.Millisecond*time.Duration(delay), func() {
						botChannel.Say(message)
					})
				}(message)
				delay += 250
			}
		default:
			log.Println("Unhandled hook type:", hookType)
		}

		w.WriteHeader(http.StatusOK)
	}
}
