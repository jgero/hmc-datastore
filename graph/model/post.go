package model

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Uuid    string `json:"uuid"`
}

type NewPost struct {
	Title      string   `json:"title"`
	Content    string   `json:"content"`
	WriterUUID string   `json:"writerUuid"`
	Keywords   []string `json:"keywords"`
}
