package kbserver

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gorilla/sessions"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Context interface {
	Login(w http.ResponseWriter, r *http.Request, u kb.User) error
	Logout(w http.ResponseWriter, r *http.Request)

	CurrentUser(w http.ResponseWriter, r *http.Request) (kb.User, error)

	Add(w http.ResponseWriter, r *http.Request, key, value string)
	Once(w http.ResponseWriter, r *http.Request, key string) string
}

type context struct {
	store sessions.Store
}

func NewContext(store sessions.Store) Context {
	return &context{store}
}

func (ctx *context) load(r *http.Request) (*sessions.Session, error) {
	s, err := ctx.store.Get(r, "context")
	s.Options.Path = "/"
	return s, err
}

func (ctx *context) Login(w http.ResponseWriter, r *http.Request, user kb.User) error {
	s, _ := ctx.load(r)

	//TODO: validate user in DB
	if !strings.HasSuffix(user.Email, "@raintreeinc.com") {
		return errors.New("Only allowed for users in @raintreinc.com.")
	}

	s.Values["user"] = user
	s.Save(r, w)
	return nil
}

func (ctx *context) Logout(w http.ResponseWriter, r *http.Request) {
	s, _ := ctx.load(r)
	for key := range s.Values {
		delete(s.Values, key)
	}
	s.Save(r, w)
}

func (ctx *context) CurrentUser(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	s, err := ctx.load(r)
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

func (ctx *context) Add(w http.ResponseWriter, r *http.Request, key, value string) {
	s, _ := ctx.load(r)
	s.Values[key] = value
	s.Save(r, w)
}

func (ctx *context) Once(w http.ResponseWriter, r *http.Request, key string) string {
	s, err := ctx.load(r)
	if err != nil {
		return ""
	}

	v, ok := s.Values[key]
	if !ok {
		return ""
	}

	x, ok := v.(string)
	if !ok {
		return ""
	}

	delete(s.Values, key)
	s.Save(r, w)

	return x
}
