package tag

import (
	"html"
	"net/http"
	"strconv"

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
		ID:          "tag",
		Name:        "Tag",
		Public:      true,
		Description: "Displays tag index.",
	}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/tag:tags", mod.tags).Methods("GET")
	mod.router.HandleFunc("/tag:{tag-id}", mod.pages).Methods("GET")
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{{
		Slug:     "tag:tags",
		Title:    "Tags",
		Synopsis: "Listing of all tags.",
	}}
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

func (mod *Module) pages(w http.ResponseWriter, r *http.Request) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	tag := kb.SlugParam(r, "tag-id")
	if tag == "" {
		http.Error(w, "tag-id missing", http.StatusBadRequest)
		return
	}

	entries, err := index.ByTag(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &kb.Page{
		Slug:  "tag:" + tag,
		Title: kb.SlugToTitle(tag),
		Story: kb.StoryFromEntries(entries),
	}

	page.WriteResponse(w)
}

func (mod *Module) tags(w http.ResponseWriter, r *http.Request) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	entries, err := index.Tags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &kb.Page{
		Slug:  "tag:tags",
		Title: "Tags",
	}

	if len(entries) == 0 {
		page.Story.Append(kb.Paragraph("No results."))
	} else {
		for _, entry := range entries {
			page.Story.Append(kb.Entry(
				html.EscapeString(entry.Name),
				strconv.Itoa(entry.Count)+" pages",
				kb.Slugify("tag:"+entry.Name)))
		}
	}

	page.WriteResponse(w)
}
