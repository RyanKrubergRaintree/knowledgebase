package kbtag

import (
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
		Description: "Displays tag indexes.",
	}
}

func (sys *System) init() {
	m := sys.router
	m.HandleFunc("/tag:tags", sys.tags).Methods("GET")
	m.HandleFunc("/tag:{tagid}", sys.pages).Methods("GET")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.router.ServeHTTP(w, r)
}

func (sys *System) pages(w http.ResponseWriter, r *http.Request) {
	user, err := sys.server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	index := sys.server.IndexByUser(user.ID)

	tagval := mux.Vars(r)["tagid"]
	if tagval == "" {
		http.Error(w, "Tag param is missing", http.StatusBadRequest)
		return
	}

	tag := kb.Slugify(tagval)
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
	user, err := sys.server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	index := sys.server.IndexByUser(user.ID)

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
			story.Append(kb.Entry(entry.Name, strconv.Itoa(entry.Count)+" pages", kb.Slugify("tag:"+entry.Name)))
		}
	}

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "tag",
		Slug:  "tag:tags",
		Title: "Tags",
		Story: story,
	})
}
