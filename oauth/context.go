package oauth

import (
	"errors"
	"net/http"

	"github.com/raintreeinc/knowledgebase/kb"

	"github.com/gorilla/sessions"
)

var _ kb.Context = &Context{}

type Context struct {
	Domain      string
	LoginURL    string
	CallbackURL string
	Sessions    sessions.Store
}

func (ctx *Context) session(r *http.Request) (*sessions.Session, error) {
	session, err := ctx.Sessions.Get(r, "session")
	session.Options.Domain = "." + ctx.Domain
	session.Options.Path = "/"
	return session, err
}

func (ctx *Context) GetUser(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	session, err := ctx.session(r)
	if err != nil {
		return kb.User{}, err
	}

	user, ok := session.Values["user"]
	if !ok {
		return kb.User{}, errors.New("user session missing")
	}

	u, ok := user.(kb.User)
	if !ok {
		return kb.User{}, errors.New("invalid type in session")
	}

	return u, nil
}

func (ctx *Context) LoggedIn(w http.ResponseWriter, r *http.Request) bool {
	_, err := ctx.GetUser(w, r)
	return err == nil
}

func (ctx *Context) LogOut(w http.ResponseWriter, r *http.Request) {
	session, err := ctx.session(r)
	if err != nil {
		return
	}
	for key := range session.Values {
		delete(session.Values, key)
	}
	session.Save(r, w)
}
