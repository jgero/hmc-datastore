package model

type Keyword struct {
	Value  string `json:"value"`
	Usages int64  `json:"usages"`
}

type SetKeywords struct {
	UUIDs     []string `json:"uuids"`
	Keywords  []string `json:"keywords"`
	Exclusive bool     `json:"exclusive"`
}

type KeywordLink interface {
	IsKeywordLink()
}
