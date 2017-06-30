package kb

import (
	"sort"
	"time"
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
	sort.Slice(xs, func(i, j int) bool { return xs[i].Modified.After(xs[j].Modified) })
}

func SortPageEntries(xs []PageEntry, fn func(a, b *PageEntry) bool) {
	sort.Slice(xs, func(i, j int) bool {
		return fn(&xs[i], &xs[j])
	})
}

func SortPageEntriesBySlug(xs []PageEntry) {
	sort.Slice(xs, func(i, j int) bool { return xs[i].Slug < xs[j].Slug })
}

func SortTagEntriesByName(xs []TagEntry) {
	sort.Slice(xs, func(i, j int) bool { return xs[i].Name < xs[j].Name })
}

func SortPageEntriesByRank(xs []PageEntry, ranking []Slug) {

	calculateRank := func(tags []string) int {
		rank := len(ranking)
		for _, tag := range tags {
			stag := Slugify(tag)
			for k, target := range ranking[:rank] {
				if stag == target {
					rank = k
					break
				}
			}
			if rank == 0 {
				break
			}
		}
		return rank
	}

	type PageRank struct {
		Entry PageEntry
		Rank  int
	}

	order := make([]PageRank, len(xs))
	for i, x := range xs {
		order[i].Entry = x
		order[i].Rank = calculateRank(x.Tags)
	}

	sort.Slice(order, func(i, k int) bool { return order[i].Rank < order[k].Rank })

	for i, x := range order {
		xs[i] = x.Entry
	}
}

func StoryFromEntries(entries []PageEntry) Story {
	story := Story{}
	if len(entries) == 0 {
		story.Append(Paragraph("No pages."))
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
