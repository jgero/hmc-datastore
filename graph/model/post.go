package model

type Post struct {
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Uuid    string  `json:"uuid"`
}
