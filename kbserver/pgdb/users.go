package pgdb

import (
	"database/sql"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

func (db *Database) Users() kbserver.Users { return &Users{db} }

type Users struct{ *Database }

// Create adds a new user to the database
// Groups field will be ignored
func (db *Users) Create(user kbserver.User) error {
	return db.exec(`
		INSERT INTO Users 
		(ID, Name, Email, Description)
		VALUES ($1, $2, $3, $4)`,
		kb.Slugify(user.Name), user.Name, user.Email, user.Description)
}

func (db *Users) Delete(name kb.Slug) error {
	return db.exec(`DELETE FROM Users WHERE Name = ?`, name)
}

func (db *Users) ByName(name kb.Slug) (kbserver.User, error) {
	var user kbserver.User
	err := db.QueryRow(`
		SELECT
			ID, Name, Email, Description,
			array_agg(Memberships.GroupName) as Groups
		FROM Users
		JOIN Memberships ON (Users.ID = Memberships.UserID)
		GROUP BY Users.ID
		WHERE Users.ID = $1
	`, name).Scan(&user.Name, &user.Name, &user.Email, &user.Description, &user.Groups)
	if err == sql.ErrNoRows {
		return user, kbserver.ErrUserNotExist
	}
	return user, nil
}

func (db *Users) List() ([]kbserver.User, error) {
	rows, err := db.Query(`
		SELECT
			ID, Name, Email, Description,
			array_agg(Memberships.GroupName) as Groups
		FROM Users
		JOIN Memberships ON (Users.ID = Memberships.UserID)
		GROUP BY Users.ID
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []kbserver.User
	for rows.Next() {
		var user kbserver.User
		rows.Scan(&user.ID, &user.Name, &user.Email, &user.Description, &user.Groups)
		users = append(users, user)
	}
	return users, rows.Err()
}
