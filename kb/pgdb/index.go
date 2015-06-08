package pgdb

import "github.com/raintreeinc/knowledgebase/kb"

type Index struct {
	Context
	UserID kb.Slug
}

const (
	publicGroup   = `SELECT ID FROM Groups WHERE Public = TRUE`
	directMember  = `SELECT GroupID FROM Membership WHERE UserID = $1`
	memberOfOwner = `
		SELECT ID FROM Groups
		JOIN Membership ON Membership.GroupID = Groups.OwnerID
		WHERE Membership.UserID = $1
	`
	memberOfCommunity = `
		SELECT ID FROM Groups
		JOIN Community  ON Community.GroupID = Groups.ID
		JOIN Membership ON Membership.GroupID = Community.MemberID
		WHERE Community.Access >= 'reader'
		  AND Membership.UserID = $1
	`
	//TODO: figure out whether this needs to be optimized
	pageVisibleToUser = `
		  (   OwnerID IN (` + publicGroup + `)
		   OR OwnerID IN (` + directMember + `)
		   OR OwnerID IN (` + memberOfOwner + `)
		   OR OwnerID IN (` + memberOfCommunity + `))
	`
)

func (db Index) List() ([]kb.PageEntry, error) {
	return db.pageEntries(`WHERE `+pageVisibleToUser+`ORDER BY Slug`, db.UserID)
}

func (db Index) Search(text string) ([]kb.PageEntry, error) {
	return db.pageEntries(`
		WHERE `+pageVisibleToUser+`
		AND	to_tsvector('english',
				coalesce(cast(Data->'title' AS TEXT),'') || ' ' ||
				coalesce(cast(Data->'story' AS TEXT), '')
			) @@ to_tsquery('english', $2);
		`, db.UserID, text)
}

func (db Index) Tags() ([]kb.TagEntry, error) {
	rows, err := db.Query(`
		SELECT
			unnest(Tags) as Tag,
			count(*) as Count
		FROM Pages
		WHERE `+pageVisibleToUser+`
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
		WHERE TagSlugs @> $2
		  AND `+pageVisibleToUser+`
	`, db.UserID, tagSlugs)
}

func (db Index) Groups(min kb.Rights) (groups []kb.Group, err error) {
	rows, err := db.Query(`
		SELECT  ID, OwnerID, Name, Public, Description
		FROM    Groups
		WHERE Public
		   OR ID IN (SELECT GroupID FROM Membership WHERE UserID = $1)
		   OR OwnerID IN (SELECT GroupID FROM Membership WHERE UserID = $1)
		   OR ID IN (
				SELECT Community.GroupID FROM Community
				JOIN  Membership ON Community.MemberID = Membership.GroupID
				WHERE Membership.UserID = $1
				  AND Community.Access >= $2
			)
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
		WHERE OwnerID = $2
		  AND `+pageVisibleToUser+`
	`, db.UserID, groupID)
}

func (db Index) RecentChanges() ([]kb.PageEntry, error) {
	n := 30
	return db.pageEntries(`
		WHERE `+pageVisibleToUser+`
		ORDER BY Modified DESC, OwnerID, Slug
		LIMIT $2
	`, db.UserID, n)
}
