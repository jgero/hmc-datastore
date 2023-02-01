package model

type NewPerson struct {
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
}

type Person struct {
	Name        string     `json:"name"`
	UUID        string     `json:"uuid"`
	Keywords    []*Keyword `json:"keywords"`
	Created     int64      `json:"created"`
	Updated     int64      `json:"updated"`
	UpdateCount int64      `json:"updateCount"`
}

func (Person) IsKeywordLink() {}
