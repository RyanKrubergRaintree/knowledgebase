package kbtag

import (
	"html"
	"net/http"
	"strconv"

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
		ID:          "tag",
		Name:        "Tag",
		Public:      true,
		Description: "Displays tag index.",
	}
}

func (sys *System) init() {
	sys.router.HandleFunc("/tag:tags", sys.tags).Methods("GET")
	sys.router.HandleFunc("/tag:{tag-id}", sys.pages).Methods("GET")
}

func (sys *System) Pages() []kb.PageEntry {
	return []kb.PageEntry{
		{
			Owner:    "tag",
			Slug:     "tag:tags",
			Title:    "Tags",
			Synopsis: "Listing of all tags.",
		},
	}
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.router.ServeHTTP(w, r)
}

func (sys *System) pages(w http.ResponseWriter, r *http.Request) {
	userID, ok := sys.server.AccessAuth(w, r)
	if !ok {
		return
	}
	index := sys.server.IndexByUser(userID)

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

func (sys *System) tags(w http.ResponseWriter, r *http.Request) {
	userID, ok := sys.server.AccessAuth(w, r)
	if !ok {
		return
	}
	index := sys.server.IndexByUser(userID)

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
