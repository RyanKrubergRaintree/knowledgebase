package kbserver

import (
	"encoding/json"
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
	group, slug := kb.SplitLink(path)
	if group == "" {
		if strings.HasPrefix(path, "/index/") {
			r.URL.Path = strings.TrimPrefix(path, "/index")
			server.HandleIndex(w, r)
			return
		}
		http.NotFound(w, r)
		return
	}

	user, err := server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	pages := server.PagesByGroup(user.ID, group)
	data, err := pages.LoadRaw(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (server *Server) HandleIndex(w http.ResponseWriter, r *http.Request) {
	user, err := server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	index := server.IndexByUser(user.ID)
	switch r.URL.Path {
	case "/all":
		entries, err := index.List()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		page := &kb.Page{
			Owner: "",
			Slug:  "index/all",
			Title: "All",
			Story: kb.StoryFromEntries(entries),
		}

		data, err := json.Marshal(page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(data)
	default:
		http.NotFound(w, r)
	}
}
