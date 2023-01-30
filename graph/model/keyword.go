package model

type Keyword struct {
	Value  string `json:"value"`
	Usages int64  `json:"usages"`
}

type SetKeywords struct {
	UUID     string   `json:"uuid"`
	Keywords []string `json:"keywords"`
}
