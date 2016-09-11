package web

import "time"

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

// PushHookResponse github json to go
type PushHookResponse struct {
	Commits    []Commit       `json:"commits"`
	HeadCommit Commit         `json:"head_commit"`
	Repository RepositoryData `json:"repository"`
	Pusher     struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"pusher"`
	Sender struct {
		Login string `json:"login"`
		ID    int    `json:"id"`
		URL   string `json:"url"`
	} `json:"sender"`
}
