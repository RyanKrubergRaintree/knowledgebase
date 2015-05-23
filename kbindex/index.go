package kbindex

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

func (sys *System) Name() string { return "Index" }

func (sys *System) init() {
	m := sys.Router
	m.HandleFunc("/index:all", sys.withIndex(sys.indexAll)).Methods("GET")
	m.HandleFunc("/index:tags", sys.withIndex(sys.indexTags)).Methods("GET")
	m.HandleFunc("/index:tag/{tag}", sys.withIndex(sys.indexTag)).Methods("GET")
	m.HandleFunc("/index:groups", sys.withIndex(sys.indexGroups)).Methods("GET")
	m.HandleFunc("/index:group/{group}", sys.withIndex(sys.indexGroup)).Methods("GET")
	m.HandleFunc("/index:recent-changes", sys.withIndex(sys.indexRecent)).Methods("GET")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.Router.ServeHTTP(w, r)
}

type IndexHandlerFunc func(http.ResponseWriter, *http.Request, kbserver.Index)

func (sys *System) withIndex(fn IndexHandlerFunc) http.HandlerFunc {
	server := sys.Server
	return func(w http.ResponseWriter, r *http.Request) {
		user, err := server.CurrentUser(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		index := server.IndexByUser(user.ID)
		fn(w, r, index)
	}
}

func (sys *System) indexAll(w http.ResponseWriter, r *http.Request, index kbserver.Index) {
	entries, err := index.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "index",
		Slug:  "index:all",
		Title: "All",
		Story: kb.StoryFromEntries(entries),
	})
}

func (sys *System) indexTags(w http.ResponseWriter, r *http.Request, index kbserver.Index) {
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
			story.Append(kb.Entry(entry.Name, strconv.Itoa(entry.Count)+" pages", kb.Slugify("index:tag/"+entry.Name)))
		}
	}

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "index",
		Slug:  "index:tags",
		Title: "Tags",
		Story: story,
	})
}

func (sys *System) indexTag(w http.ResponseWriter, r *http.Request, index kbserver.Index) {
	tag := mux.Vars(r)["tag"]
	if tag == "" {
		http.Error(w, "Tag param is missing", http.StatusBadRequest)
		return
	}

	entries, err := index.ByTag(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "index",
		Slug:  kb.Slugify("index:tag/" + tag),
		Title: tag,
		Story: kb.StoryFromEntries(entries),
	})
}

func (sys *System) indexGroups(w http.ResponseWriter, r *http.Request, index kbserver.Index) {
	entries, err := index.Groups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	story := kb.Story{}
	if len(entries) == 0 {
		story.Append(kb.Paragraph("No results."))
	} else {
		for _, entry := range entries {
			story.Append(kb.Entry(entry.Name, entry.Description, "index:group/"+entry.ID))
		}
	}

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "index",
		Slug:  "index:groups",
		Title: "Groups",
		Story: story,
	})
}

func (sys *System) indexGroup(w http.ResponseWriter, r *http.Request, index kbserver.Index) {
	server := sys.Server
	groupval := mux.Vars(r)["group"]
	if groupval == "" {
		http.Error(w, "Group param is missing", http.StatusBadRequest)
		return
	}
	group := kb.Slugify(groupval)

	entries, err := index.ByGroup(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	info, err := server.Groups().ByID(group)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner:    "index",
		Slug:     "index:group/" + group,
		Title:    info.Name,
		Synopsis: info.Description,
		Story:    kb.StoryFromEntries(entries),
	})
}

func (sys *System) indexRecent(w http.ResponseWriter, r *http.Request, index kbserver.Index) {
	entries, err := index.RecentChanges(30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "index",
		Slug:  "index:recent-changes",
		Title: "Recent Changes",
		Story: kb.StoryFromEntries(entries),
	})
}
