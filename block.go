package notion

import "time"

// Block represents content on the Notion platform.
// See: https://developers.notion.com/reference/block
type Block struct {
	Object         string     `json:"object"`
	ID             string     `json:"id,omitempty"`
	Type           BlockType  `json:"type"`
	CreatedTime    *time.Time `json:"created_time,omitempty"`
	LastEditedTime *time.Time `json:"last_edited_time,omitempty"`
	HasChildren    bool       `json:"has_children,omitempty"`

	Paragraph        *RichTextBlock `json:"paragraph,omitempty"`
	Heading1         *Heading       `json:"heading_1,omitempty"`
	Heading2         *Heading       `json:"heading_2,omitempty"`
	Heading3         *Heading       `json:"heading_3,omitempty"`
	BulletedListItem *RichTextBlock `json:"bulleted_list_item,omitempty"`
	NumberedListItem *RichTextBlock `json:"numbered_list_item,omitempty"`
	ToDo             *ToDo          `json:"to_do,omitempty"`
	Toggle           *RichTextBlock `json:"toggle,omitempty"`
	ChildPage        *ChildPage     `json:"rich_text,omitempty"`
}

type RichTextBlock struct {
	Text     []RichText `json:"text"`
	Children []Block    `json:"children,omitempty"`
}

type Heading struct {
	Text []RichText `json:"text"`
}

type ToDo struct {
	RichTextBlock
	Checked *bool `json:"checked,omitempty"`
}

type ChildPage struct {
	Title string `json:"title"`
}

type BlockType string

const (
	BlockTypeParagraph        BlockType = "paragraph"
	BlockTypeHeading1         BlockType = "heading_1"
	BlockTypeHeading2         BlockType = "heading_2"
	BlockTypeHeading3         BlockType = "heading_3"
	BlockTypeBulletedListItem BlockType = "bulleted_list_item"
	BlockTypeNumberedListItem BlockType = "numbered_list_item"
	BlockTypeToDo             BlockType = "to_do"
	BlockTypeToggle           BlockType = "toggle"
	BlockTypeChildPage        BlockType = "child_page"
	BlockTypeUnsupported      BlockType = "unsupported"
)

type PaginationQuery struct {
	StartCursor string
	PageSize    int
}

// BlockChildrenResponse contains results (block children) and pagination data returned from a find request.
type BlockChildrenResponse struct {
	Results    []Block `json:"results"`
	HasMore    bool    `json:"has_more"`
	NextCursor string  `json:"next_cursor"`

	RawJSON []byte `json:"-"`
}
