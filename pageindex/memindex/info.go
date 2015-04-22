package memindex

import (
	"sort"

	"github.com/egonelbre/fedwiki"
)

type citations map[fedwiki.Slug][]fedwiki.Slug
type tags map[string][]fedwiki.Slug
type headers []*fedwiki.PageHeader

type info struct {
	store fedwiki.PageStore

	byDate []*fedwiki.PageHeader
	bySlug []*fedwiki.PageHeader
	header map[fedwiki.Slug]fedwiki.PageHeader

	citations citations
	tags      tags
}

func newInfo(store fedwiki.PageStore) *info {
	info := &info{}

	info.byDate, _ = store.List()
	info.bySlug = make([]*fedwiki.PageHeader, len(info.byDate))
	copy(info.bySlug, info.byDate)

	info.citations = make(citations, len(info.byDate))
	info.tags = make(tags, len(info.byDate))

	sort.Sort(byDate(info.byDate))
	sort.Sort(bySlug(info.byDate))

	info.updateTags()
	info.updateCitations()

	return info
}

func (info *info) updateTags() {
	for _, header := range info.byDate {
		tags := header.Meta["tags"]
		if tagnames, ok := tags.([]string); ok {
			for _, tag := range tagnames {
				info.tags[tag] = append(info.tags[tag], header.Slug)
			}
		}
	}
}

func (info *info) updateCitations() {
	//TODO
}

type byDate []*fedwiki.PageHeader

func (s byDate) Len() int           { return len(s) }
func (s byDate) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s byDate) Less(i, j int) bool { return s[j].Date.Before(s[i].Date.Time) }

// TODO: use humanized sorting
type bySlug []*fedwiki.PageHeader

func (s bySlug) Len() int           { return len(s) }
func (s bySlug) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s bySlug) Less(i, j int) bool { return s[i].Slug < s[j].Slug }
