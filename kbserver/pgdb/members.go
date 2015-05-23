package pgdb

import (
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

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

func (db *Groups) MembersOf(group kb.Slug) ([]kbserver.User, error) {
	rows, err := db.Query(`
		SELECT
			Users.ID, Users.Name, Users.Email, Users.Description
		FROM Users
		JOIN Memberships ON (Users.ID = Memberships.UserID)
		Where Memberships.GroupID = $1
		ORDER BY Users.ID
	`, group)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []kbserver.User
	for rows.Next() {
		var user kbserver.User
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Description)
		users = append(users, user)
	}
	return users, rows.Err()
}
