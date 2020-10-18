package model

// Some Const data for tag
const (
	idPrefix        = "id"
	timeStampPrefix = "timestamp"
	slugPrefix      = "bySlug"
	resourcePrefix  = "byResource"
	typePrefix      = "byType"
	tagCountPrefix  = "tagCount"
	childrenByTag   = "childrenByTag"
)

//DB to handle DB
type DB struct {
}

//NewService return a model context
func NewService() *DB {
	return &DB{}
}

//Tag for article
type Tag struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Type            string `json:"type"`
	Count           int64  `json:"count"`
	CreateTimestamp int64  `json:"create_timestamp"`
	UpdateTimestamp int64  `json:"update_timestamp"`
}
