package dita

import (
	"github.com/raintreeinc/ditaconvert"
	"github.com/raintreeinc/knowledgebase/extra/index"
)

func EntryToIndexItem(entry *ditaconvert.Entry) *index.Item {
	item := &index.Item{
		Title: entry.Title,
	}
	if entry.Topic != nil {
		item.Slug = mapping.ByTopic[entry.Topic]
	}

	for _, child := range entry.Children {
		if !child.TOC {
			continue
		}
		childitem := EntryToIndexItem(mapping, child)
		item.Children = append(item.Children, childitem)
	}

	return item
}
