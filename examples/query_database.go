package main

import (
	"context"

	"github.com/kjk/notion"
)

// Test database: https://www.notion.so/b8d975b27cdd441da97e035ecbb04ee7
// b8d975b27cdd441da97e035ecbb04ee7

func showDatabaseQueryInfo(dqr *notion.DatabaseQueryResponse) {
	logf("showDatabaseQueryInfo:\n")
	logf("  hasMore: %v\n", dqr.HasMore)
	if dqr.NextCursor != "" {
		logf("  nextCursor: %s\n", dqr.NextCursor)
	}
	logf("  %d rows:\n", len(dqr.Results))
	for _, p := range dqr.Results {
		showPageInfo(&p)
	}
}

func queryDatabase2(apiKey string, id string) {
	logf("queryDatabase: id='%s'\n", id)
	c := getClient(apiKey)
	ctx := context.Background()
	dqr, err := c.QueryDatabase(ctx, id, nil)
	if err != nil {
		logf("QueryDatabase() failed with '%s'\n", err)
		logf("RawJSON: '%s'\n", dqr.RawJSON)
		ppJSON(dqr.RawJSON)
		return
	}
	showDatabaseQueryInfo(dqr)
}

func queryDatabase(apiKey string, id string) {
	if id == "" {
		queryDatabase2(apiKey, "b8d975b27cdd441da97e035ecbb04ee7")
	} else {
		queryDatabase2(apiKey, id)
	}
}
