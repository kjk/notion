package main

import "context"

func getUser2(apiKey string, id string) {
	logf("getUser: '%s'\n", id)

	c := getClient(apiKey)
	ctx := context.Background()
	res, err := c.GetUser(ctx, id)
	if err != nil {
		logf("GetUser() failed with: '%s'\n", err)
		logf("res.RawJSON: '%s'\n", res.RawJSON)
		ppJSON(res.RawJSON)
		return
	}
	showUser(res)
}

func getUser(apiKey string, id string) {
	if id == "" {
		// user id of 'api client test' bot
		getUser2(apiKey, "c9cebcd2-9fc0-4092-aa6e-c2b505c57021")
	} else {
		getUser2(apiKey, id)
	}
}
