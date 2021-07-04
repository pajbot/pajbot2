package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pajbot/pajbot2/pkg/apirequest"
	pb2twitch "github.com/pajbot/pajbot2/pkg/twitch"
	"golang.org/x/oauth2"
)

type botConfig struct {
	databaseID  int
	account     *pb2twitch.TwitchAccount
	tokenSource oauth2.TokenSource
	token       *oauth2.Token
}

func newBotConfig(databaseID int, account *pb2twitch.TwitchAccount, credentials pb2twitch.BotCredentials, oauthConfig *oauth2.Config) botConfig {
	databaseToken := &oauth2.Token{
		AccessToken:  credentials.AccessToken,
		TokenType:    "bearer",
		RefreshToken: credentials.RefreshToken,
		Expiry:       credentials.Expiry.Time,
	}

	return botConfig{
		databaseID:  databaseID,
		account:     account,
		tokenSource: oauth2.ReuseTokenSource(databaseToken, oauthConfig.TokenSource(context.Background(), databaseToken)),
		token:       databaseToken,
	}
}

func (bc *botConfig) Validate(sqlClient *sql.DB) error {
	token, err := bc.tokenSource.Token()
	if err != nil {
		return fmt.Errorf("error validating botconfig for %s: %w", bc.account.Name(), err)
	}

	// Update token in database
	const queryF = `
UPDATE "bot"
SET
	twitch_access_token = $1,
	twitch_refresh_token = $2,
	twitch_access_token_expiry = $3
WHERE
id=$4`
	_, err = sqlClient.Exec(queryF, token.AccessToken, token.RefreshToken, token.Expiry, bc.databaseID)
	if err != nil {
		return err
	}

	bc.token = token

	isValid, self, err := apirequest.TwitchWrapper.HelixBot().ValidateToken(token.AccessToken)
	if err != nil {
		return fmt.Errorf("error validating oauth token for bot '%s': %w", bc.account.Name(), err)
	}
	if !isValid {
		return fmt.Errorf("token for '%s' is invalid", bc.account.Name())
	}

	if self.Data.UserID != bc.account.ID() {
		return fmt.Errorf("mismatching user ID for %s (%s) - doesn't match the API response (%s)", bc.account.Name(), bc.account.ID(), self.Data.UserID)
	}

	return nil
}
