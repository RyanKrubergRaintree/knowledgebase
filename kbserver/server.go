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

type System interface {
	Name() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	Domain string
	Database
	Presenter
	Context

	Systems map[kb.Slug]System
}

func New(domain string, db Database, presenter Presenter, context Context) *Server {
	return &Server{
		Domain:    domain,
		Database:  db,
		Presenter: presenter,
		Context:   context,

		Systems: make(map[kb.Slug]System),
	}
}

func (server *Server) AddSystem(system System) {
	slug := kb.Slugify(system.Name())
	_, exists := server.Systems[slug]
	if exists {
		panic("System " + system.Name() + " already exists.")
	}
	server.Systems[slug] = system
}

func tokenizeLink(link string) (owner kb.Slug, page kb.Slug) {
	if strings.HasPrefix(link, "/") {
		link = link[1:]
	}
	slug := kb.Slugify(link)

	i := strings.LastIndex(string(slug), ":")
	if i < 0 {
		return "", slug
	}
	return slug[:i], slug
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		err := server.Present(w, r, "index.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	path := strings.TrimSuffix(r.URL.Path, ".json")

	group, slug := tokenizeLink(path)
	if group == "" {
		http.Error(w, "No owner specified", http.StatusBadRequest)
		return
	}

	if sys, ok := server.Systems[group]; ok {
		sys.ServeHTTP(w, r)
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

func WriteJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
