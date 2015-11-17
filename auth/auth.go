package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/raintreeinc/knowledgebase/auth/session"
	"github.com/raintreeinc/knowledgebase/kb"
)

var ErrUnauthorized = errors.New("Unauthorized")

type Provider interface {
	Verify(user, pass string) (kb.User, error)
}

type Rules interface {
	Login(user kb.User, db kb.Database) error
}

type Server struct {
	Rules    Rules
	DB       kb.Database
	Provider map[string]Provider

	sessions *session.Store
}

func NewServer(rules Rules, db kb.Database) *Server {
	return &Server{
		Rules:    rules,
		DB:       db,
		Provider: make(map[string]Provider),
		sessions: session.NewStore(time.Hour),
	}
}

func (server *Server) Verify(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	token, err := session.TokenFromString(r.Header.Get("X-Auth-Token"))
	if err != nil {
		return kb.User{}, ErrUnauthorized
	}

	user, ok := server.sessions.Load(token)
	if !ok {
		return kb.User{}, ErrUnauthorized
	}

	return user, nil
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	providername := strings.TrimPrefix(r.URL.Path, "/")

	provider, ok := server.Provider[providername]
	if !ok {
		http.Error(w, "Unknown authorization provider: "+providername, http.StatusNotFound)
		return
	}

	user, err := provider.Verify(
		r.FormValue("user"),
		r.FormValue("code"))
	if err != nil {
		http.Error(w, "Unable to verify: "+err.Error(), http.StatusUnauthorized)
		return
	}

	if err := server.Rules.Login(user, server.DB); err != nil {
		http.Error(w, "Unknown user: "+err.Error(), http.StatusUnauthorized)
		return
	}

	token, err := server.sessions.New(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var result struct {
		Token string  `json:"token"`
		User  kb.User `json:"user"`
	}

	result.Token = token.String()
	result.User = user

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&result)
}
