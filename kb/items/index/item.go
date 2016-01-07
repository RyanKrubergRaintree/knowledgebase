package index

import "github.com/raintreeinc/knowledgebase/kb"

type Item struct {
	Title    string  `json:"title"`
	Slug     kb.Slug `json:"slug"`
	Children []*Item `json:"children,omitempty"`
}

func New(id string, item *Item) kb.Item {
	return kb.Item{
		"type": "index",
		"id":   id,
		"root": item,
	}
}
