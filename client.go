package notion

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

const (
	baseURL       = "https://api.notion.com/v1"
	apiVersion    = "2021-05-13"
	clientVersion = "0.0.0"
)

// Client is used for HTTP requests to the Notion API.
type Client struct {
	apiKey     string
	httpClient *http.Client
}

// ClientOptions describes options when creating client
type ClientOptions struct {
	HTTPClient *http.Client
}

// NewClient returns a new Client.
func NewClient(apiKey string, opts *ClientOptions) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}

	if opts != nil {
		if opts.HTTPClient != nil {
			c.httpClient = opts.HTTPClient
		}
	}

	return c
}

func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	return reflect.ValueOf(v).Kind() == reflect.Ptr && reflect.ValueOf(v).IsNil()
}

func (c *Client) newRequestJSON(ctx context.Context, method, url string, params interface{}) (*http.Request, error) {
	if isNil(params) {
		return c.newRequest(ctx, method, url, nil)
	}

	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to encode body params to JSON: %w", err)
	}
	return c.newRequest(ctx, method, url, body)
}

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, baseURL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiKey))
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("User-Agent", "go-notion/"+clientVersion)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (c *Client) doHTTPAndUnmarshalResponse(req *http.Request, val interface{}, op string) ([]byte, error) {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}

	d, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return d, err
	}

	if resp.StatusCode != http.StatusOK {
		return d, fmt.Errorf("notion: failed to '%s': %w", op, parseErrorResponseJSON(d))
	}

	err = json.Unmarshal(d, val)
	if err != nil {
		return d, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}
	return d, nil
}

// GetDatabase fetches information about a database given its ID.
// See: https://developers.notion.com/reference/get-database
func (c *Client) GetDatabase(ctx context.Context, id string) (*Database, error) {

	uri := "/databases/" + id
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	var res Database
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "find database")
	return &res, err
}

// QueryDatabase returns database contents, with optional filters, sorts and pagination.
// See: https://developers.notion.com/reference/post-database-query
func (c *Client) QueryDatabase(ctx context.Context, id string, query *DatabaseQuery) (*DatabaseQueryResponse, error) {
	uri := "/databases/" + id + "/query"
	req, err := c.newRequestJSON(ctx, http.MethodPost, uri, query)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	var res DatabaseQueryResponse
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "query database")
	return &res, err
}

// GetPage fetches information about a page by ID
// See: https://developers.notion.com/reference/get-page
func (c *Client) GetPage(ctx context.Context, id string) (*Page, error) {
	uri := "/pages/" + id
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	var res Page
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "find page")
	return &res, err
}

// CreatePage creates a new page in the specified database or as a child of an existing page.
// See: https://developers.notion.com/reference/post-page
func (c *Client) CreatePage(ctx context.Context, params CreatePageParams) (*Page, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("notion: invalid page params: %w", err)
	}

	uri := "/pages"
	req, err := c.newRequestJSON(ctx, http.MethodPost, uri, params)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	var res Page
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "create page")
	return &res, err
}

// UpdatePageProps updates page property values for a page.
// See: https://developers.notion.com/reference/patch-page
func (c *Client) UpdatePageProps(ctx context.Context, pageID string, params UpdatePageParams) (*Page, error) {
	if err := params.Validate(); err != nil {
		return nil, fmt.Errorf("notion: invalid page params: %w", err)
	}

	uri := "/pages/" + pageID
	req, err := c.newRequestJSON(ctx, http.MethodPatch, uri, params)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	var res Page
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "update page properties")
	return &res, err
}

// GetBlockChildren returns a list of block children for a given block ID.
// See: https://developers.notion.com/reference/get-block-children
func (c *Client) GetBlockChildren(ctx context.Context, blockID string, query *PaginationQuery) (*BlockChildrenResponse, error) {
	uri := "/blocks/" + blockID + "/children"
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	if query != nil {
		q := url.Values{}
		if query.StartCursor != "" {
			q.Set("start_cursor", query.StartCursor)
		}
		if query.PageSize != 0 {
			q.Set("page_size", strconv.Itoa(query.PageSize))
		}
		req.URL.RawQuery = q.Encode()
	}

	var res BlockChildrenResponse
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "find block children")
	return &res, err
}

// AppendBlockChildren appends child content (blocks) to an existing block.
// See: https://developers.notion.com/reference/patch-block-children
func (c *Client) AppendBlockChildren(ctx context.Context, blockID string, children []Block) (*Block, error) {
	type PostBody struct {
		Children []Block `json:"children"`
	}
	dto := PostBody{children}
	uri := "/blocks/" + blockID + "/children"
	req, err := c.newRequestJSON(ctx, http.MethodPatch, uri, dto)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}
	var res Block
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "append block children")
	return &res, err
}

// FindUserByID fetches a user by ID.
// See: https://developers.notion.com/reference/get-user
func (c *Client) GetUser(ctx context.Context, id string) (*User, error) {

	uri := "/users/" + id
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	var res User
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "find user")
	return &res, err
}

// ListUsers returns a list of all users, and pagination metadata.
// See: https://developers.notion.com/reference/get-users
func (c *Client) ListUsers(ctx context.Context, query *PaginationQuery) (*ListUsersResponse, error) {

	uri := "/users"
	req, err := c.newRequest(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	if query != nil {
		q := url.Values{}
		if query.StartCursor != "" {
			q.Set("start_cursor", query.StartCursor)
		}
		pageSize := query.PageSize
		if pageSize > 0 {
			// limit to max allowed
			if pageSize > 100 {
				pageSize = 100
			}
			q.Set("page_size", strconv.Itoa(pageSize))
		}
		req.URL.RawQuery = q.Encode()
	}

	var res ListUsersResponse
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "list users")
	return &res, err
}

// Search fetches all pages and child pages that are shared with the integration. Optionally uses query, filter and
// pagination options.
// See: https://developers.notion.com/reference/post-search
func (c *Client) Search(ctx context.Context, opts *SearchOpts) (*SearchResponse, error) {
	uri := "/search"
	req, err := c.newRequestJSON(ctx, http.MethodPost, uri, opts)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}
	var res SearchResponse
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "search")
	return &res, err
}
