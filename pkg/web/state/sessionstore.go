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
	const queryF = `SELECT UserSession.user_id, User.twitch_userid, User.twitch_username FROM UserSession INNER JOIN User ON User.id=UserSession.user_id WHERE UserSession.id=?`

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
		// some query or mysql error occured
		return nil
	default:
		return &Session{
			UserID:         userID,
			TwitchUserID:   twitchUserID,
			TwitchUserName: twitchUserName,
		}
	}
}
