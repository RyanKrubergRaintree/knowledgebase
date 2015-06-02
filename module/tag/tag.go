package tag

import (
	"html"
	"net/http"
	"strconv"

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
	return []kb.PageEntry{
		{
			Owner:    "tag",
			Slug:     "tag:tags",
			Title:    "Tags",
			Synopsis: "Listing of all tags.",
		},
	}
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

func (mod *Module) pages(w http.ResponseWriter, r *http.Request) {
	userID, ok := mod.server.AccessAuth(w, r)
	if !ok {
		return
	}
	index := mod.server.IndexByUser(userID)

	tag := kbserver.SlugParam(r, "tag-id")
	if tag == "" {
		http.Error(w, "tag-id missing", http.StatusBadRequest)
		return
	}

	entries, err := index.ByTag(string(tag))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "tag",
		Slug:  "tag:" + tag,
		Title: kb.SlugToTitle(tag),
		Story: kb.StoryFromEntries(entries),
	})
}

func (mod *Module) tags(w http.ResponseWriter, r *http.Request) {
	userID, ok := mod.server.AccessAuth(w, r)
	if !ok {
		return
	}
	index := mod.server.IndexByUser(userID)

	entries, err := index.Tags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	story := kb.Story{}
	if len(entries) == 0 {
		story.Append(kb.Paragraph("No results."))
	} else {
		for _, entry := range entries {
			story.Append(kb.Entry(
				html.EscapeString(entry.Name),
				strconv.Itoa(entry.Count)+" pages",
				kb.Slugify("tag:"+entry.Name)))
		}
	}

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "tag",
		Slug:  "tag:tags",
		Title: "Tags",
		Story: story,
	})
}
