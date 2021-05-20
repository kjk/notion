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
	return notion.NewClient(apiKey, nil)
}

func main() {
	var (
		flgGetPageInfo      bool
		flgGetBlockChildren bool
		flgGetDatabase      bool
		flgListUsers        bool
		flgGetUser          bool
		flgQueryDatabase    bool
		flgSearch           bool

		flgAPIKey string
		flgID     string
	)
	{
		flag.BoolVar(&flgGetPageInfo, "get-page-info", false, "get information about a page (use -id for page id)")
		flag.BoolVar(&flgGetBlockChildren, "get-block-children", false, "get children of a block (use -id for parent id)")
		flag.BoolVar(&flgGetDatabase, "get-database", false, "get database information (use -id for database id)")
		flag.BoolVar(&flgListUsers, "list-users", false, "list users")
		flag.BoolVar(&flgGetUser, "get-user", false, "get user (use -id for user id)")
		flag.BoolVar(&flgQueryDatabase, "query-db", false, "query database (use -id for database id)")
		flag.BoolVar(&flgSearch, "search", false, "search")
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

	if flgGetDatabase {
		getDatabaseInfo(flgAPIKey, flgID)
		return
	}

	if flgListUsers {
		listUsers(flgAPIKey)
		return
	}

	if flgGetUser {
		getUser(flgAPIKey, flgID)
		return
	}

	if flgQueryDatabase {
		queryDatabase(flgAPIKey, flgID)
		return
	}

	if flgSearch {
		search(flgAPIKey)
		return
	}

	flag.Usage()
}
