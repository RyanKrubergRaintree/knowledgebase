package pageindex

import (
	"github.com/egonelbre/fedwiki"
	"github.com/egonelbre/fedwiki/item"
)

func HeadersToStory(headers []*fedwiki.PageHeader) fedwiki.Story {
	story := fedwiki.Story{}
	if len(headers) == 0 {
		story.Append(item.Paragraph("No results found."))
		return story
	}

	for _, h := range headers {
		story.Append(Entry(h.Title, h.Synopsis, h.Slug))
	}

	return story
}

func Entry(title, synopsis string, slug fedwiki.Slug) fedwiki.Item {
	return fedwiki.Item{
		"type":  "entry",
		"id":    slug,
		"title": title,
		"text":  synopsis,
		"rank":  0,
		"url":   slug,
	}
}
