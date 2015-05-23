package pgdb

import (
	"database/sql"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

func (db *Database) Groups() kbserver.Groups { return &Groups{db} }

type Groups struct{ *Database }

func (db *Groups) ByID(name kb.Slug) (kbserver.Group, error) {
	var group kbserver.Group
	err := db.QueryRow(`
		SELECT
			ID, Name, Public, Description
		FROM Groups
		WHERE ID = $1
	`, name).Scan(&group.ID, &group.Name, &group.Public, &group.Description)
	if err == sql.ErrNoRows {
		return group, kbserver.ErrUserNotExist
	}
	if err != nil {
		return group, err
	}
	return group, nil
}

func (db *Groups) Create(group kbserver.Group) error {
	return db.exec(`
		INSERT INTO Groups 
		(ID, Name, Public, Description)
		VALUES ($1, $2, $3, $4)`,
		kb.Slugify(group.Name), group.Name, group.Public, group.Description)
}

func (db *Groups) Delete(group kb.Slug) error {
	return db.exec(`DELETE FROM Groups WHERE Name = $1`, group)
}

func (db *Groups) List() ([]kbserver.Group, error) {
	rows, err := db.Query(`
		SELECT ID, Name, Public, Description
		FROM Groups
	`)
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

func (db *Groups) AddMember(group, user kb.Slug) error {
	return db.exec(`
		INSERT INTO Memberships
		(GroupID, UserID)
		VALUES ($1, $2)`,
		group, user)
}

func (db *Groups) RemoveMember(group, user kb.Slug) error {
	return db.exec(`
		DELETE
		FROM Memberships
		WHERE GroupID = $1 AND UserID = $2
	`, group, user)
}
