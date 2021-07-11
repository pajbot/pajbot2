package twitch

import (
	"net/http"

	"github.com/nicklaw5/helix"
	"github.com/pajbot/pajbot2/pkg/web/state"
	"github.com/pajbot/utils"
	"golang.org/x/oauth2"
)

func onBotAuthenticated(
	w http.ResponseWriter, r *http.Request,
	self *helix.ValidateTokenResponse, oauth2Token *oauth2.Token, stateData *stateData) {
	const queryF = `
INSERT INTO bot (twitch_userid, twitch_username, twitch_access_token, twitch_refresh_token, twitch_access_token_expiry)
    VALUES ($1, $2, $3, $4, $5) ON CONFLICT (twitch_userid)
    DO
    UPDATE
    SET
        twitch_username = $2,
        twitch_access_token = $3,
        twitch_refresh_token = $4,
        twitch_access_token_expiry = $5`
	c := state.Context(w, r)
	_, err := c.SQL.Exec(queryF, self.Data.UserID, self.Data.Login, oauth2Token.AccessToken, oauth2Token.RefreshToken, oauth2Token.Expiry) // GOOD
	if err != nil {
		_ = utils.WebWriteError(w, 500, "Error inserting bot, admin should check console logs NaM")
		return
	}

	_, _ = w.Write([]byte("Bot added/updated! Restart the bot for the changes to take effect"))
}
