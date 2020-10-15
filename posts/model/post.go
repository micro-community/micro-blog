package model

//Post for article
type Post struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Content         string   `json:"content"`
	CreateTimestamp int64    `json:"create_timestamp"`
	UpdateTimestamp int64    `json:"update_timestamp"`
	TagNames        []string `json:"tagNames"`
}
