package pgdb

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

func (db *Database) PagesByGroup(user, group kb.Slug) kbserver.Pages {
	return &Pages{db, user, group}
}

type Pages struct {
	*Database
	User  kb.Slug
	Group kb.Slug
}

func (db *Pages) tx() (*sql.Tx, error) {
	tx, err := db.Begin()
	if err != nil {
		return tx, err
	}

	err = tx.QueryRow(`
		SELECT
		FROM Memberships
		WHERE UserID = $1 AND GroupID = $2
	`, db.User, db.Group).Scan()

	if err == sql.ErrNoRows {
		tx.Rollback()
		return nil, kbserver.ErrUserNotAllowed
	}
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	return tx, nil
}

func (db *Pages) CanWrite() (result bool) {
	err := db.QueryRow(`SELECT
		EXISTS(SELECT FROM Memberships WHERE UserID = $1 AND GroupID = $2)
	`, db.User, db.Group).Scan(&result)
	return (err == nil) && result
}

func (db *Pages) CanRead() (result bool) {
	err := db.QueryRow(`SELECT
		(SELECT Public FROM Groups WHERE ID = $2)
     	OR EXISTS( SELECT FROM Memberships WHERE UserID = $1 AND GroupID = $2 )
	`, db.User, db.Group).Scan(&result)
	return (err == nil) && result
}

func (db *Pages) Create(page *kb.Page) error {
	if page.Owner != db.Group {
		return fmt.Errorf("mismatching page.Owner (%s) and group (%s)", page.Owner, db.Group)
	}

	if !db.CanWrite() {
		return kbserver.ErrUserNotAllowed
	}

	tags := kb.ExtractTags(page)
	data, err := json.Marshal(page)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
		INSERT INTO Pages
		(Owner, Slug, Data, Tags)
		VALUES ($1, $2, $3, $4)
	`, page.Owner, page.Slug, data, stringSlice(tags))

	return err
}

func (db *Pages) Load(slug kb.Slug) (*kb.Page, error) {
	data, err := db.LoadRaw(slug)
	if err != nil {
		return nil, err
	}

	page := &kb.Page{}
	err = json.Unmarshal(data, page)
	return page, err
}

func (db *Pages) LoadRaw(slug kb.Slug) ([]byte, error) {
	if !db.CanRead() {
		return nil, kbserver.ErrUserNotAllowed
	}

	var data []byte
	err := db.QueryRow(`
		SELECT Data
		FROM Pages
		WHERE Owner = $1 AND Slug = $2
	`, db.Group, slug).Scan(&data)

	if err == sql.ErrNoRows {
		return nil, kbserver.ErrPageNotExist
	}

	return data, err
}

func (db *Pages) Save(slug kb.Slug, page *kb.Page) error {
	if page.Owner != db.Group {
		return fmt.Errorf("mismatching page.Owner (%s) and group (%s)", page.Owner, db.Group)
	}
	if page.Slug != slug {
		return fmt.Errorf("mismatching page.Slug (%s) and slug (%s)", page.Slug, slug)
	}
	if !db.CanWrite() {
		return kbserver.ErrUserNotAllowed
	}

	tags := kb.ExtractTags(page)
	data, err := json.Marshal(page)
	if err != nil {
		return err
	}

	_, err = db.Exec(`
		UPDATE Pages
		SET	Tags = $1,
			Data = $2,
			Version = (Version + 1)
		WHERE Owner = $3 AND Slug = $4
	`, stringSlice(tags), data, db.Group, slug)

	return err
}

//func (db *Pages) SaveRaw(slug kb.Slug, page []byte) error {}

func (db *Pages) List() ([]kb.PageEntry, error) {
	if !db.CanRead() {
		return nil, kbserver.ErrUserNotAllowed
	}

	index := &Index{db.Database, db.User}
	return index.selectPages(`
		WHERE Owner = $1
		ORDER BY Slug
	`, db.Group)
}
