package pgdb

import (
	"errors"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

func (db *Database) IndexByUser(user kb.Slug) kbserver.Index {
	return &Index{db, user}
}

type Index struct {
	*Database
	User kb.Slug
}

const selectPages = `
	SELECT
		Owner,
		Slug,
		Data->'title' as Title,
		Data->'synopsis' as Synopsis,
		Tags,
		Modified
	FROM Pages
`

func (db *Index) selectPages(filter string, args ...interface{}) ([]kb.PageEntry, error) {
	rows, err := db.Query(`SELECT
		Owner,
		Slug,
		Data->'title' as Title,
		Data->'synopsis' as Synopsis,
		Tags,
		Modified
	FROM Pages
	`+filter, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []kb.PageEntry
	for rows.Next() {
		var entry kb.PageEntry

		xtags := stringSlice{}
		err := rows.Scan(
			&entry.Owner,
			&entry.Slug,
			&entry.Title,
			&entry.Synopsis,
			&xtags,
			&entry.Modified,
		)
		entry.Tags = []string(xtags)

		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (db *Index) List() ([]kb.PageEntry, error) {
	return db.selectPages(`
		WHERE  Owner IN (SELECT Name    FROM Groups      WHERE Public = TRUE)
			OR Owner IN (SELECT GroupID FROM Memberships WHERE User = $1)
		ORDER BY Owner, Slug
	`, db.User)
}

func (db *Index) Search(text string) ([]kb.PageEntry, error) {
	return nil, errors.New("not implemented")
}

func (db *Index) Tags() ([]kb.TagEntry, error) {
	rows, err := db.Query(`
		SELECT
			unnest(Tags) as Tag,
			count(*) as Count
		FROM Pages
		WHERE Owner IN (SELECT GroupID FROM Memberships WHERE User = $1)
		   OR Owner IN (SELECT Name    FROM Groups WHERE Public = True)
		GROUP BY Tag
	`, db.User)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []kb.TagEntry
	for rows.Next() {
		var tag kb.TagEntry
		err := rows.Scan(&tag.Name, &tag.Count)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (db *Index) ByTag(tag string) ([]kb.PageEntry, error) {
	ntag := string(kb.Slugify(tag))
	return db.selectPages(`
		WHERE (NTags @> $1) 
		  AND (    Owner IN (SELECT Name    FROM Groups      WHERE Public = TRUE)
				OR Owner IN (SELECT GroupID FROM Memberships WHERE User = $2))
		ORDER BY Owner, Slug
	`, ntag, db.User)
}

func (db *Index) RecentChanges(n int) ([]kb.PageEntry, error) {
	return db.selectPages(`
		WHERE  Owner IN (SELECT Name    FROM Groups      WHERE Public = TRUE)
			OR Owner IN (SELECT GroupID FROM Memberships WHERE User = $1)
		ORDER BY Modified DESC, Owner, Slug
		LIMIT $2
	`, db.User, n)
}
