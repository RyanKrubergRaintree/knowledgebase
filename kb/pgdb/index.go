package pgdb

import "github.com/raintreeinc/knowledgebase/kb"

type Index struct {
	Context
	UserID kb.Slug
}

func (db Index) List() ([]kb.PageEntry, error) {
	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		ORDER BY Slug`, db.UserID)
}

func (db Index) Search(text string) ([]kb.PageEntry, error) {
	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		  AND Content @@ plainto_tsquery('english', $2)
		ORDER BY ts_rank(Content, plainto_tsquery('english', $2)) DESC
		`, db.UserID, text)
}

func (db Index) SearchFilter(text, exclude, include string) ([]kb.PageEntry, error) {
	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		  AND (OwnerID NOT LIKE $3 || '%' OR OwnerID = $4)
		  AND Content @@ plainto_tsquery('english', $2)
		ORDER BY ts_rank(Content, plainto_tsquery('english', $2)) DESC
		`, db.UserID, text, exclude, include)
}

func (db Index) Tags() ([]kb.TagEntry, error) {
	rows, err := db.Query(`
		SELECT
			unnest(Tags) as Tag,
			count(*) as Count
		FROM Pages
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1 AND AccessView.Access >= 'reader'
		GROUP BY Tag
		ORDER BY Tag
	`, db.UserID)
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

func (db Index) ByTag(tag kb.Slug) ([]kb.PageEntry, error) {
	tags := kb.SlugifyTags([]string{string(tag)})
	tagSlugs := stringSlice(tags)

	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		  AND TagSlugs && $2
	`, db.UserID, tagSlugs)
}

func (db Index) ByTagFilter(tags []kb.Slug, exclude, include string) ([]kb.PageEntry, error) {
	if len(tags) == 0 {
		return nil, nil
	}

	tagSlugs := make([]string, len(tags))
	for i, tag := range tags {
		tagSlugs[i] = string(tag)
	}

	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		  AND (OwnerID NOT LIKE $3 || '%' OR OwnerID = $4)
		  AND TagSlugs && $2
		`, db.UserID, stringSlice(tagSlugs), exclude, include)
}

func (db Index) readable() (groups []kb.Group, err error) {
	rows, err := db.Query(`
		SELECT  ID, OwnerID, Name, Public, Description
		FROM    Groups
		JOIN AccessView ON Groups.ID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
	`, db.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group kb.Group
		rows.Scan(&group.ID, &group.OwnerID, &group.Name, &group.Public, &group.Description)
		groups = append(groups, group)
	}
	return groups, nil
}

func (db Index) Groups(min kb.Rights) (groups []kb.Group, err error) {
	if min == kb.Reader {
		return db.readable()
	}

	user, err := db.Users().ByID(db.UserID)
	if err != nil || user.MaxAccess.Level() < min.Level() {
		return []kb.Group{}, err
	}

	rows, err := db.Query(`
		SELECT  ID, OwnerID, Name, Public, Description
		FROM    Groups
		JOIN AccessView ON Groups.ID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= $2
	`, db.UserID, string(min))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group kb.Group
		rows.Scan(&group.ID, &group.OwnerID, &group.Name, &group.Public, &group.Description)
		groups = append(groups, group)
	}
	return groups, nil
}

func (db Index) ByGroup(groupID kb.Slug) ([]kb.PageEntry, error) {
	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		  AND OwnerID = $2
	`, db.UserID, groupID)
}

func (db Index) ByTitle(suffix kb.Slug) ([]kb.PageEntry, error) {
	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		  AND Slug LIKE '%=' || $2
	`, db.UserID, suffix)
}

func (db Index) RecentChanges(n int) ([]kb.PageEntry, error) {
	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		ORDER BY Modified DESC, OwnerID, Slug
		LIMIT $2
	`, db.UserID, n)
}

func (db Index) RecentChangesByGroup(n int, groupID kb.Slug) ([]kb.PageEntry, error) {
	return db.pageEntries(`
		JOIN AccessView ON OwnerID = AccessView.GroupID
		WHERE AccessView.UserID = $1
		  AND AccessView.Access >= 'reader'
		  AND OwnerID = $2
		ORDER BY Modified DESC, OwnerID, Slug
		LIMIT $3
	`, db.UserID, groupID, n)
}
