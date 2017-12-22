package pgdb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/raintreeinc/knowledgebase/kb"
)

type pageInfo struct {
	Page     *kb.Page
	Tags     []string
	TagSlugs []string
	Data     []byte
	Hash     []byte
}

func (db Pages) createPageInfos(pages map[kb.Slug]*kb.Page) (map[kb.Slug]*pageInfo, error) {
	infos := make(map[kb.Slug]*pageInfo, len(pages))
	for slug, page := range pages {
		if owner, _ := kb.TokenizeLink(string(slug)); owner != db.GroupID {
			return nil, errors.New("Invalid group replacement.")
		}

		data, err := json.Marshal(page)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize page: %v", err)
		}

		hash, err := page.Hash()
		if err != nil {
			return nil, fmt.Errorf("failed to get page hash: %v", err)
		}

		tags := kb.ExtractTags(page)
		tagSlugs := kb.SlugifyTags(tags)

		infos[slug] = &pageInfo{
			Page:     page,
			Tags:     tags,
			TagSlugs: tagSlugs,
			Data:     data,
			Hash:     hash,
		}
	}

	return infos, nil
}

func (db Pages) BatchReplace(pages map[kb.Slug]*kb.Page, complete func(string, kb.Slug)) error {
	infos, err := db.createPageInfos(pages)
	if err != nil {
		return err
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
			Created, Modified, Hash
		) VALUES (
			$1, $2, $3, $4,
			$5, $6,
			$7, $8, $9
		)
	`)
	if err != nil {
		insert.Close()
		return fmt.Errorf("failed to create prepared statement: %v", err)
	}

	for _, info := range infos {
		_, err = insert.Exec(
			db.GroupID, info.Page.Slug, info.Data, info.Page.Version,
			stringSlice(info.Tags), stringSlice(info.TagSlugs),
			info.Page.Modified, info.Page.Modified, info.Hash)
		if err != nil {
			insert.Close()
			return fmt.Errorf("failed to insert: %v", err)
		}

		complete("inserted", info.Page.Slug)
	}

	insert.Close()
	return tx.Commit()
}

func (db Pages) BatchReplaceDelta(pages map[kb.Slug]*kb.Page, complete func(string, kb.Slug)) error {
	infos, err := db.createPageInfos(pages)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	oldHashes := map[kb.Slug][]byte{}
	{
		rows, err := tx.Query("SELECT Slug, Hash FROM Pages WHERE OwnerID = $1", db.GroupID)
		if err != nil {
			return fmt.Errorf("failed to get current headers: %v", err)
		}

		for rows.Next() {
			var slug kb.Slug
			var hash []byte

			if err := rows.Scan(&slug, &hash); err != nil {
				return fmt.Errorf("failed to get header row: %v", err)
			}
			oldHashes[slug] = hash
		}

		if err := rows.Close(); err != nil {
			return fmt.Errorf("failed to close rows query: %v", err)
		}
	}

	{ // remove missing pages
		del, err := tx.Prepare("DELETE FROM Pages WHERE Slug = $1")
		if err != nil {
			return fmt.Errorf("failed to create delete statement: %v", err)
		}

		for oldslug, oldHash := range oldHashes {
			info, stillExists := infos[oldslug]
			if stillExists && bytes.Equal(info.Hash, oldHash) {
				continue
			}
			if _, err := del.Exec(oldslug); err != nil {
				return fmt.Errorf("failed to exec del statement: %v", err)
			}

			if !stillExists {
				complete("deleted", oldslug)
			}
		}

		if err := del.Close(); err != nil {
			return fmt.Errorf("failed to close del statement: %v", err)
		}
	}

	insert, err := tx.Prepare(`
		INSERT INTO Pages(
			OwnerID, Slug,
			Data, Version, Tags, TagSlugs,
			Created, Modified,
			Hash
		) VALUES (
			$1, $2,
			$3, $4, $5, $6,
			$7, $8,
			$9
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create insert statement: %v", err)
	}

	for _, info := range infos {
		oldHash, exists := oldHashes[info.Page.Slug]
		if exists && bytes.Equal(info.Hash, oldHash) {
			complete("unchanged", info.Page.Slug)
			continue
		}

		_, err = insert.Exec(
			db.GroupID, info.Page.Slug, info.Data, info.Page.Version,
			stringSlice(info.Tags), stringSlice(info.TagSlugs),
			info.Page.Modified, info.Page.Modified,
			info.Hash)
		if err != nil {
			insert.Close()
			return fmt.Errorf("failed to insert: %v", err)
		}

		if exists {
			complete("updated", info.Page.Slug)
		} else {
			complete("added", info.Page.Slug)
		}
	}

	insert.Close()
	return tx.Commit()
}
