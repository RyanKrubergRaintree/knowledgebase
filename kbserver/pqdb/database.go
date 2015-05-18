package pqdb

import (
	"database/sql"
	"fmt"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

type Database struct {
	DB *sql.DB
}

func New(params string) *Database {
	db, err := sql.Open("postgres", params)
	if err != nil {
		return nil, fmt.Errorf("failed to load DB: %s", err)
	}
	return Database{db}
}

func (db *Database) ValidUser(user string) error {
	//TODO
	return nil
}

func (db *Database) CanAccess(user, owner string) error {
	if err := db.ValidUser(user); err != nil {
		return kbserver.ErrInvalidUser
	}

	//TODO: check whether belongs to owner
	return nil
}

func (db *Database) PagesByOwner(user, owner string) (kbserver.Pages, error) {
	if err := db.CanAccess(user, owner); err != nil {
		return nil, err
	}
	return &Pages{db, 0}, nil
}

func (db *Database) IndexByUser(user string) (kbserver.Index, error) {
	if err := db.ValidUser(user); err != nil {
		return nil, err
	}
	return &Index{db, 0}, nil
}

type Pages struct {
	DB    *sql.DB
	Owner int
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
	DB   *sql.DB
	User int
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
