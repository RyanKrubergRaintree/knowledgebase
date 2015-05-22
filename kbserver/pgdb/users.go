package pgdb

import (
	"database/sql"

	"github.com/raintreeinc/knowledgebase/kbserver"
)

func (db *Database) Users() kbserver.Users { return &Users{db} }

type Users struct{ *Database }

// Create adds a new user to the database
// Groups field will be ignored
func (db *Users) Create(user kbserver.User) error {
	return db.exec(`
		INSERT INTO Users 
		(Name, Email, Description)
		VALUES ($1, $2, $3)`,
		user.Name, user.Email, user.Description)
}

func (db *Users) Delete(name string) error {
	return db.exec(`DELETE FROM Users WHERE Name = ?`, name)
}

func (db *Users) ByName(name string) (kbserver.User, error) {
	var user kbserver.User
	err := db.QueryRow(`
		SELECT
			Name, Email, Description,
			array_agg(Memberships.GroupName) as Groups
		FROM Users
		JOIN Memberships ON (Users.Name = Memberships.UserName)
		GROUP BY Users.Name
		WHERE Users.Name = $1
	`, name).Scan(&user.Name, &user.Email, &user.Description, &user.Groups)
	if err == sql.ErrNoRows {
		return user, kbserver.ErrUserNotExist
	}
	return user, nil
}

func (db *Users) List() ([]kbserver.User, error) {
	rows, err := db.Query(`
		SELECT
			Name, Email, Description,
			array_agg(Memberships.GroupName) as Groups
		FROM Users
		JOIN Memberships ON (Users.Name = Memberships.UserName)
		GROUP BY Users.Name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []kbserver.User
	for rows.Next() {
		var user kbserver.User
		rows.Scan(&user.Name, &user.Email, &user.Description, &user.Groups)
		users = append(users, user)
	}
	return users, rows.Err()
}
