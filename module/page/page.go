package page

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/kb"
)

var _ kb.Module = &Module{}

type Module struct {
	server *kb.Server
	router *mux.Router
}

func New(server *kb.Server) *Module {
	mod := &Module{
		server: server,
		router: mux.NewRouter(),
	}
	mod.init()
	return mod
}

func (mod *Module) Info() kb.Group {
	return kb.Group{
		ID:          "page",
		Name:        "Page",
		Public:      true,
		Description: "Displays page listing and information.",
	}
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{{
		Slug:     "page=pages",
		Title:    "Pages",
		Synopsis: "List of all pages.",
	}, {
		Slug:     "page=recent-changes",
		Title:    "Recent Changes",
		Synopsis: "Shows recently changed pages.",
	}}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/page=pages", mod.pages).Methods("GET")
	mod.router.HandleFunc("/page=recent-changes", mod.recentChanges).Methods("GET")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

func (mod *Module) pages(w http.ResponseWriter, r *http.Request) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	page := &kb.Page{
		Slug:  "page=pages",
		Title: "Pages",
	}

	entries, err := index.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page.Story = kb.StoryFromEntries(entries)
	page.WriteResponse(w)
}

func (mod *Module) recentChanges(w http.ResponseWriter, r *http.Request) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	entries, err := index.RecentChanges()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &kb.Page{
		Slug:  "page=recent-changes",
		Title: "Recent Changes",
		Story: kb.StoryFromEntries(entries),
	}
	page.WriteResponse(w)
}
