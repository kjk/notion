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

// ClientOption is used to override default client behavior.
type ClientOption func(*Client)

// NewClient returns a new Client.
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:     apiKey,
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithHTTPClient overrides the default http.Client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
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

// FindDatabaseByID fetches a database by ID.
// See: https://developers.notion.com/reference/get-database
func (c *Client) FindDatabaseByID(ctx context.Context, id string) (db Database, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/databases/"+id, nil)
	if err != nil {
		return Database{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Database{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Database{}, fmt.Errorf("notion: failed to find database: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&db)
	if err != nil {
		return Database{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return db, nil
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

// FindPageByID fetches a page by ID.
// See: https://developers.notion.com/reference/get-page
func (c *Client) FindPageByID(ctx context.Context, id string) (page Page, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/pages/"+id, nil)
	if err != nil {
		return Page{}, fmt.Errorf("notion: invalid request: %w", err)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Page{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}

	d, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("notion: failed to find page: %w", parseErrorResponse(res))
	}
	//fmt.Printf("Body:\n%s\n", string(d))

	page.RawJSON = d
	if err != nil {
		return page, err
	}

	err = json.Unmarshal(d, &page)
	page.RawJSON = d
	if err != nil {
		return page, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return page, nil
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

	d, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("notion: failed to create page: %w", parseErrorResponse(res))
	}

	page.RawJSON = d
	if err != nil {
		return page, err
	}

	err = json.Unmarshal(d, &page)
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

// FindBlockChildrenByID returns a list of block children for a given block ID.
// See: https://developers.notion.com/reference/post-database-query
func (c *Client) FindBlockChildrenByID(ctx context.Context, blockID string, query *PaginationQuery) (result BlockChildrenResponse, err error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/blocks/%v/children", blockID), nil)
	if err != nil {
		return BlockChildrenResponse{}, fmt.Errorf("notion: invalid request: %w", err)
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
		return BlockChildrenResponse{}, fmt.Errorf("notion: failed to make HTTP request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return BlockChildrenResponse{}, fmt.Errorf("notion: failed to find block children: %w", parseErrorResponse(res))
	}

	err = json.NewDecoder(res.Body).Decode(&result)
	if err != nil {
		return BlockChildrenResponse{}, fmt.Errorf("notion: failed to parse HTTP response: %w", err)
	}

	return result, nil
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
