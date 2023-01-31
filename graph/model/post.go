package model

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	UUID    string `json:"uuid"`
}

func (Post) IsKeywordLink() {}

type NewPost struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	WriterUUID string   `json:"writerUuid"`
	Keywords   []string `json:"keywords"`
}
