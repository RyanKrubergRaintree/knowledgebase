package pqindex

import (
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/egonelbre/fedwiki"
)

/*
type Index interface {
	All() ([]*fedwiki.PageHeader, error)

	Tags() ([]TagInfo, error)
	PagesByTag(tag string) ([]*fedwiki.PageHeader, error)
	PagesCiting(slug fedwiki.Slug) ([]*fedwiki.PageHeader, error)

	Search(content string) ([]*fedwiki.PageHeader, error)
	RecentChanges(n int) ([]*fedwiki.PageHeader, error)
}

type PageStore interface {
	// Exists checks whether a page with `slug` exists
	Exists(slug Slug) bool
	// Create adds a new `page` with `slug`
	Create(slug Slug, page *Page) error
	// Load loads the page with identified by `slug`
	Load(slug Slug) (*Page, error)
	// Save saves the new page to `slug`
	Save(slug Slug, page *Page) error
	// List lists all the page headers
	List() ([]*PageHeader, error)
}
*/

type Store struct {
	DB *sql.DB
}

func New(params string) (*Store, error) {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, err
	}

	store := &Store{db}
	return store, store.init()
}

func (store *Store) init() error {

	return nil
}

type Index struct {
	store fedwiki.PageStore
}

func NewIndex()
