package notion

type RichText struct {
	Type        RichTextType `json:"type,omitempty"`
	Annotations *Annotations `json:"annotations,omitempty"`

	PlainText string    `json:"plain_text,omitempty"`
	HRef      *string   `json:"href,omitempty"`
	Text      *Text     `json:"text,omitempty"`
	Mention   *Mention  `json:"mention,omitempty"`
	Equation  *Equation `json:"equation,omitempty"`
}

type Equation struct {
	Expression string `json:"expression"`
}

type Annotations struct {
	Bold          bool  `json:"bold,omitempty"`
	Italic        bool  `json:"italic,omitempty"`
	Strikethrough bool  `json:"strikethrough,omitempty"`
	Underline     bool  `json:"underline,omitempty"`
	Code          bool  `json:"code,omitempty"`
	Color         Color `json:"color,omitempty"`
}

type Mention struct {
	Type MentionType `json:"type"`

	User     *User `json:"user,omitempty"`
	Page     *ID   `json:"page,omitempty"`
	Database *ID   `json:"database,omitempty"`
	Date     *Date `json:"date,omitempty"`
}

type Date struct {
	Start Time  `json:"start"`
	End   *Time `json:"end,omitempty"`
}

type Text struct {
	Content string `json:"content"`
	Link    *Link  `json:"link,omitempty"`
}

type Link struct {
	URL string `json:"url"`
}

type ID struct {
	ID string `json:"id"`
}

type (
	RichTextType string
	MentionType  string
	Color        string
)

const (
	RichTextTypeText     RichTextType = "text"
	RichTextTypeMention  RichTextType = "mention"
	RichTextTypeEquation RichTextType = "equation"
)

const (
	MentionTypeUser     MentionType = "user"
	MentionTypePage     MentionType = "page"
	MentionTypeDatabase MentionType = "database"
	MentionTypeDate     MentionType = "date"
)

const (
	ColorDefault  Color = "default"
	ColorGray     Color = "gray"
	ColorBrown    Color = "brown"
	ColorOrange   Color = "orange"
	ColorYellow   Color = "yellow"
	ColorGreen    Color = "green"
	ColorBlue     Color = "blue"
	ColorPurple   Color = "purple"
	ColorPink     Color = "pink"
	ColorRed      Color = "red"
	ColorGrayBg   Color = "gray_background"
	ColorBrownBg  Color = "brown_background"
	ColorOrangeBg Color = "orange_background"
	ColorYellowBg Color = "yellow_background"
	ColorGreenBg  Color = "green_background"
	ColorBlueBg   Color = "blue_background"
	ColorPurpleBg Color = "purple_background"
	ColorPinkBg   Color = "pink_background"
	ColorRedBg    Color = "red_background"
)
