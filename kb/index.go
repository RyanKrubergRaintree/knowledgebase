package kb

import (
	"time"

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

func SortPageEntries(xs []PageEntry, fn func(a, b *PageEntry) bool) {
	slice.Sort(xs, func(i, j int) bool {
		return fn(&xs[i], &xs[j])
	})
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
			entry.Title,
			entry.Synopsis,
			entry.Slug,
		))
	}
	return story
}
