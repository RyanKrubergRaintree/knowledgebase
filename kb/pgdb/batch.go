package pgdb

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/raintreeinc/knowledgebase/kb"
)

func (db Pages) BatchReplace(pages map[kb.Slug]*kb.Page, complete func(kb.Slug)) error {
	for slug := range pages {
		if owner, _ := kb.TokenizeLink(string(slug)); owner != db.GroupID {
			return errors.New("Invalid group replacement.")
		}
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// remove old data
	if _, err := tx.Exec("DELETE FROM Pages WHERE OwnerID = $1", db.GroupID); err != nil {
		return fmt.Errorf("failed to delete old pages: %v", err)
	}

	insert, err := tx.Prepare(`
		INSERT INTO Pages(
			OwnerID, Slug, Data, Version,
			Tags, TagSlugs,
			Created, Modified
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8
		)
	`)
	if err != nil {
		insert.Close()
		return fmt.Errorf("failed to create prepared statement: %v", err)
	}

	for _, page := range pages {
		tags := kb.ExtractTags(page)
		tagSlugs := kb.SlugifyTags(tags)
		data, err := json.Marshal(page)
		if err != nil {
			insert.Close()
			return fmt.Errorf("failed to serialize page: %v", err)
		}

		_, err = insert.Exec(
			db.GroupID, page.Slug, data, page.Version,
			stringSlice(tags), stringSlice(tagSlugs),
			page.Modified, page.Modified)
		if err != nil {
			insert.Close()
			return fmt.Errorf("failed to insert: %v", err)
		}

		complete(page.Slug)
	}

	insert.Close()
	return tx.Commit()
}
