package kbserver

import (
	"errors"

	"github.com/raintreeinc/knowledgebase/kb"
)

var (
	ErrInvalidUser    = errors.New("invalid user")
	ErrUserNotAllowed = errors.New("user does not have sufficient permissions")
	ErrPageMissing    = errors.New("page does not exist")
)

type Database interface {
	PagesByOwner(user, owner string) (Pages, error)
	IndexByUser(user string) (Index, error)
}

type Pages interface {
	All() ([]kb.PageEntry, error)
	Exists(slug kb.Slug) bool
	Create(slug kb.Slug, page *kb.Page) error
	Load(slug kb.Slug) (*kb.Page, error)
	LoadRaw(slug kb.Slug) ([]byte, error)
	Save(slug kb.Slug, page *kb.Page) error
	// SaveRaw(slug kb.Slug, page []byte) error
}

type Index interface {
	All() ([]kb.PageEntry, error)
	Search(text string) ([]kb.PageEntry, error)

	Tags() ([]kb.TagEntry, error)
	ByTag(tag string) ([]kb.PageEntry, error)

	RecentChanges(n int) ([]kb.PageEntry, error)
}
