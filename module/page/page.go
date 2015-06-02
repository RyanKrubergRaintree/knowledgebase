package page

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ kbserver.Module = &Module{}

type Module struct {
	server *kbserver.Server
	router *mux.Router
}

func New(server *kbserver.Server) *Module {
	mod := &Module{
		server: server,
		router: mux.NewRouter(),
	}
	mod.init()
	return mod
}

func (mod *Module) Info() kbserver.Group {
	return kbserver.Group{
		ID:          "page",
		Name:        "Page",
		Public:      true,
		Description: "Displays page listing and information.",
	}
}

func (mod *Module) Pages() []kb.PageEntry {
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

func (mod *Module) init() {
	mod.router.HandleFunc("/page:pages", mod.pages).Methods("GET")
	mod.router.HandleFunc("/page:recent-changes", mod.recentChanges).Methods("GET")
	mod.router.HandleFunc("/page:search", mod.search).Methods("GET")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

func (mod *Module) pages(w http.ResponseWriter, r *http.Request) {
	index, ok := mod.server.AccessIndex(w, r)
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

func (mod *Module) search(w http.ResponseWriter, r *http.Request) {
	index, ok := mod.server.AccessIndex(w, r)
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

func (mod *Module) recentChanges(w http.ResponseWriter, r *http.Request) {
	index, ok := mod.server.AccessIndex(w, r)
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
