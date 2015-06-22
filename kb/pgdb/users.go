package pgdb

import (
	"database/sql"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Users struct{ Context }

func (db Users) ByID(id kb.Slug) (user kb.User, err error) {
	var maxaccess string
	err = db.QueryRow(`
		SELECT
			ID, Email, Name, Company, Admin, MaxAccess,
			AuthID, AuthProvider
		FROM    Users
		WHERE   ID = $1
	`, id).Scan(
		&user.ID, &user.Email, &user.Name, &user.Company, &user.Admin, &maxaccess,
		&user.AuthID, &user.AuthProvider,
	)
	user.MaxAccess = kb.Rights(maxaccess)

	if err == sql.ErrNoRows {
		return user, kb.ErrUserNotExist
	}
	return user, err
}

func (db Users) Create(user kb.User) error {
	_, err := db.Exec(`
		INSERT INTO Users(
			ID, Email, Name, Company, Admin, MaxAccess,
			AuthID, AuthProvider
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8
		)
	`, user.ID, user.Email, user.Name, user.Company, user.Admin, string(user.MaxAccess),
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
			ID, Email, Name, Company, Admin, MaxAccess,
			AuthID, AuthProvider
		FROM    Users
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var maxaccess string
		var user kb.User
		err := rows.Scan(
			&user.ID, &user.Email, &user.Name, &user.Company, &user.Admin, &maxaccess,
			&user.AuthID, &user.AuthProvider,
		)
		user.MaxAccess = kb.Rights(maxaccess)
		if err != nil {
			return users, err
		}
		users = append(users, user)
	}
	return users, nil
}
