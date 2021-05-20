# notion

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/kjk/notion?label=go%20module)](https://github.com/kjk/notion/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/kjk/notion.svg)](https://pkg.go.dev/github.com/kjk/notion)
[![GitHub](https://img.shields.io/github/license/kjk/notion)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/kjk/notion)](https://goreportcard.com/report/github.com/kjk/notion)

**notion** is a Go client for the
 [Notion API](https://developers.notion.com/reference). Based on https://github.com/dstotijn/go-notion (as of https://github.com/dstotijn/go-notion/commit/55aa9db5c7a72af2a57ac953ebbbdbdec3e1efa1, May 19 2021)

## Features

The client supports all (non-deprecated) endpoints available in the Notion API,
as of May 15, 2021:

- [x] [Get database info](https://pkg.go.dev/github.com/kjk/notion#Client.GetDatabase), [example](https://github.com/kjk/notion/blob/master/examples/get_database_info.go)
- [x] [Query a database](https://pkg.go.dev/github.com/kjk/notion#Client.QueryDatabase), [example](https://github.com/kjk/notion/blob/master/examples/query_database.go)
- [x] [Retrieve page info](https://pkg.go.dev/github.com/kjk/notion#Client.GetPage). [example](https://github.com/kjk/notion/blob/master/examples/get_page_info.go)
- [x] [Retrieve children of a block](https://pkg.go.dev/github.com/kjk/notion#Client.GetBlockChildren), [example](https://github.com/kjk/notion/blob/master/examples/get_block_children.go)
- [x] [Create a page](https://pkg.go.dev/github.com/kjk/notion#Client.CreatePage), [example](https://github.com/kjk/notion/blob/master/examples/create_page.go)
- [x] [Update page properties](https://pkg.go.dev/github.com/kjk/notion#Client.UpdatePageProps)
- [x] [Append block children](https://pkg.go.dev/github.com/kjk/notion#Client.AppendBlockChildren)
- [x] [Get user info](https://pkg.go.dev/github.com/kjk/notion#Client.GetUser), [example](https://github.com/kjk/notion/blob/master/examples/get_user.go)
- [x] [List all users](https://pkg.go.dev/github.com/kjk/notion#Client.ListUsers), [example](https://github.com/kjk/notion/blob/master/examples/list_users.go)
- [x] [Search](https://pkg.go.dev/github.com/kjk/notion#Client.Search), [example](https://github.com/kjk/notion/blob/master/examples/search.go)

## Getting started

To obtain an API key (required for API calls), follow Notion’s [getting started guide](https://developers.notion.com/docs/getting-started).

### Code example

First, construct a new `Client`:

```go
import "github.com/kjk/notion"

(...)

client := notion.NewClient("secret-api-key")
```

Then, use the methods defined on `Client` to make requests to the API. For
example:

```go
page, err := client.GetPage(context.Background(), "18d35eb5-91f1-4dcb-85b0-c340fd965015")
if err != nil {
    // Handle error...
}
```

👉 Check out the docs on
[pkg.go.dev](https://pkg.go.dev/github.com/kjk/notion) for further
reference and examples.

## Status

The Notion API is currently in _public beta_.

⚠️ Although the API itself is versioned, this client **will** make breaking
changes in its code until `v1.0` of the module is released.

Official Notion API is still limited:
* not all block types are supported
* no way to avoid re-downloading data we already have

For more functionality use unofficial API client
https://github.com/kjk/notionapi

### To do

- [ ] Write tests

## License

[MIT License](LICENSE)
