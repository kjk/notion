package notion_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kjk/notion"
)

type mockRoundtripper struct {
	fn func(*http.Request) (*http.Response, error)
}

func (m *mockRoundtripper) RoundTrip(r *http.Request) (*http.Response, error) {
	return m.fn(r)
}

func mustParseTime(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func stringPtr(s string) *string {
	return &s
}

func TestFindDatabaseByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expDatabase    *notion.Database
		expError       error
	}{
		{
			name: "successful response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "database",
						"id": "668d797c-76fa-4934-9b05-ad288df2d136",
						"created_time": "2020-03-17T19:10:04.968Z",
						"last_edited_time": "2020-03-17T21:49:37.913Z",
						"title": [
							{
								"type": "text",
								"text": {
									"content": "Grocery List",
									"link": null
								},
								"annotations": {
									"bold": false,
									"italic": false,
									"strikethrough": false,
									"underline": false,
									"code": false,
									"color": "default"
								},
								"plain_text": "Grocery List",
								"href": null
							}
						],
						"properties": {
							"Name": {
								"id": "title",
								"type": "title",
								"title": {}
							},
							"Description": {
								"id": "J@cS",
								"type": "rich_text",
								"text": {}
							},
							"In stock": {
								"id": "{xYx",
								"type": "checkbox",
								"checkbox": {}
							},
							"Food group": {
								"id": "TJmr",
								"type": "select",
								"select": {
									"options": [
										{
											"id": "96eb622f-4b88-4283-919d-ece2fbed3841",
											"name": "🥦Vegetable",
											"color": "green"
										},
										{
											"id": "bb443819-81dc-46fb-882d-ebee6e22c432",
											"name": "🍎Fruit",
											"color": "red"
										},
										{
											"id": "7da9d1b9-8685-472e-9da3-3af57bdb221e",
											"name": "💪Protein",
											"color": "yellow"
										}
									]
								}
							},
							"Price": {
								"id": "cU^N",
								"type": "number",
								"number": {
									"format": "dollar"
								}
							},
							"Cost of next trip": {
								"id": "p:sC",
								"type": "formula",
								"formula": {
									"expression": "if(prop(\"In stock\"), 0, prop(\"Price\"))"
								}
							},
							"Last ordered": {
								"id": "]\\R[",
								"type": "date",
								"date": {}
							},
							"Meals": {
								"id": "lV]M",
								"type": "relation",
								"relation": {
									"database_id": "668d797c-76fa-4934-9b05-ad288df2d136",
									"synced_property_name": "Related to Test database (Relation Test)",
									"synced_property_id": "IJi<"
								}
							},
							"Number of meals": {
								"id": "Z\\Eh",
								"type": "rollup",
								"rollup": {
									"rollup_property_name": "Name",
									"relation_property_name": "Meals",
									"rollup_property_id": "title",
									"relation_property_id": "mxp^",
									"function": "count"
								}
							},
							"Store availability": {
								"id": "=_>D",
								"type": "multi_select",
								"multi_select": {
									"options": [
										{
											"id": "d209b920-212c-4040-9d4a-bdf349dd8b2a",
											"name": "Duc Loi Market",
											"color": "blue"
										},
										{
											"id": "6c3867c5-d542-4f84-b6e9-a420c43094e7",
											"name": "Gus's Community Market",
											"color": "yellow"
										}
									]
								}
							},
							"+1": {
								"id": "aGut",
								"type": "people",
								"people": {}
							},
							"Photo": {
								"id": "aTIT",
								"type": "files",
								"files": {}
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expDatabase: &notion.Database{
				ID:             "668d797c-76fa-4934-9b05-ad288df2d136",
				CreatedTime:    mustParseTime(time.RFC3339, "2020-03-17T19:10:04.968Z"),
				LastEditedTime: mustParseTime(time.RFC3339, "2020-03-17T21:49:37.913Z"),
				Title: []notion.RichText{
					{
						Type: notion.RichTextTypeText,
						Text: &notion.Text{
							Content: "Grocery List",
						},
						Annotations: &notion.Annotations{
							Color: notion.ColorDefault,
						},
						PlainText: "Grocery List",
					},
				},
				Properties: notion.DatabaseProperties{
					"Name": notion.DatabaseProperty{
						ID:   "title",
						Type: notion.DBPropTypeTitle,
					},
					"Description": notion.DatabaseProperty{
						ID:   "J@cS",
						Type: notion.DBPropTypeRichText,
					},
					"In stock": notion.DatabaseProperty{
						ID:   "{xYx",
						Type: notion.DBPropTypeCheckbox,
					},
					"Food group": notion.DatabaseProperty{
						ID:   "TJmr",
						Type: notion.DBPropTypeSelect,
						Select: &notion.SelectMetadata{
							Options: []notion.SelectOptions{
								{
									ID:    "96eb622f-4b88-4283-919d-ece2fbed3841",
									Name:  "🥦Vegetable",
									Color: notion.ColorGreen,
								},
								{
									ID:    "bb443819-81dc-46fb-882d-ebee6e22c432",
									Name:  "🍎Fruit",
									Color: notion.ColorRed,
								},
								{
									ID:    "7da9d1b9-8685-472e-9da3-3af57bdb221e",
									Name:  "💪Protein",
									Color: notion.ColorYellow,
								},
							},
						},
					},
					"Price": notion.DatabaseProperty{
						ID:   "cU^N",
						Type: notion.DBPropTypeNumber,
						Number: &notion.NumberMetadata{
							Format: notion.NumberFormatDollar,
						},
					},
					"Cost of next trip": {
						ID:   "p:sC",
						Type: notion.DBPropTypeFormula,
						Formula: &notion.FormulaMetadata{
							Expression: `if(prop("In stock"), 0, prop("Price"))`,
						},
					},
					"Last ordered": notion.DatabaseProperty{
						ID:   "]\\R[",
						Type: notion.DBPropTypeDate,
					},
					"Meals": notion.DatabaseProperty{
						ID:   "lV]M",
						Type: notion.DBPropTypeRelation,
						Relation: &notion.RelationMetadata{
							DatabaseID:     "668d797c-76fa-4934-9b05-ad288df2d136",
							SyncedPropName: "Related to Test database (Relation Test)",
							SyncedPropID:   "IJi<",
						},
					},
					"Number of meals": notion.DatabaseProperty{
						ID:   "Z\\Eh",
						Type: notion.DBPropTypeRollup,
						Rollup: &notion.RollupMetadata{
							RollupPropName:   "Name",
							RelationPropName: "Meals",
							RollupPropID:     "title",
							RelationPropID:   "mxp^",
							Function:         "count",
						},
					},
					"Store availability": notion.DatabaseProperty{
						ID:   "=_>D",
						Type: notion.DBPropTypeMultiSelect,
						MultiSelect: &notion.SelectMetadata{
							Options: []notion.SelectOptions{
								{
									ID:    "d209b920-212c-4040-9d4a-bdf349dd8b2a",
									Name:  "Duc Loi Market",
									Color: notion.ColorBlue,
								},
								{
									ID:    "6c3867c5-d542-4f84-b6e9-a420c43094e7",
									Name:  "Gus's Community Market",
									Color: notion.ColorYellow,
								},
							},
						},
					},
					"+1": notion.DatabaseProperty{
						ID:   "aGut",
						Type: notion.DBPropTypePeople,
					},
					"Photo": {
						ID:   "aTIT",
						Type: "files",
					},
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "error",
						"status": 400,
						"code": "validation_error",
						"message": "foobar"
					}`,
				)
			},
			respStatusCode: http.StatusBadRequest,
			expDatabase:    &notion.Database{},
			expError:       errors.New("notion: failed to find database: foobar (code: validation_error, status: 400)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			opts := notion.ClientOptions{
				HTTPClient: httpClient,
			}
			client := notion.NewClient("secret-api-key", &opts)
			db, err := client.GetDatabase(context.Background(), "00000000-0000-0000-0000-000000000000")
			db.RawJSON = nil

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expDatabase, db); diff != "" {
				t.Fatalf("database not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestQueryDatabase(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		query          *notion.DatabaseQuery
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    *notion.DatabaseQueryResponse
		expError       error
	}{
		{
			name: "with query, successful response",
			query: &notion.DatabaseQuery{
				Filter: &notion.DatabaseQueryFilter{
					Property: "Name",
					Text: &notion.TextDatabaseQueryFilter{
						Contains: "foobar",
					},
				},
				Sorts: []notion.DatabaseQuerySort{
					{
						Property:  "Name",
						Timestamp: notion.SortTimeStampCreatedTime,
						Direction: notion.SortDirAsc,
					},
					{
						Property:  "Date",
						Timestamp: notion.SortTimeStampLastEditedTime,
						Direction: notion.SortDirDesc,
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [
							{
								"object": "page",
								"id": "7c6b1c95-de50-45ca-94e6-af1d9fd295ab",
								"created_time": "2021-05-18T17:50:22.371Z",
								"last_edited_time": "2021-05-18T17:50:22.371Z",
								"parent": {
									"type": "database_id",
									"database_id": "39ddfc9d-33c9-404c-89cf-79f01c42dd0c"
								},
								"archived": false,
								"properties": {
									"Date": {
										"id": "Q]uT",
										"type": "date",
										"date": {
											"start": "2021-05-18T12:49:00.000-05:00",
											"end": null
										}
									},
									"Name": {
										"id": "title",
										"type": "title",
										"title": [
											{
												"type": "text",
												"text": {
													"content": "Foobar",
													"link": null
												},
												"annotations": {
													"bold": false,
													"italic": false,
													"strikethrough": false,
													"underline": false,
													"code": false,
													"color": "default"
												},
												"plain_text": "Foobar",
												"href": null
											}
										]
									}
								}
							}
						],
						"next_cursor": "A^hd",
						"has_more": true
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"filter": map[string]interface{}{
					"property": "Name",
					"text": map[string]interface{}{
						"contains": "foobar",
					},
				},
				"sorts": []interface{}{
					map[string]interface{}{
						"property":  "Name",
						"timestamp": "created_time",
						"direction": "ascending",
					},
					map[string]interface{}{
						"property":  "Date",
						"timestamp": "last_edited_time",
						"direction": "descending",
					},
				},
			},
			expResponse: &notion.DatabaseQueryResponse{
				Results: []notion.Page{
					{
						ID:             "7c6b1c95-de50-45ca-94e6-af1d9fd295ab",
						CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-18T17:50:22.371Z"),
						LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-18T17:50:22.371Z"),
						Parent: notion.PageParent{
							Type:       notion.ParentTypeDatabase,
							DatabaseID: stringPtr("39ddfc9d-33c9-404c-89cf-79f01c42dd0c"),
						},
						Archived: false,
						Properties: notion.DatabasePageProperties{
							"Date": notion.DatabasePageProperty{
								ID:   "Q]uT",
								Type: notion.DBPropTypeDate,
								Date: &notion.Date{
									Start: notion.Time(mustParseTime(time.RFC3339Nano, "2021-05-18T12:49:00.000-05:00")),
								},
							},
							"Name": notion.DatabasePageProperty{
								ID:   "title",
								Type: notion.DBPropTypeTitle,
								Title: []notion.RichText{
									{
										Type: notion.RichTextTypeText,
										Text: &notion.Text{
											Content: "Foobar",
										},
										PlainText: "Foobar",
										Annotations: &notion.Annotations{
											Color: notion.ColorDefault,
										},
									},
								},
							},
						},
					},
				},
				HasMore:    true,
				NextCursor: "A^hd",
			},
			expError: nil,
		},
		{
			name:  "without query, successful response",
			query: nil,
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "list",
						"results": [],
						"next_cursor": null,
						"has_more": false
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody:    nil,
			expResponse: &notion.DatabaseQueryResponse{
				Results:    []notion.Page{},
				HasMore:    false,
				NextCursor: "",
			},
			expError: nil,
		},
		{
			name: "error response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "error",
						"status": 400,
						"code": "validation_error",
						"message": "foobar"
					}`,
				)
			},
			respStatusCode: http.StatusBadRequest,
			expResponse:    &notion.DatabaseQueryResponse{},
			expError:       errors.New("notion: failed to query database: foobar (code: validation_error, status: 400)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					postBody := make(map[string]interface{})

					err := json.NewDecoder(r.Body).Decode(&postBody)
					if err != nil && err != io.EOF {
						t.Fatal(err)
					}

					if len(tt.expPostBody) == 0 && len(postBody) != 0 {
						t.Errorf("unexpected post body: %+v", postBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) == 0 {
						t.Errorf("post body not equal (expected %+v, got: nil)", tt.expPostBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) != 0 {
						if diff := cmp.Diff(tt.expPostBody, postBody); diff != "" {
							t.Errorf("post body not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			opts := notion.ClientOptions{
				HTTPClient: httpClient,
			}
			client := notion.NewClient("secret-api-key", &opts)
			resp, err := client.QueryDatabase(context.Background(), "00000000-0000-0000-0000-000000000000", tt.query)
			if resp != nil {
				resp.RawJSON = nil
			}

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, resp); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestFindPageByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPage        *notion.Page
		expError       error
	}{
		{
			name: "successful response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "page",
						"id": "606ed832-7d79-46de-bbed-5b4896e7bc02",
						"created_time": "2021-05-19T18:34:00.000Z",
						"last_edited_time": "2021-05-19T18:34:00.000Z",
						"parent": {
							"type": "page_id",
							"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
						},
						"archived": false,
						"properties": {
							"title": {
								"id": "title",
								"type": "title",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Lorem ipsum",
											"link": null
										},
										"annotations": {
											"bold": false,
											"italic": false,
											"strikethrough": false,
											"underline": false,
											"code": false,
											"color": "default"
										},
										"plain_text": "Lorem ipsum",
										"href": null
									}
								]
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPage: &notion.Page{
				ID:             "606ed832-7d79-46de-bbed-5b4896e7bc02",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-19T18:34:00.000Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-19T18:34:00.000Z"),
				Parent: notion.PageParent{
					Type:   notion.ParentTypePage,
					PageID: stringPtr("b0668f48-8d66-4733-9bdb-2f82215707f7"),
				},
				Properties: notion.PageProperties{
					Title: notion.PageTitle{
						Title: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: "Lorem ipsum",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
								PlainText: "Lorem ipsum",
							},
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "error",
						"status": 404,
						"code": "object_not_found",
						"message": "foobar"
					}`,
				)
			},
			respStatusCode: http.StatusNotFound,
			expPage:        &notion.Page{},
			expError:       errors.New("notion: failed to find page: foobar (code: object_not_found, status: 404)"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			opts := notion.ClientOptions{
				HTTPClient: httpClient,
			}
			client := notion.NewClient("secret-api-key", &opts)
			page, err := client.GetPage(context.Background(), "00000000-0000-0000-0000-000000000000")
			page.RawJSON = nil

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expPage, page); diff != "" {
				t.Fatalf("page not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}

func TestCreatePage(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		params         notion.CreatePageParams
		respBody       func(r *http.Request) io.Reader
		respStatusCode int
		expPostBody    map[string]interface{}
		expResponse    *notion.Page
		expError       error
	}{
		{
			name: "successful response",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypePage,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
				Children: []notion.Block{
					{
						Object: "block",
						Type:   notion.BlockTypeParagraph,
						Paragraph: &notion.RichTextBlock{
							Text: []notion.RichText{
								{
									Text: &notion.Text{
										Content: "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "page",
						"id": "276ee233-e426-4ed0-9986-6b22af8550df",
						"created_time": "2021-05-19T19:34:05.068Z",
						"last_edited_time": "2021-05-19T19:34:05.069Z",
						"parent": {
							"type": "page_id",
							"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7"
						},
						"archived": false,
						"properties": {
							"title": {
								"id": "title",
								"type": "title",
								"title": [
									{
										"type": "text",
										"text": {
											"content": "Foobar",
											"link": null
										},
										"annotations": {
											"bold": false,
											"italic": false,
											"strikethrough": false,
											"underline": false,
											"code": false,
											"color": "default"
										},
										"plain_text": "Foobar",
										"href": null
									}
								]
							}
						}
					}`,
				)
			},
			respStatusCode: http.StatusOK,
			expPostBody: map[string]interface{}{
				"parent": map[string]interface{}{
					"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				"properties": map[string]interface{}{
					"title": []interface{}{
						map[string]interface{}{
							"text": map[string]interface{}{
								"content": "Foobar",
							},
						},
					},
				},
				"children": []interface{}{
					map[string]interface{}{
						"object": "block",
						"type":   "paragraph",
						"paragraph": map[string]interface{}{
							"text": []interface{}{
								map[string]interface{}{
									"text": map[string]interface{}{
										"content": "Lorem ipsum dolor sit amet.",
									},
								},
							},
						},
					},
				},
			},
			expResponse: &notion.Page{
				ID:             "276ee233-e426-4ed0-9986-6b22af8550df",
				CreatedTime:    mustParseTime(time.RFC3339Nano, "2021-05-19T19:34:05.068Z"),
				LastEditedTime: mustParseTime(time.RFC3339Nano, "2021-05-19T19:34:05.069Z"),
				Parent: notion.PageParent{
					Type:   notion.ParentTypePage,
					PageID: stringPtr("b0668f48-8d66-4733-9bdb-2f82215707f7"),
				},
				Properties: notion.PageProperties{
					Title: notion.PageTitle{
						Title: []notion.RichText{
							{
								Type: notion.RichTextTypeText,
								Text: &notion.Text{
									Content: "Foobar",
								},
								Annotations: &notion.Annotations{
									Color: notion.ColorDefault,
								},
								PlainText: "Foobar",
							},
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "error response",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypePage,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
			},
			respBody: func(_ *http.Request) io.Reader {
				return strings.NewReader(
					`{
						"object": "error",
						"status": 400,
						"code": "validation_error",
						"message": "foobar"
					}`,
				)
			},
			respStatusCode: http.StatusBadRequest,
			expPostBody: map[string]interface{}{
				"parent": map[string]interface{}{
					"page_id": "b0668f48-8d66-4733-9bdb-2f82215707f7",
				},
				"properties": map[string]interface{}{
					"title": []interface{}{
						map[string]interface{}{
							"text": map[string]interface{}{
								"content": "Foobar",
							},
						},
					},
				},
			},
			expResponse: &notion.Page{},
			expError:    errors.New("notion: failed to create page: foobar (code: validation_error, status: 400)"),
		},
		{
			name: "parent type required error",
			params: notion.CreatePageParams{
				ParentID: "b0668f48-8d66-4733-9bdb-2f82215707f7",
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
			},
			expResponse: nil,
			expError:    errors.New("notion: invalid page params: parent type is required"),
		},
		{
			name: "parent id required error",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypePage,
				Title: []notion.RichText{
					{
						Text: &notion.Text{
							Content: "Foobar",
						},
					},
				},
			},
			expResponse: nil,
			expError:    errors.New("notion: invalid page params: parent ID is required"),
		},
		{
			name: "page title required error",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypePage,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
			},
			expResponse: nil,
			expError:    errors.New("notion: invalid page params: title is required when parent type is page"),
		},
		{
			name: "database properties required error",
			params: notion.CreatePageParams{
				ParentType: notion.ParentTypeDatabase,
				ParentID:   "b0668f48-8d66-4733-9bdb-2f82215707f7",
			},
			expResponse: nil,
			expError:    errors.New("notion: invalid page params: database page properties is required when parent type is database"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			httpClient := &http.Client{
				Transport: &mockRoundtripper{fn: func(r *http.Request) (*http.Response, error) {
					postBody := make(map[string]interface{})

					err := json.NewDecoder(r.Body).Decode(&postBody)
					if err != nil && err != io.EOF {
						t.Fatal(err)
					}

					if len(tt.expPostBody) == 0 && len(postBody) != 0 {
						t.Errorf("unexpected post body: %#v", postBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) == 0 {
						t.Errorf("post body not equal (expected %+v, got: nil)", tt.expPostBody)
					}

					if len(tt.expPostBody) != 0 && len(postBody) != 0 {
						if diff := cmp.Diff(tt.expPostBody, postBody); diff != "" {
							t.Errorf("post body not equal (-exp, +got):\n%v", diff)
						}
					}

					return &http.Response{
						StatusCode: tt.respStatusCode,
						Status:     http.StatusText(tt.respStatusCode),
						Body:       ioutil.NopCloser(tt.respBody(r)),
					}, nil
				}},
			}
			opts := notion.ClientOptions{
				HTTPClient: httpClient,
			}
			client := notion.NewClient("secret-api-key", &opts)
			page, err := client.CreatePage(context.Background(), tt.params)
			if page != nil {
				page.RawJSON = nil
			}

			if tt.expError == nil && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if tt.expError != nil && err == nil {
				t.Fatalf("error not equal (expected: %v, got: nil)", tt.expError)
			}
			if tt.expError != nil && err != nil && tt.expError.Error() != err.Error() {
				t.Fatalf("error not equal (expected: %v, got: %v)", tt.expError, err)
			}

			if diff := cmp.Diff(tt.expResponse, page); diff != "" {
				t.Fatalf("response not equal (-exp, +got):\n%v", diff)
			}
		})
	}
}
