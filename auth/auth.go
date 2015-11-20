package auth

import (
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/raintreeinc/knowledgebase/auth/session"
	"github.com/raintreeinc/knowledgebase/kb"
)

var ErrUnauthorized = errors.New("Unauthorized")

type Provider interface {
	Boot() template.HTML
	Verify(user, pass string) (kb.User, error)
}

type Rules interface {
	Login(user kb.User, db kb.Database) error
}

type Server struct {
	Rules    Rules
	DB       kb.Database
	Provider map[string]Provider

	Sessions *session.Store
}

func NewServer(rules Rules, db kb.Database) *Server {
	return &Server{
		Rules:    rules,
		DB:       db,
		Provider: make(map[string]Provider),
		Sessions: session.NewStore(time.Hour),
	}
}

func (server *Server) params(w http.ResponseWriter, r *http.Request) (kb.User, session.Token, error) {
	token, err := session.TokenFromString(r.Header.Get("X-Auth-Token"))
	if err != nil {
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	user, ok := server.Sessions.Load(token)
	if !ok {
		return kb.User{}, token, ErrUnauthorized
	}

	return user, token, nil
}

func (server *Server) Verify(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	user, _, err := server.params(w, r)
	return user, err
}

type SessionInfo struct {
	Token string  `json:"token"`
	User  kb.User `json:"user"`

	Param map[string]string `json:"param,omitempty"`
}

func (server *Server) info(w http.ResponseWriter, r *http.Request) {
	tokstring := r.FormValue("token")
	if tokstring == "" {
		tokstring = r.Header.Get("X-Auth-Token")
	}
	token, err := session.TokenFromString(tokstring)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := server.Sessions.Load(token)
	if !ok {
		http.Error(w, "Session expired.", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&SessionInfo{
		Token: token.String(),
		User:  user,
	})
}

func (server *Server) logout(w http.ResponseWriter, r *http.Request) {
	tokstring := r.FormValue("token")
	if tokstring == "" {
		tokstring = r.Header.Get("X-Auth-Token")
	}
	token, err := session.TokenFromString(tokstring)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusBadRequest)
		return
	}

	server.Sessions.Delete(token)
}

func (server *Server) login(providername, username, pass string) (kb.User, session.Token, error) {
	provider, ok := server.Provider[providername]
	if !ok {
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	user, err := provider.Verify(username, pass)
	if err != nil {
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	if err := server.Rules.Login(user, server.DB); err != nil {
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	token, err := server.Sessions.New(user)
	if err != nil {
		return kb.User{}, session.ZeroToken, err
	}

	return user, token, nil
}

func (server *Server) InitialSession(r *http.Request) (*SessionInfo, bool) {
	auth := r.Header.Get("Authorization")

	args := strings.SplitN(auth, " ", 3)
	if len(args) != 3 {
		return nil, false
	}

	user, token, err := server.login(args[0], args[1], args[2])
	if err != nil {
		return nil, false
	}

	return &SessionInfo{
		User:  user,
		Token: token.String(),
	}, true
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/info") {
		server.info(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/logout") {
		server.logout(w, r)
		return
	}

	providername := strings.TrimPrefix(r.URL.Path, "/")
	user, token, err := server.login(
		providername,
		r.FormValue("user"),
		r.FormValue("code"),
	)

	if err != nil {
		if err == ErrUnauthorized {
			http.Error(w, "Access denied.", http.StatusForbidden)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&SessionInfo{
		Token: token.String(),
		User:  user,
	})
}
