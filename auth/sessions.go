package auth

import (
	"crypto/rand"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/raintreeinc/knowledgebase/kb"
)

type Sessions struct {
	sessions.Store
}

func NewSessions() *Sessions {
	appkey := os.Getenv("APPKEY")
	if appkey == "" {
		var key [32]byte
		_, err := rand.Read(key[:])
		if err != nil {
			panic(err)
		}
		appkey = string(key[:])
	}
	store := sessions.NewFilesystemStore("", []byte(appkey))

	return &Sessions{store}
}

func (store *Sessions) load(r *http.Request) (*sessions.Session, error) {
	s, err := store.Get(r, "context")
	s.Options.Path = "/"
	return s, err
}

func (store *Sessions) SaveUser(w http.ResponseWriter, r *http.Request, user kb.User) error {
	s, _ := store.load(r)
	s.Values["user"] = user
	s.Save(r, w)
	return nil
}

func (store *Sessions) GetUser(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	s, err := store.load(r)
	if err != nil {
		return kb.User{}, err
	}

	v, ok := s.Values["user"]
	if !ok {
		return kb.User{}, errors.New("user session missing")
	}

	user, ok := v.(kb.User)
	if !ok {
		delete(s.Values, "user")
		s.Save(r, w)
		return kb.User{}, errors.New("invalid type in session")
	}

	return user, nil
}

func (store *Sessions) ClearUser(w http.ResponseWriter, r *http.Request) error {
	s, _ := store.load(r)
	for key := range s.Values {
		delete(s.Values, key)
	}
	s.Save(r, w)
	return nil
}
