package main

import (
	"context"

	"github.com/kjk/notion"
)

func showUser(user *notion.User) {
	logf("  user: %s\n", user.ID)
	logf("    type: %s\n", user.Type)
	logf("    name: %s\n", user.Name)
	if user.Person != nil {
		logf("    email: %s\n", user.Person.Email)
	}
}

func showListUsers(rsp *notion.ListUsersResponse) {
	logf("showListUsers:\n")
	logf("  hasMore: %v\n", rsp.HasMore)
	if rsp.NextCursor != "" {
		logf("  nextCursor: '%s'\n", rsp.NextCursor)
	}
	logf("  nextCursor: %s\n", rsp.NextCursor)
	logf("  %d users:\n", len(rsp.Results))
	for _, u := range rsp.Results {
		showUser(&u)
	}
}

func listUsers(apiKey string) {
	logf("listUsers:\n")

	c := getClient(apiKey)
	ctx := context.Background()
	res, err := c.ListUsers(ctx, nil)
	if err != nil {
		logf("ListUsers() failed with: '%s'\n", err)
		logf("res.RawJSON: '%s'\n", res.RawJSON)
		ppJSON(res.RawJSON)
		return
	}
	showListUsers(res)
}
