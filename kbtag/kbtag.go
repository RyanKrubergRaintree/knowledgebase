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
	Server *kbserver.Server
	Router *mux.Router
}

func New(server *kbserver.Server) *System {
	sys := &System{
		Server: server,
		Router: mux.NewRouter(),
	}
	sys.init()
	return sys
}

func (sys *System) Name() string { return "Tag" }

func (sys *System) init() {
	m := sys.Router
	m.HandleFunc("/tag:tags", sys.tags).Methods("GET")
	m.HandleFunc("/tag:{tagid}", sys.pages).Methods("GET")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.Router.ServeHTTP(w, r)
}

func (sys *System) pages(w http.ResponseWriter, r *http.Request) {
	user, err := sys.Server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	index := sys.Server.IndexByUser(user.ID)

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
	user, err := sys.Server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	index := sys.Server.IndexByUser(user.ID)

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
