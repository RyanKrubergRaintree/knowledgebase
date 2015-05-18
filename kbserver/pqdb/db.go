package pqdb

import (
	"database/sql"
	"fmt"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

type Database struct {
	db *sql.DB
}

func New(params string) *Database {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, fmt.Errorf("failed to load DB: %s", err)
	}
	return Database{db}
}

func (db *Database) PagesByOwner(user, owner int) (kbserver.Pages, error) {
	// TODO: verify access
	return &Pages{db, owner}, nil
}

func (db *Database) IndexByUser(user int) (kbserver.Index, error) {
	// TODO: verify access
	return &Index{db, owner}, nil
}

type Pages struct {
	db    *sql.DB
	owner int
}

func (pages *Pages) All() ([]int, error) {

}

func (pages *Pages) Exists(slug kb.Slug) bool {

}

func (pages *Pages) Create(slug kb.Slug, page *kb.Page) error {

}

func (pages *Pages) Load(slug kb.Slug) (*kb.Page, error) {

}

func (pages *Pages) LoadRaw(slug kb.Slug) ([]byte, error) {

}

func (pages *Pages) Save(slug kb.Slug, page *kb.Page) error {

}

type Index struct {
	db   *sql.DB
	user int
}

func (index *Index) All() ([]int, error) {

}

func (index *Index) Search(text string) ([]int, error) {

}

func (index *Index) Tags() ([]int, error) {

}

func (index *Index) ByTag(tag string) ([]int, error) {

}

func (index *Index) RecentChanges(n int) ([]int, error) {

}
