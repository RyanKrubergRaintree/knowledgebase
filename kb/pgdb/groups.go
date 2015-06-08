package pgdb

import (
	"database/sql"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Groups struct{ Context }

func (db Groups) ByID(id kb.Slug) (group kb.Group, err error) {
	err = db.QueryRow(`
		SELECT  ID, OwnerID, Name, Public, Description
		FROM    Groups
		WHERE   ID = $1
	`, id).Scan(&group.ID, &group.OwnerID, &group.Name, &group.Public, &group.Description)

	if err == sql.ErrNoRows {
		return group, kb.ErrGroupNotExist
	}
	return group, err
}

func (db Groups) Create(group kb.Group) error {
	_, err := db.Exec(`
		INSERT INTO
		Groups (ID, OwnerID, Name, Public, Description)
		VALUES ($1, $2, $3, $4, $5)
	`, group.ID, group.OwnerID, group.Name, group.Public, group.Description)

	if dupkey(err) {
		return kb.ErrGroupExists
	}
	return err
}

func (db Groups) Delete(id kb.Slug) error {
	_, err := db.Exec(`DELETE FROM Groups WHERE ID = $1`, id)
	return err
}

func (db Groups) List() (groups []kb.Group, err error) {
	rows, err := db.Query(`
		SELECT  ID, OwnerID, Name, Public, Description
		FROM    Groups
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var group kb.Group
		err := rows.Scan(&group.ID, &group.OwnerID, &group.Name, &group.Public, &group.Description)
		if err != nil {
			return groups, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}
