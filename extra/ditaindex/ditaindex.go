package ditaindex

import (
	"github.com/raintreeinc/knowledgebase/ditaconv"
	"github.com/raintreeinc/knowledgebase/kb"
)

type Item struct {
	Title    string  `json:"title"`
	Slug     kb.Slug `json:"slug"`
	Children []*Item `json:"children,omitempty"`
}

func New(id string, item *Item) kb.Item {
	return kb.Item{
		"type": "dita-index",
		"id":   id,
		"root": item,
	}
}

func EntryToItem(mapping *ditaconv.Mapping, entry *ditaconv.Entry) *Item {
	item := &Item{
		Title: entry.Title,
	}
	if entry.Topic != nil {
		item.Slug = mapping.ByTopic[entry.Topic]
	}

	for _, child := range entry.Children {
		if !child.TOC {
			continue
		}
		childitem := EntryToItem(mapping, child)
		item.Children = append(item.Children, childitem)
	}

	return item
}
