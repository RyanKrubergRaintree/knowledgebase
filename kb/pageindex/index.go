package pageindex

import "github.com/egonelbre/fedwiki"

type TagInfo struct {
	Name  string
	Count int
}

type Index interface {
	All() ([]*fedwiki.PageHeader, error)

	Tags() ([]TagInfo, error)
	PagesByTag(tag string) ([]*fedwiki.PageHeader, error)
	PagesCiting(slug fedwiki.Slug) ([]*fedwiki.PageHeader, error)

	Search(content string) ([]*fedwiki.PageHeader, error)
	RecentChanges(n int) ([]*fedwiki.PageHeader, error)
}
