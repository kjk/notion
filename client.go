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

func (c *Client) newRequest(ctx context.Context, method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, baseURL+url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", c.apiKey))
	req.Header.Set("Notion-Version", apiVersion)
	req.Header.Set("User-Agent", "go-notion/"+clientVersion)

	if method == http.MethodPost || method == http.MethodPatch {
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
func (c *Client) QueryDatabase(ctx context.Context, id string, query DatabaseQuery) (*DatabaseQueryResponse, error) {
	body := &bytes.Buffer{}
	err := json.NewEncoder(body).Encode(query)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to encode filter to JSON: %w", err)
	}

	uri := "/databases/" + id + "/query"
	req, err := c.newRequest(ctx, http.MethodPost, uri, body)
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

	body := &bytes.Buffer{}

	err := json.NewEncoder(body).Encode(params)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to encode body params to JSON: %w", err)
	}

	uri := "/pages"
	req, err := c.newRequest(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	var res Page
	res.RawJSON, err = c.doHTTPAndUnmarshalResponse(req, &res, "create page")
	return &res, err
}

// UpdatePageProps updates page property values for a page.
// See: https://developers.notion.com/reference/patch-page
func (c *Client) UpdatePageProps(ctx context.Context, pageID string, params UpdatePageParams) (page Page, err error) {
	if err := params.Validate(); err != nil {
		return Page{}, fmt.Errorf("notion: invalid page params: %w", err)
	}

	body := &bytes.Buffer{}

	err = json.NewEncoder(body).Encode(params)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to encode body params to JSON: %w", err)
	}

	uri := "/pages/" + pageID
	req, err := c.newRequest(ctx, http.MethodPatch, uri, body)
	if err != nil {
		return Page{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("notion: failed to update page properties: %w", parseErrorResponse(resp))
	}

	err = json.NewDecoder(resp.Body).Decode(&page)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return page, nil
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
func (c *Client) AppendBlockChildren(ctx context.Context, blockID string, children []Block) (block Block, err error) {
	type PostBody struct {
		Children []Block `json:"children"`
	}

	dto := PostBody{children}
	body := &bytes.Buffer{}

	err = json.NewEncoder(body).Encode(dto)
	if err != nil {
		return Block{}, fmt.Errorf("notion: failed to encode body params to JSON: %w", err)
	}

	uri := "/blocks/" + blockID + "/children"
	req, err := c.newRequest(ctx, http.MethodPatch, uri, body)
	if err != nil {
		return Block{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return Block{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Block{}, fmt.Errorf("notion: failed to append block children: %w", parseErrorResponse(resp))
	}

	err = json.NewDecoder(resp.Body).Decode(&block)
	if err != nil {
		return Block{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return block, nil
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}

	var res ListUsersResponse
	res.RawJSON, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if err != nil {
		return &res, err
	}

	if resp.StatusCode != http.StatusOK {
		return &res, fmt.Errorf("notion: failed to list users: %w", parseErrorResponseJSON(res.RawJSON))
	}

	err = json.Unmarshal(res.RawJSON, &res)
	if err != nil {
		return &res, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return &res, nil
}

// Search fetches all pages and child pages that are shared with the integration. Optionally uses query, filter and
// pagination options.
// See: https://developers.notion.com/reference/post-search
func (c *Client) Search(ctx context.Context, opts *SearchOpts) (result SearchResponse, err error) {
	body := &bytes.Buffer{}

	if opts != nil {
		err = json.NewEncoder(body).Encode(opts)
		if err != nil {
			return SearchResponse{}, fmt.Errorf("notion: failed to encode filter to JSON: %w", err)
		}
	}

	uri := "/search"
	req, err := c.newRequest(ctx, http.MethodPost, uri, body)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return SearchResponse{}, fmt.Errorf("notion: failed to search: %w", parseErrorResponse(resp))
	}

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return result, nil
}
