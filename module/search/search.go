package search

import (
	"net/http"
	"sort"
	"strings"

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
		ID:          "search",
		Name:        "Search",
		Public:      true,
		Description: "For searching pages.",
	}
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/search=search", mod.search).Methods("GET")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

//nolint:unused // vsc does not find reference
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
	//nolint:errcheck
	page.WriteResponse(w)
}

func ImproveSearchResults(query string, xs []kb.PageEntry) {
	type stub struct {
		entry   kb.PageEntry
		quality int
	}

	stubs := make([]stub, len(xs))
	query = strings.ToLower(query)
	for i := range stubs {
		entry := &xs[i]
		stubs[i].entry = *entry
		title := strings.ToLower(entry.Title)
		stubs[i].quality = strings.Count(strings.ToLower(entry.Title), query)*100 +
			strings.Count(strings.ToLower(entry.Synopsis), query)*10
		if title == query {
			stubs[i].quality += 1000
		}
	}

	sort.SliceStable(stubs, func(i, j int) bool {
		return stubs[i].quality > stubs[j].quality
	})

	for i := range xs {
		xs[i] = stubs[i].entry
	}
}

func (mod *Module) search(w http.ResponseWriter, r *http.Request) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	q := r.URL.Query().Get("q")
	filter := r.Header.Get("X-Filter")

	var entries []kb.PageEntry
	var err error
	if filter == "" {
		entries, err = index.Search(q)
	} else {
		filter = string(kb.Slugify(filter))
		entries, err = index.SearchFilter(q, "help-", "help-"+filter)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ImproveSearchResults(q, entries)

	page := &kb.Page{
		Slug:  "search=search",
		Title: "Search \"" + q + "\"",
		Story: kb.StoryFromEntries(entries),
	}
	//nolint:errcheck
	page.WriteResponse(w)
}
