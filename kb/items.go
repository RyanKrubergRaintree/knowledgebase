// This package implements common federated wiki types
package kb

import (
	"strings"
	"unicode"
)

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
		"link":  slug,
	}
}

func Tags(tags ...string) Item {
	return Item{
		"type": "tags",
		"id":   NewID(),
		"text": strings.Join(tags, ", "),
	}
}

func ExtractTags(page *Page) []string {
	tags := make(map[string]string)
	for _, item := range page.Story {
		if item.Type() == "tags" {
			for _, tag := range strings.Split(item.Val("text"), ",") {
				ntag := string(Slugify(tag))
				tags[ntag] = strings.TrimSpace(tag)
			}
		}
	}

	result := make([]string, 0, len(tags))
	for _, tag := range tags {
		result = append(result, tag)
	}
	return result
}

func SlugifyTags(tags []string) []string {
	normalized := make([]string, 0, len(tags))
	for _, tag := range tags {
		normalized = append(normalized, string(Slugify(tag)))
	}

	return normalized
}

func limitWords(text string, limit int) string {
	words := strings.Fields(text)
	if len(words) > limit {
		words = words[:limit]
	}
	r := strings.Join(words, " ")
	if len(r) > 0 && unicode.IsLetter(rune(r[len(r)-1])) {
		r += "..."
	}
	return r
}

func ExtractSynopsis(page *Page) string {
	for _, item := range page.Story {
		if item.Type() == "paragraph" {
			text := item.Val("text")
			if text != "" {
				return limitWords(text, 30)
			}
		}
	}
	return ""
}
