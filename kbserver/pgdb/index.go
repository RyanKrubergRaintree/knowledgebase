package pgdb

import (
	"errors"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

func (db *Database) IndexByUser(user string) kbserver.Index {
	return &Index{db, user}
}

type Index struct {
	*Database
	User string
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
		err := rows.Scan(
			&entry.Owner,
			&entry.Slug,
			&entry.Title,
			&entry.Synopsis,
			&entry.Tags,
			&entry.Modified,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, rows.Err()
}

func (db *Index) List() ([]kb.PageEntry, error) {
	return db.selectPages(`
		WHERE  Owner IN (SELECT Name      FROM Groups      WHERE Public = TRUE)
			OR Owner IN (SELECT GroupName FROM Memberships WHERE User = ?)
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
		WHERE Owner IN (SELECT GroupName FROM Memberships WHERE User = ?)
		   OR Owner IN (SELECT Name FROM Groups WHERE Public = True)
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
	return db.selectPages(`
		WHERE (Tags @> ?) 
		  AND (    Owner IN (SELECT Name      FROM Groups      WHERE Public = TRUE)
				OR Owner IN (SELECT GroupName FROM Memberships WHERE User = ?))
		ORDER BY Owner, Slug
	`, tag, db.User)
}

func (db *Index) RecentChanges(n int) ([]kb.PageEntry, error) {
	return db.selectPages(`
		WHERE  Owner IN (SELECT Name      FROM Groups      WHERE Public = TRUE)
			OR Owner IN (SELECT GroupName FROM Memberships WHERE User = ?)
		ORDER BY Modified DESC, Owner, Slug
		LIMIT ?
	`, db.User, n)
}
