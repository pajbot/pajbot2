package main

import (
	"fmt"
	"os"
)

type fixType func(app *Application) error

var fixes []fixType

func fixCmd() {
	fixes = append(fixes, fix1)

	app := newApplication()

	err := app.LoadConfig(*configPath)
	if err != nil {
		fmt.Println("An error occurred while loading the config file: ", err)
		os.Exit(1)
	}

	err = app.InitializeAPIs()
	if err != nil {
		fmt.Println("An error occurred while initializing APIs: ", err)
		os.Exit(1)
	}

	err = app.InitializeSQL()
	if err != nil {
		fmt.Println("Error starting SQL client:", err)
		os.Exit(1)
	}

	for i, fix := range fixes {
		fmt.Printf("Fix #%d: Running\n", i+1)
		err := fix(app)
		if err != nil {
			fmt.Printf("Fix #%d: Error: %s\n", i+1, err.Error())
			break
		}
	}
}

// Fix #1: Migration failed in line 0: CREATE UNIQUE INDEX `itwitchuid` ON Bot(twitch_userid); (details: Error 1062: Duplicate entry '' for key 'itwitchuid')")
func fix1(app *Application) error {
	db := app.SQL()
	const queryF = `SELECT version, dirty FROM schema_migrations`
	row := db.QueryRow(queryF)
	var version int64
	var dirty bool
	err := row.Scan(&version, &dirty)
	if err != nil {
		return err
	}

	if version != 20190118204509 || !dirty {
		fmt.Println("Fix #1: Skipping (Schema version and dirty value are not correct)")
		return nil
	}

	rows, err := db.Query("SELECT id, twitch_userid, name FROM Bot")
	if err != nil {
		return err
	}
	for rows.Next() {
		var id int64
		var oldUserID string
		var name string
		err = rows.Scan(&id, &oldUserID, &name)
		if err != nil {
			return err
		}

		if oldUserID != "" {
			fmt.Printf("Skipping %s because it already has an ID set (%s)\n", name, oldUserID)
			continue
		}

		twitchUserID := app.UserStore().GetID(name)
		if twitchUserID == "" {
			return fmt.Errorf("Unable to get ID for user '%s' - fix and rerun", name)
		}

		_, err = db.Exec("UPDATE bot SET twitch_userid=$1 WHERE id=$2", twitchUserID, id)
		if err != nil {
			fmt.Println("Error updating user ID:", err)
			return err
		}

		fmt.Println("Updated User ID for", name, "to", twitchUserID)
	}

	fmt.Println("Attempting to revert the schema migration (maybe use migration .Down here or something)")

	_, err = db.Exec("UPDATE schema_migrations SET version=20190118000448, dirty=0")
	if err != nil {
		return err
	}

	fmt.Println("Done! now just run ./bot like normal")
	return nil
}
