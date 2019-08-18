package state

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"
)

const SessionIDCookie = "pb2sessionid"

type SessionStore struct {
}

func getCookie(r *http.Request, cookieName string) *string {
	c, _ := r.Cookie(cookieName)
	if c == nil {
		return nil
	}

	return &c.Value
}

func IsValidSessionID(sessionID string) bool {
	// TODO: Make this more sophisticated
	return sessionID != ""
}

func SetSessionCookies(w http.ResponseWriter, sessionID string, twitchUserName string) {
	sessionIDCookie := &http.Cookie{
		Name:  SessionIDCookie,
		Value: sessionID,
		Path:  "/",
	}
	http.SetCookie(w, sessionIDCookie)

	userNameCookie := &http.Cookie{
		Name:  "pb2username",
		Value: twitchUserName,
		Path:  "/",
	}
	http.SetCookie(w, userNameCookie)
}

func ClearSessionCookies(w http.ResponseWriter) {
	sessionIDCookie := &http.Cookie{
		Name:  SessionIDCookie,
		Value: "",
		Path:  "/",

		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, sessionIDCookie)

	userNameCookie := &http.Cookie{
		Name:  "pb2username",
		Value: "",
		Path:  "/",

		Expires: time.Unix(0, 0),
	}
	http.SetCookie(w, userNameCookie)
}

func (s *SessionStore) Get(db *sql.DB, sessionID string) *Session {
	const queryF = `SELECT user_session.user_id, "user".twitch_userid, "user".twitch_username FROM user_session INNER JOIN "user" ON "user".id=user_session.user_id WHERE user_session.id=$1`

	var userID uint64
	var twitchUserID string
	var twitchUserName string

	err := db.QueryRow(queryF, sessionID).Scan(&userID, &twitchUserID, &twitchUserName)
	switch {
	case err == sql.ErrNoRows:
		// invalid session ID
		return nil
	case err != nil:
		fmt.Println("SQL Error:", err)
		// some query or mysql error occurred
		return nil
	default:
		return &Session{
			UserID:         userID,
			TwitchUserID:   twitchUserID,
			TwitchUserName: twitchUserName,
		}
	}
}
