// This package implements common federated wiki types
package kb

import "strings"

func Paragraph(text string) Item {
	return Item{
		"type": "paragraph",
		"id":   NewID(),
		"text": text,
	}
}

func HTML(text string) Item {
	return Item{
		"type": "html",
		"id":   NewID(),
		"text": text,
	}
}

func Reference(title, site, text string) Item {
	return Item{
		"type":  "reference",
		"id":    NewID(),
		"title": title,
		"site":  site,
		"text":  text,
	}
}

func Image(caption, url, text string) Item {
	return Item{
		"type":    "image",
		"id":      NewID(),
		"url":     url,
		"text":    text,
		"caption": caption,
	}
}

func Entry(title, synopsis string, slug Slug) Item {
	return Item{
		"type":  "entry",
		"id":    slug,
		"title": title,
		"text":  synopsis,
		"rank":  0,
		"url":   slug,
	}
}

func Tags(tags ...string) Item {
	return Item{
		"type": "tags",
		"id":   NewID(),
		"text": strings.Join(tags, ", "),
	}
}
