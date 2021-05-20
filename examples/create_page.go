package main

import (
	"context"
	"fmt"

	"github.com/kjk/notion"
)

func createPage(apiKey string, parentID string) {
	fmt.Printf("createPage\n")

	if parentID == "" {
		// https://www.notion.so/Test-pages-for-notionapi-0367c2db381a4f8b9ce360f388a6b2e3
		parentID = "0367c2db381a4f8b9ce360f388a6b2e3"
	}

	c := getClient(apiKey)
	ctx := context.Background()
	titleText := &notion.Text{
		Content: "this is a title of the created page",
	}
	title := []notion.RichText{
		{
			Type: notion.RichTextTypeText,
			Text: titleText,
		},
	}
	params := notion.CreatePageParams{
		ParentID:   parentID,
		ParentType: notion.ParentTypePage,
		Title:      title,
	}
	rsp, err := c.CreatePage(ctx, params)
	if err != nil {
		logf("CreatePage() failed with '%s'\n", err)
		logf("rsp.RawJSON: '%s'\n", rsp.RawJSON)
		ppJSON(rsp.RawJSON)
		return
	}
	logf("Created a page!\n")
	showPageInfo(rsp)
}
