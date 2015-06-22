package pgdb

import (
	"crypto/rand"
	"crypto/subtle"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/raintreeinc/knowledgebase/kb"

	"golang.org/x/crypto/scrypt"
)

func encrypt(password string, salt []byte) ([]byte, error) {
	const sN, sR, sP, sLen = 16384, 8, 1, 32
	return scrypt.Key([]byte(password), salt, sN, sR, sP, sLen)
}

type GuestLogin struct{ Context }

func (db GuestLogin) Add(name, email, password string) error {
	const sSaltSize = 8
	var salt [sSaltSize]byte
	if _, err := rand.Read(salt[:]); err != nil {
		return err
	}
	var authCode [8]byte
	if _, err := rand.Read(authCode[:]); err != nil {
		return err
	}

	dk, err := encrypt(password, salt[:])
	if err != nil {
		return err
	}
	authID := hex.EncodeToString(authCode[:])

	_, err = db.Exec(`
		INSERT INTO
		GuestLogin (AuthID, Name, Email, Salt, DK)
		Values ($1, $2, $3, $4, $5)
	`, authID, name, email, salt[:], dk)
	return err
}

func (db GuestLogin) Verify(name, password string) (kb.User, error) {
	var salt, expect []byte
	var email, authID string
	err := db.QueryRow(`
		SELECT AuthID, Email, Salt, DK
		FROM GuestLogin
		WHERE LOWER(Name) = LOWER($1)`, name,
	).Scan(&authID, &email, &salt, &expect)

	if err == sql.ErrNoRows {
		return kb.User{}, kb.ErrUserNotExist
	}
	if err != nil {
		return kb.User{}, err
	}

	dk, err := encrypt(password, salt)
	if err != nil {
		return kb.User{}, err
	}

	if subtle.ConstantTimeCompare(dk, expect) != 1 {
		return kb.User{}, errors.New("Invalid password for " + name)
	}

	return kb.User{
		AuthID:       authID,
		AuthProvider: "guest",

		ID:    kb.Slugify(name),
		Email: email,
		Name:  name,

		MaxAccess: kb.Reader,
	}, nil
}
