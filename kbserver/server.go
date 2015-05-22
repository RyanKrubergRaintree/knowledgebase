package kbserver

import (
	"log"
	"net/http"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Presenter interface {
	Present(w http.ResponseWriter, r *http.Request, tname string, data interface{}) error
}

type Server struct {
	Domain string
	Database
	Presenter
	Context
}

func New(domain string, db Database, presenter Presenter, context Context) *Server {
	return &Server{
		Domain:    domain,
		Database:  db,
		Presenter: presenter,
		Context:   context,
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		err := server.Present(w, r, "index.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	path := strings.TrimSuffix(r.URL.Path, ".json")
	owner, slug := kb.SplitOwner(path)
	if owner == "" && strings.HasSuffix(r.Host, "."+server.Domain) {
		owner = strings.TrimSuffix(r.Host, "."+server.Domain)
	}
	if owner == "" {
		http.NotFound(w, r)
	}

	user, err := server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pages := server.PagesByGroup(user.Name, owner)
	data, err := pages.LoadRaw(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
