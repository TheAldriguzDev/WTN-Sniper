package notifier

type Discord_Params struct {
	Webhook_url  string
	Webhook_data Webhook
	Retries      int
}

type Webhook struct {
	Content     interface{}   `json:"content,omitempty"`
	Embeds      []Embed       `json:"embeds,omitempty"`
	Attachments []interface{} `json:"attachments,omitempty"`
}

type Embed struct {
	Title       string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	URL         string  `json:"url,omitempty"`
	Color       int64   `json:"color,omitempty"`
	Fields      []Field `json:"fields,omitempty"`
	Author      Author  `json:"author,omitempty"`
	Footer      Footer  `json:"footer,omitempty"`
	Timestamp   string  `json:"timestamp,omitempty"`
	Image       Image   `json:"image,omitempty"`
	Thumbnail   Image   `json:"thumbnail,omitempty"`
}

type Author struct {
	Name    string `json:"name,omitempty"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type Field struct {
	Name   string `json:"name,omitempty"`
	Value  string `json:"value,omitempty"`
	Inline bool   `json:"inline,omitempty"`
}

type Footer struct {
	Text    string `json:"text,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

type Image struct {
	URL string `json:"url,omitempty"`
}
