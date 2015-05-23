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

func (db *Index) selectPages(filter string, args ...interface{}) ([]kb.PageEntry, error) {
	rows, err := db.Query(`
	SELECT
		Owner,
		Slug,
		Data->>'title' as Title,
		Data->>'synopsis' as Synopsis,
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
			OR Owner IN (SELECT GroupID FROM Memberships WHERE UserID = $1)
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
		WHERE Owner IN (SELECT Name    FROM Groups      WHERE Public = True)
		   OR Owner IN (SELECT GroupID FROM Memberships WHERE UserID = $1)
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
	ntags := stringSlice{ntag}
	return db.selectPages(`
		WHERE (NormTags @> $1) 
		  AND (    Owner IN (SELECT Name    FROM Groups      WHERE Public = TRUE)
				OR Owner IN (SELECT GroupID FROM Memberships WHERE UserID = $2))
		ORDER BY Owner, Slug
	`, ntags, db.User)
}

func (db *Index) Groups() ([]kbserver.Group, error) {
	rows, err := db.Query(`
		SELECT ID, Name, Public, Description
		FROM Groups
		WHERE Public = True
		   OR ID IN (SELECT GroupID FROM Memberships WHERE UserID = $1)
	`, db.User)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []kbserver.Group
	for rows.Next() {
		var group kbserver.Group
		rows.Scan(&group.ID, &group.Name, &group.Public, &group.Description)
		groups = append(groups, group)
	}
	return groups, nil
}

func (db *Index) ByGroup(group kb.Slug) ([]kb.PageEntry, error) {
	if !db.CanRead(db.User, group) {
		return nil, kbserver.ErrUserNotAllowed
	}

	return db.selectPages(`
		WHERE (Owner = $1) 
		ORDER BY Owner, Slug
	`, group)
}

func (db *Index) RecentChanges(n int) ([]kb.PageEntry, error) {
	return db.selectPages(`
		WHERE  Owner IN (SELECT Name    FROM Groups      WHERE Public = TRUE)
			OR Owner IN (SELECT GroupID FROM Memberships WHERE UserID = $1)
		ORDER BY Modified DESC, Owner, Slug
		LIMIT $2
	`, db.User, n)
}
