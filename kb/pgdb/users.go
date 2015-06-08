package pgdb

import (
	"database/sql"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Users struct{ Context }

func (db Users) ByID(id kb.Slug) (user kb.User, err error) {
	err = db.QueryRow(`
		SELECT
			ID, Email, Name, Company, Admin,
			AuthID, AuthProvider
		FROM    Users
		WHERE   ID = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Company, &user.Admin,
		&user.AuthID, &user.AuthProvider,
	)

	if err == sql.ErrNoRows {
		return user, kb.ErrUserNotExist
	}
	return user, err
}

func (db Users) Create(user kb.User) error {
	_, err := db.Exec(`
		INSERT INTO Users(
			ID, Email, Name, Company, Admin,
			AuthID, AuthProvider
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7
		)
	`, user.ID, user.Email, user.Name, user.Company, user.Admin,
		user.AuthID, user.AuthProvider,
	)
	if dupkey(err) {
		return kb.ErrUserExists
	}
	return err
}

func (db Users) Delete(id kb.Slug) error {
	_, err := db.Exec(`DELETE FROM Users WHERE ID = $1`, id)
	return err
}

func (db Users) List() (users []kb.User, err error) {
	rows, err := db.Query(`
		SELECT
			ID, Email, Name, Company, Admin,
			AuthID, AuthProvider
		FROM    Users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user kb.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Company, &user.Admin,
			&user.AuthID, &user.AuthProvider,
		)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
