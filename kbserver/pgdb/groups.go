package pgdb

import "github.com/raintreeinc/knowledgebase/kbserver"

func (db *Database) Groups() kbserver.Groups { return &Groups{db} }

type Groups struct{ *Database }

func (db *Groups) Create(group kbserver.Group) error {
	return db.exec(`
		INSERT INTO Groups 
		(Name, Public, Description)
		VALUES ($1, $2, $3)`,
		group.Name, group.Public, group.Description)
}

func (db *Groups) Delete(name string) error {
	return db.exec(`DELETE FROM Groups WHERE Name = $1`, name)
}

func (db *Groups) List() ([]kbserver.Group, error) {
	rows, err := db.Query(`
		SELECT Name, Public, Description
		FROM Groups
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []kbserver.Group
	for rows.Next() {
		var group kbserver.Group
		rows.Scan(&group.Name, &group.Public, &group.Description)
		groups = append(groups, group)
	}
	return groups, nil
}

func (db *Groups) AddMember(group string, user string) error {
	return db.exec(`
		INSERT INTO Memberships
		(GroupName, UserName)
		VALUES ($1, $2)`,
		group, user)
}

func (db *Groups) RemoveMember(group string, user string) error {
	return db.exec(`
		DELETE
		FROM Memberships
		WHERE GroupName = $1 AND UserName = $2
	`, group, user)
}
