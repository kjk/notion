package main

import (
	"context"

	"github.com/kjk/notion"
	"github.com/tidwall/pretty"
)

// shows how to use https://developers.notion.com/reference/get-page API

// shows the info about:
//   regular page https://www.notion.so/Test-pages-for-notionapi-0367c2db381a4f8b9ce360f388a6b2e3
//   page in a database https://www.notion.so/A-row-that-is-not-empty-page-e56b74a6398a43848137cca2a0de20b2
// or the page given with -id argument (in which case also needs )

func getIndent(n int) string {
	s := ""
	for n > 0 {
		n -= 1
		s += "  "
	}
	return s
}

func showRichText(indent int, name string, richText []notion.RichText) {
	s := getIndent(indent)
	// TODO: better implementation
	if name != "" {
		logf("%s%s: %v\n", s, name, richText)
		return
	}
	logf("%s%v\n", s, richText)
}

func ppJSON(d []byte) {
	res := pretty.Pretty(d)
	logf("pretty printed JSON:\n%s\n", res)
}

func showPageInfo(page *notion.Page) {
	logf("showPageInfo:\n")
	logf("  ID: '%s'\n", page.ID)
	logf("  CreatedTime: '%s'\n", page.CreatedTime)
	logf("  LastEditedTime: '%s'\n", page.LastEditedTime)
	if page.Parent.PageID != nil {
		logf("  Parent: page with ID '%s'\n", *page.Parent.PageID)
	} else if page.Parent.DatabaseID != nil {
		logf("  Parent: database with ID '%s'\n", *page.Parent.DatabaseID)
	} else {
		panicIf(true, "both page.Parent.PageID or page.Parent.DatabaseID are nil")
	}
	logf("  Archived: %v\n", page.Archived)
	switch prop := page.Properties.(type) {
	case notion.PageProperties:
		logf("  page properties:\n")
		showRichText(2, "Title", prop.Title.Title)
	case notion.DatabasePageProperties:
		logf("  database properties (NYI):\n")
	}
}

func getPageInfo2(apiKey string, pageID string) {
	logf("getPageInfo: pageID='%s'\n", pageID)

	c := getClient(apiKey)
	ctx := context.Background()
	page, err := c.GetPage(ctx, pageID)
	if err != nil {
		logf("GetPage() failed with: '%s'\n", err)
		logf("page.RawJSON: '%s'\n", page.RawJSON)
		ppJSON(page.RawJSON)
		return
	}
	showPageInfo(page)
}

func getPageInfo(apiKey string, pageID string) {
	if pageID == "" {
		// show info about regular page
		getPageInfo2(apiKey, "0367c2db381a4f8b9ce360f388a6b2e3")
		// show info about page in a database
		getPageInfo2(apiKey, "e56b74a6398a43848137cca2a0de20b2")
	} else {
		getPageInfo2(apiKey, pageID)
	}
}
