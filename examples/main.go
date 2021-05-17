package main

import (
	"flag"

	"github.com/kjk/notion"
	"github.com/kjk/u"
)

// apiKey is needed to access pages. It comes from an internal integration
// I've created for my test pages
// Read https://developers.notion.com/docs/getting-started and
// https://developers.notion.com/docs/authorization
// to learn more about authorization
const apiKeyForTestPages = "secret_P6H8sOo7kkN6fn4Jd70Axlc0tVBG9bOnL8ZxTprm7rR"

var must = u.Must
var logf = u.Logf
var panicIf = u.PanicIf

func getClient(apiKey string) *notion.Client {
	return notion.NewClient(apiKey)
}

func main() {
	var (
		flgGetPageInfo      bool
		flgGetBlockChildren bool

		flgAPIKey string
		flgID     string
	)
	{
		flag.BoolVar(&flgGetPageInfo, "get-page-info", false, "get information about a page")
		flag.BoolVar(&flgGetBlockChildren, "get-block-children", false, "get children of a block")
		flag.StringVar(&flgID, "id", "", "id of page or block (if not using default test pages)")
		flag.StringVar(&flgAPIKey, "api-key", "", "api key for authentication (if not using default test page)")
		flag.Parse()

		if flgAPIKey == "" {
			flgAPIKey = apiKeyForTestPages
		}
	}

	if flgGetPageInfo {
		getPageInfo(flgAPIKey, flgID)
		return
	}

	if flgGetBlockChildren {
		getBlockChildren(flgAPIKey, flgID)
		return
	}

	flag.Usage()
}
