package logout

import (
	"fmt"
	"net/http"

	"github.com/pajlada/pajbot2/web/router"
	"github.com/pajlada/pajbot2/web/state"
)

func Load() {
	router.Get("/logout", handleLogout)
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	c := state.Context(w, r)
	if c.SessionID != nil {
		sessionID := *c.SessionID

		if state.IsValidSessionID(sessionID) {
			// Remove session from database
			const queryF = `DELETE FROM UserSession WHERE id=?`

			_, err := c.SQL.Exec(queryF, sessionID)
			if err != nil {
				fmt.Println("Error deleting session ID")
			}
		}

		state.ClearSessionCookies(w)
	}

	// Redirect user
	redirectURL := r.FormValue("redirect")

	if redirectURL == "" {
		redirectURL = "/"
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}
