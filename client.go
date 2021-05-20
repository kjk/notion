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

// GetDatabase fetches a database by ID.
// See: https://developers.notion.com/reference/get-database
func (c *Client) GetDatabase(ctx context.Context, id string) (*Database, error) {

	req, err := c.newRequest(ctx, http.MethodGet, "/databases/"+id, nil)
	if err != nil {
		return nil, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	var db Database

	db.RawJSON, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return &db, err
	}

	if res.StatusCode != http.StatusOK {
		return &db, fmt.Errorf("notion: failed to find database: %w", parseErrorResponseJSON(db.RawJSON))
	}

	err = json.Unmarshal(db.RawJSON, &db)
	if err != nil {
		return &db, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return &db, nil
}

// QueryDatabase returns database contents, with optional filters, sorts and pagination.
// See: https://developers.notion.com/reference/post-database-query
func (c *Client) QueryDatabase(ctx context.Context, id string, query DatabaseQuery) (result DatabaseQueryResponse, err error) {
	body := &bytes.Buffer{}

	err = json.NewEncoder(body).Encode(query)
	if err != nil {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: failed to encode filter to JSON: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, fmt.Sprintf("/databases/%v/query", id), body)
	if err != nil {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: failed to find database: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return DatabaseQueryResponse{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return result, nil
}

// GetPage fetches information about a page by ID
// See: https://developers.notion.com/reference/get-page
func (c *Client) GetPage(ctx context.Context, id string) (*Page, error) {
	var p Page
	req, err := c.newRequest(ctx, http.MethodGet, "/pages/"+id, nil)
	if err != nil {
		return &p, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}

	p.RawJSON, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return &p, err
	}

	if res.StatusCode != http.StatusOK {
		return &p, fmt.Errorf("notion: failed to find page: %w", parseErrorResponseJSON(p.RawJSON))
	}

	err = json.Unmarshal(p.RawJSON, &p)
	if err != nil {
		return &p, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return &p, nil
}

// CreatePage creates a new page in the specified database or as a child of an existing page.
// See: https://developers.notion.com/reference/post-page
func (c *Client) CreatePage(ctx context.Context, params CreatePageParams) (page Page, err error) {
	if err := params.Validate(); err != nil {
		return Page{}, fmt.Errorf("notion: invalid page params: %w", err)
	}

	body := &bytes.Buffer{}

	err = json.NewEncoder(body).Encode(params)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to encode body params to JSON: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/pages", body)
	if err != nil {
		return Page{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}

	page.RawJSON, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return page, err
	}

	if res.StatusCode != http.StatusOK {
		return page, fmt.Errorf("notion: failed to create page: %w", parseErrorResponseJSON(page.RawJSON))
	}

	err = json.Unmarshal(page.RawJSON, &page)
	if err != nil {
		return page, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return page, nil
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

	req, err := c.newRequest(ctx, http.MethodPatch, "/pages/"+pageID, body)
	if err != nil {
		return Page{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("notion: failed to update page properties: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&page)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return page, nil
}

// GetBlockChildren returns a list of block children for a given block ID.
// See: https://developers.notion.com/reference/get-block-children
func (c *Client) GetBlockChildren(ctx context.Context, blockID string, query *PaginationQuery) (*BlockChildrenResponse, error) {
	var bcr BlockChildrenResponse
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/blocks/%v/children", blockID), nil)
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

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}

	bcr.RawJSON, err = ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		return &bcr, err
	}

	if res.StatusCode != http.StatusOK {
		return &bcr, fmt.Errorf("notion: failed to find block children: %w", parseErrorResponseJSON(bcr.RawJSON))
	}

	err = json.Unmarshal(bcr.RawJSON, &bcr)
	if err != nil {
		return &bcr, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return &bcr, nil
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

	req, err := c.newRequest(ctx, http.MethodPatch, fmt.Sprintf("/blocks/%v/children", blockID), body)
	if err != nil {
		return Block{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Block{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Block{}, fmt.Errorf("notion: failed to append block children: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&block)
	if err != nil {
		return Block{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return block, nil
}

// FindUserByID fetches a user by ID.
// See: https://developers.notion.com/reference/get-user
func (c *Client) FindUserByID(ctx context.Context, id string) (user User, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/users/"+id, nil)
	if err != nil {
		return User{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return User{}, fmt.Errorf("notion: failed to find user: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&user)
	if err != nil {
		return User{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return user, nil
}

// ListUsers returns a list of all users, and pagination metadata.
// See: https://developers.notion.com/reference/get-users
func (c *Client) ListUsers(ctx context.Context, query *PaginationQuery) (result ListUsersResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/users", nil)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("notion: invalid request: %w", err)
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

	res, err := c.httpClient.Do(req)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return ListUsersResponse{}, fmt.Errorf("notion: failed to list users: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return ListUsersResponse{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return result, nil
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

	req, err := c.newRequest(ctx, http.MethodPost, "/search", body)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return SearchResponse{}, fmt.Errorf("notion: failed to search: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return SearchResponse{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return result, nil
}
