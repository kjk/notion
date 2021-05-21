package main

import (
	"context"

	"github.com/kjk/notion"
)

// Test database: https://www.notion.so/b8d975b27cdd441da97e035ecbb04ee7
// b8d975b27cdd441da97e035ecbb04ee7

func showDatabaseProperty(name string, prop notion.DatabaseProperty) {
	logf("    property: '%s'\n", name)
	logf("      id: %s\n", prop.ID)
	logf("      type: %s\n", prop.Type)
	if prop.Number != nil {
		num := prop.Number
		logf("      format: %s\n", num.Format)
	} else if prop.Select != nil {
		sel := prop.Select
		for _, selopt := range sel.Options {
			logf("      sel opt: %s\n", selopt.ID)
			logf("        name: %s\n", selopt.Name)
			logf("        color: %s\n", selopt.Color)
		}
	} else if prop.MultiSelect != nil {
		msel := prop.MultiSelect
		for _, selopt := range msel.Options {
			logf("      sel opt: %s\n", selopt.ID)
			logf("        name: %s\n", selopt.Name)
			logf("        color: %s\n", selopt.Color)
		}
	} else if prop.Formula != nil {
		f := prop.Formula
		logf("      expression: %s\n", f.Expression)
	} else if prop.Relation != nil {
		r := prop.Relation
		logf("    database id: %s\n", r.DatabaseID)
		logf("    synced prop idd: %s\n", r.SyncedPropID)
		logf("    synced prop name: %s\n", r.SyncedPropName)
	} else if prop.Rollup != nil {
		r := prop.Rollup
		logf("      relation prop id: %s\n", r.RelationPropID)
		logf("      relation prop name: %s\n", r.RelationPropName)
		logf("      rollup prop id: %s\n", r.RollupPropID)
		logf("      rollup prop name: %s\n", r.RollupPropName)
		logf("      function: %s\n", r.Function)
	}
}

func showDatabaseInfo(db *notion.Database) {
	logf("database:\n")
	logf("  ID: '%s'\n", db.ID)
	logf("  CreatedTime: '%s'\n", db.CreatedTime)
	logf("  LastEditedTime: '%s'\n", db.LastEditedTime)
	showRichText(1, "Title", db.Title)
	logf("  %d properties:\n", len(db.Properties))
	for name, prop := range db.Properties {
		showDatabaseProperty(name, prop)
	}
}

func getDabaseInfo2(apiKey string, id string) {
	logf("getDatabaseInfo: id='%s'\n", id)
	c := getClient(apiKey)
	ctx := context.Background()
	db, err := c.GetDatabase(ctx, id)
	if err != nil {
		logf("c.GetDatabase() failed with '%s'\n", err)
		logf("db.RawJSON: '%s'\n", db.RawJSON)
		ppJSON(db.RawJSON)
		return
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
