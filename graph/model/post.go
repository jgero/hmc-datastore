package model

type Post struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	UUID        string `json:"uuid"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
	UpdateCount int64  `json:"updateCount"`
}

func (Post) IsKeywordLink() {}

type NewPost struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	WriterUUID string   `json:"writerUuid"`
	Keywords   []string `json:"keywords"`
}

type UpdatePost struct {
	UUID    string  `json:"uuid"`
	Title   *string `json:"title"`
	Content *string `json:"content"`
}
