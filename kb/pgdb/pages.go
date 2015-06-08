package pgdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Pages struct {
	Context
	GroupID kb.Slug
}

func (db Pages) record(action string, slug kb.Slug, version int, v interface{}) {
	data, _ := json.Marshal(v)
	_, err := db.Exec(`
		INSERT INTO
		PageJournal(Actor, Slug, Version, Action, Data)
		VALUES($1, $2, $3, $4, $5)
	`, db.ActiveUser, slug, version, action, data)
	if err != nil {
		log.Println(err)
	}
}
func (db Pages) Create(page *kb.Page) error {
	owner, _ := kb.TokenizeLink(string(page.Slug))
	if owner != db.GroupID {
		return fmt.Errorf("mismatching page.Slug (%s) and group (%s)", page.Slug, db.GroupID)
	}
	if err := kb.ValidateSlug(page.Slug); err != nil {
		return kb.ErrInvalidSlug
	}

	tags := kb.ExtractTags(page)
	tagSlugs := kb.SlugifyTags(tags)

	data, err := json.Marshal(page)
	if err != nil {
		return fmt.Errorf("failed to serialize page: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO Pages(
			OwnerID, Slug, Data, Version,
			Tags, TagSlugs,
			Created, Modified
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8
		)
	`, db.GroupID, page.Slug, data, page.Version,
		stringSlice(tags), stringSlice(tagSlugs),
		page.Modified, page.Modified)

	if dupkey(err) {
		return kb.ErrPageExists
	}
	if err == nil {
		db.record("create", page.Slug, 0, page)
	}
	return err
}

func (db Pages) Load(id kb.Slug) (*kb.Page, error) {
	data, err := db.LoadRaw(id)
	if err != nil {
		return nil, err
	}
	page := &kb.Page{}
	err = json.Unmarshal(data, page)
	return page, err
}

func (db Pages) LoadRaw(id kb.Slug) ([]byte, error) {
	var data []byte
	err := db.QueryRow(`
		SELECT Data
		FROM Pages
		Where OwnerID = $1 AND Slug = $2
	`, db.GroupID, id).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, kb.ErrPageNotExist
	}
	return data, err
}

func (db Pages) Overwrite(id kb.Slug, version int, page *kb.Page) error {
	owner, _ := kb.TokenizeLink(string(page.Slug))
	if owner != db.GroupID {
		return fmt.Errorf("mismatching page.Slug (%s) and group (%s)", page.Slug, db.GroupID)
	}

	tags := kb.ExtractTags(page)
	tagSlugs := kb.SlugifyTags(tags)
	//TODO: extract/update synopsis here
	// synopsis := kb.ExtractSynopsis(page)

	data, err := json.Marshal(page)
	if err != nil {
		return fmt.Errorf("failed to serialize page: %v", err)
	}

	r, err := db.Exec(`
		UPDATE Pages
		SET Data = $4,
			Version = $5,
			Tags = $6,
			TagSlugs = $7,
			Created = $8,
			Modified = $9
		WHERE OwnerID = $1 AND Slug = $2 AND Version = $3
	`, db.GroupID, page.Slug, version,
		data, page.Version, stringSlice(tags), stringSlice(tagSlugs),
		page.Modified, page.Modified)

	affected, _ := r.RowsAffected()
	if affected == 0 {
		return kb.ErrConcurrentEdit
	}
	if err == nil {
		db.record("overwrite", page.Slug, version, page)
	}
	return err
}

func (db Pages) Edit(id kb.Slug, version int, action kb.Action) error {
	page, err := db.Load(id)
	if err != nil {
		return err
	}
	if version > 0 && page.Version != version {
		return kb.ErrConcurrentEdit
	}
	version = page.Version
	page.Modified = time.Now()
	if err := page.Apply(action); err != nil {
		return err
	}

	db.record("try-edit", id, version, action)
	return db.Overwrite(id, version, page)
}

func (db Pages) Delete(id kb.Slug, version int) error {
	r, err := db.Exec(`
		DELETE FROM Pages
		WHERE Slug = $1 AND Version = $2
	`, id, version)
	affected, _ := r.RowsAffected()
	if err == sql.ErrNoRows || affected == 0 {
		return kb.ErrConcurrentEdit
	}
	if err != nil {
		db.record("delete", id, version, "")
	}
	return err
}

func (db Pages) BatchReplace(pages map[kb.Slug]*kb.Page) error {
	return ErrNotImplemented
}

func (db Pages) List() ([]kb.PageEntry, error) {
	return db.pageEntries(`
		WHERE OwnerID = $1
		ORDER BY Slug
	`, db.GroupID)
}

func (db Pages) Journal(id kb.Slug) ([]kb.Action, error) {
	return []kb.Action{}, ErrNotImplemented
}
