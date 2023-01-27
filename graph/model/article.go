package model

type Article struct {
    Title string `json:"title"`
    Content string `json:"content"`
    Writer *Person `json:"writer"`
}
