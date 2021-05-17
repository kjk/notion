package main

import (
	"context"
	"fmt"

	"github.com/kjk/notion"
)

func showBlockChildren(bcr *notion.BlockChildrenResponse) {
	logf("showBlockChildren:\n")
	logf("  %d children:\n", len(bcr.Results))
	for _, b := range bcr.Results {
		logf("  %v\n", b)
	}
}

func getBlockChildren2(apiKey string, blockID string) {
	fmt.Printf("getBlockChildren: blockID='%s'\n", blockID)

	c := getClient(apiKey)
	ctx := context.Background()
	rsp, err := c.GetBlockChildren(ctx, blockID, nil)
	if err != nil {
		logf("page.RawJSON: '%s'\n", rsp.RawJSON)
		ppJSON(rsp.RawJSON)
		must(err)
	}
	showBlockChildren(rsp)
}

func getBlockChildren(apiKey string, blockID string) {
	if blockID == "" {
		// show info about regular page
		getBlockChildren2(apiKey, "0367c2db381a4f8b9ce360f388a6b2e3")
		// show info about page in a database
		getBlockChildren2(apiKey, "e56b74a6398a43848137cca2a0de20b2")
	} else {
		getBlockChildren2(apiKey, blockID)
	}
}
