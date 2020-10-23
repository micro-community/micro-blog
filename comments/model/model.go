package model

//NewComment Return Comments
func NewComment(opts ...Option) *Comment {
	options := &Options{}
	for _, o := range opts {
		o(options)
	}
	return &Comment{
		Options: options,
	}
}

//Comment for post
type Comment struct {
	ID              string   `json:"id"`
	ResourceID      string   `json:"resource_id"` //post id of a article
	Index           int32    `json:"index"`
	Content         string   `json:"content"`
	CreateTimestamp int64    `json:"create_timestamp"`
	UpdateTimestamp int64    `json:"update_timestamp"`
	Options         *Options `json:"-"`
}
