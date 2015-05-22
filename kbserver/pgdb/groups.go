package pgdb

import "github.com/raintreeinc/knowledgebase/kbserver"

func (db *Database) Groups() kbserver.Groups { return &Groups{db} }

type Groups struct{ *Database }

func (db *Groups) Create(group kbserver.Group) error {
	return db.exec(`
		INSERT INTO Group 
		(Name, Public, Description)
		VALUES
		(?, ?, ?)`,
		group.Name, group.Public, group.Description)
}

func (db *Groups) Delete(name string) error {
	return db.exec(`DELETE FROM Group WHERE Name = ?`, name)
}

func (db *Groups) List() ([]kbserver.Group, error) {
	rows, err := db.Query(`
		SELECT Name, Public, Description
		FROM Group
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
