package kbserver

import (
	"errors"

	"github.com/raintreeinc/knowledgebase/kb"
)

var (
	ErrUserNotExist   = errors.New("user does not exist")
	ErrUserNotAllowed = errors.New("user does not have sufficient permissions")
	ErrPageMissing    = errors.New("page does not exist")
)

type Database interface {
	PagesByOwner(user, owner int) (Pages, error)
	IndexByUser(user int) (Index, error)
}

type Pages interface {
	All() ([]int, error)
	Exists(slug kb.Slug) bool
	Create(slug kb.Slug, page *kb.Page) error
	Load(slug kb.Slug) (*kb.Page, error)
	LoadRaw(slug kb.Slug) ([]byte, error)
	Save(slug kb.Slug, page *kb.Page) error
	// SaveRaw(slug kb.Slug, page []byte) error
}

type Index interface {
	All() ([]int, error)
	Search(text string) ([]int, error)

	Tags() ([]int, error)
	ByTag(tag string) ([]int, error)

	RecentChanges(n int) ([]int, error)
}
