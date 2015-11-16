package auth

import (
	"net/http"

	"github.com/raintreeinc/knowledgebase/kb"
)

type DB struct {
	kb.Database
}

func NewDB(db kb.Database) *DB {
	return &DB{db}
}

func (db *DB) Start(w http.ResponseWriter, r *http.Request) {
	panic("unimplemented")
}

func (db *DB) Finish(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	return db.Context("admin").GuestLogin().Verify(username, password)
}

func (db *DB) Logins() (logins []kb.AuthLogin) {
	panic("unimplemented")
}
