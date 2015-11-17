package auth

import (
	"net/http"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"
)

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

	sessions *Sessions
}

func NewServer(rules Rules, db kb.Database) *Server {
	return &Server{
		Rules:    rules,
		DB:       db,
		Provider: make(map[string]Provider),
		sessions: NewSessions(),
	}
}

func (server *Server) Verify(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	return server.sessions.GetUser(w, r)
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

	if err := server.sessions.SaveUser(w, r, user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
