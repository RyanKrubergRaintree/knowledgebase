package kb

import (
	"html"
	"time"

	"github.com/raintreeinc/knowledgebase/internal/natural"

	"github.com/bradfitz/slice"
)

type PageEntry struct {
	Slug     Slug      `json:"slug"`
	Title    string    `json:"title"`
	Synopsis string    `json:"synopsis"`
	Tags     []string  `json:"tags"`
	Modified time.Time `json:"modified"`
}

func (page *PageEntry) HasTag(tag string) bool {
	for _, t := range page.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func PageEntryFrom(page *Page) PageEntry {
	return PageEntry{
		Slug:     page.Slug,
		Title:    page.Title,
		Synopsis: page.Synopsis,
		Tags:     ExtractTags(page),
		Modified: page.Modified,
	}
}

type TagEntry struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func SortPageEntriesByDate(xs []PageEntry) {
	slice.Sort(xs, func(i, j int) bool { return xs[i].Modified.After(xs[j].Modified) })
}

func SortPageEntriesByTitle(xs []PageEntry) {
	slice.Sort(xs, func(i, j int) bool {
		return natural.Less(xs[i].Title, xs[j].Title)
	})
}

func ReversePageEntries(xs []PageEntry) {
	for i, j := 0, len(xs)-1; i < j; i, j = i+1, j-1 {
		xs[i], xs[j] = xs[j], xs[i]
	}
}

func SortPageEntriesBySlug(xs []PageEntry) {
	slice.Sort(xs, func(i, j int) bool { return xs[i].Slug < xs[j].Slug })
}

func SortTagEntriesByName(xs []TagEntry) {
	slice.Sort(xs, func(i, j int) bool { return xs[i].Name < xs[j].Name })
}

func StoryFromEntries(entries []PageEntry) Story {
	story := Story{}
	if len(entries) == 0 {
		story.Append(Paragraph("No results."))
		return story
	}
	for _, entry := range entries {
		story.Append(Entry(
			html.EscapeString(entry.Title),
			html.EscapeString(entry.Synopsis),
			entry.Slug,
		))
	}
	return story
}
