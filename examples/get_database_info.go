package main

import (
	"context"

	"github.com/kjk/notion"
)

// Test database: https://www.notion.so/b8d975b27cdd441da97e035ecbb04ee7
// b8d975b27cdd441da97e035ecbb04ee7

func showDatabaseInfo(db *notion.Database) {
	logf("database:\n")
	logf("  ID: '%s'\n", db.ID)
	logf("  CreatedTime: '%s'\n", db.CreatedTime)
	logf("  LastEditedTime: '%s'\n", db.LastEditedTime)
	showRichText(1, "Title", db.Title)
	logf("  %d properties:\n", len(db.Properties))
	for name, prop := range db.Properties {
		// TODO: better display of properties
		logf("    %s: %v\n", name, prop)
	}
}

func getDabaseInfo2(apiKey string, id string) {
	logf("getDatabaseInfo: id='%s'\n", id)
	c := getClient(apiKey)
	ctx := context.Background()
	db, err := c.GetDatabase(ctx, id)
	if err != nil {
		logf("db.RawJSON: '%s'\n", db.RawJSON)
		ppJSON(db.RawJSON)
		must(err)
	}
	showDatabaseInfo(db)
}

func getDatabaseInfo(apiKey string, id string) {
	if id == "" {
		getDabaseInfo2(apiKey, "b8d975b27cdd441da97e035ecbb04ee7")
		getDabaseInfo2(apiKey, "ffbfda6791d34147b44a57ef83ab907a")
		getDabaseInfo2(apiKey, "3acbc0fae5e34dfa9f3960d91cfb018a")
		getDabaseInfo2(apiKey, "509fe00ee06448249687a4eb26bf9579")

	} else {
		getDabaseInfo2(apiKey, id)
	}
}
