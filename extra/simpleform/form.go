package simpleform

import "github.com/raintreeinc/knowledgebase/kb"

func Field(id, label string) Item {
	return Item{
		"type":  "field",
		"id":    id,
		"label": label,
	}
}

func Option(id string, values []string) Item {
	return Item{
		"type":   "option",
		"id":     id,
		"values": values,
	}
}

func Button(action, caption string) Item {
	return Item{
		"type":    "button",
		"action":  action,
		"caption": caption,
	}
}

type Item map[string]interface{}

func New(url, text string, items ...Item) kb.Item {
	return kb.Item{
		"type":  "simple-form",
		"id":    kb.NewID(),
		"url":   url,
		"text":  text,
		"items": items,
	}
}
