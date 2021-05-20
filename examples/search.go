package main

import (
	"context"

	"github.com/kjk/notion"
)

func search(apiKey string) {
	logf("search:\n")
	c := getClient(apiKey)
	ctx := context.Background()
	res, err := c.Search(ctx, nil)
	if err != nil {
		logf("Search() failed with '%s'\n", err)
		logf("res.RawJSON: '%s'\n", res.RawJSON)
		ppJSON(res.RawJSON)
		return
	}
	showSearchResponse(res)
}

func showSearchResponse(sr *notion.SearchResponse) {
	logf("showSearchResponse:\n")
	logf("  hasMore: %v\n", sr.HasMore)
	if sr.NextCursor != "" {
		logf("  nextCursor: %s\n", sr.NextCursor)
	}
	for _, res := range sr.Results {
		showSearchResult(res)
	}
}

func showSearchResult(sr interface{}) {
	logf("  showSearchResult: %T\n", sr)
}
