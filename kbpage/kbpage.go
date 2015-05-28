package kbpage

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ kbserver.System = &System{}

type System struct {
	server *kbserver.Server
	router *mux.Router
}

func New(server *kbserver.Server) *System {
	sys := &System{
		server: server,
		router: mux.NewRouter(),
	}
	sys.init()
	return sys
}

func (sys *System) Info() kbserver.Group {
	return kbserver.Group{
		ID:          "page",
		Name:        "Page",
		Public:      true,
		Description: "Displays page listing and information.",
	}
}

func (sys *System) Pages() []kb.PageEntry {
	return []kb.PageEntry{
		{
			Owner:    "page",
			Slug:     "page:pages",
			Title:    "Pages",
			Synopsis: "List of all pages.",
		},
		{
			Owner:    "page",
			Slug:     "page:recent-changes",
			Title:    "Recent Changes",
			Synopsis: "Shows recently changed pages.",
		},
	}
}

func (sys *System) init() {
	sys.router.HandleFunc("/page:pages", sys.pages).Methods("GET")
	sys.router.HandleFunc("/page:recent-changes", sys.recentChanges).Methods("GET")
	sys.router.HandleFunc("/page:search", sys.search).Methods("GET")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.router.ServeHTTP(w, r)
}

func (sys *System) pages(w http.ResponseWriter, r *http.Request) {
	index, ok := sys.server.AccessIndex(w, r)
	if !ok {
		return
	}

	entries, err := index.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "page",
		Slug:  "page:pages",
		Title: "Pages",
		Story: kb.StoryFromEntries(entries),
	})
}

func (sys *System) search(w http.ResponseWriter, r *http.Request) {
	index, ok := sys.server.AccessIndex(w, r)
	if !ok {
		return
	}

	q := r.URL.Query().Get("q")
	entries, err := index.Search(q)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "page",
		Slug:  "page:search",
		Title: "Search \"" + q + "\"",
		Story: kb.StoryFromEntries(entries),
	})
}

func (sys *System) recentChanges(w http.ResponseWriter, r *http.Request) {
	index, ok := sys.server.AccessIndex(w, r)
	if !ok {
		return
	}

	entries, err := index.RecentChanges(30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "page",
		Slug:  "page:recent-changes",
		Title: "Recent Changes",
		Story: kb.StoryFromEntries(entries),
	})
}
