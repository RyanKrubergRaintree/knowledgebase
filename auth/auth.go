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
	provider, ok := server.Provider[providername]
	if !ok {
		http.Error(w, "Unknown authorization provider: "+providername, http.StatusBadRequest)
		return
	}

	user, err := provider.Verify(
		r.FormValue("user"),
		r.FormValue("code"))
	if err != nil {
		http.Error(w, "Unable to verify: "+err.Error(), http.StatusForbidden)
		return
	}

	if err := server.Rules.Login(user, server.DB); err != nil {
		http.Error(w, "Unknown user: "+err.Error(), http.StatusForbidden)
		return
	}

	token, err := server.Sessions.New(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&SessionInfo{
		Token: token.String(),
		User:  user,
	})
}
