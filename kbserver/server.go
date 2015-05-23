package kbserver

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/kb"
)

type Presenter interface {
	Present(w http.ResponseWriter, r *http.Request, tname string, data interface{}) error
}

type Server struct {
	Domain string
	Database
	Presenter
	Context

	Router *mux.Router
}

func New(domain string, db Database, presenter Presenter, context Context) *Server {
	server := &Server{
		Domain:    domain,
		Database:  db,
		Presenter: presenter,
		Context:   context,
		Router:    mux.NewRouter(),
	}
	server.init()
	return server
}

func (server *Server) init() {
	m := server.Router
	s := server

	m.HandleFunc("/", s.main)
	m.HandleFunc("/index/all", s.withIndex(s.indexAll)).Methods("GET")
	m.HandleFunc("/index/tags", s.withIndex(s.indexTags)).Methods("GET")
	m.HandleFunc("/index/tag/{tag}", s.withIndex(s.indexTag)).Methods("GET")
	m.HandleFunc("/index/groups", s.withIndex(s.indexGroups)).Methods("GET")
	m.HandleFunc("/index/group/{group}", s.withIndex(s.indexGroup)).Methods("GET")
	m.HandleFunc("/index/recent-changes", s.withIndex(s.indexRecent)).Methods("GET")
}

func (server *Server) main(w http.ResponseWriter, r *http.Request) {
	err := server.Present(w, r, "index.html", nil)
	if err != nil {
		log.Println(err)
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, ".json")
	group, slug := kb.SplitLink(path)
	if group == "" {
		server.Router.ServeHTTP(w, r)
		return
	}

	user, err := server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	pages := server.PagesByGroup(user.ID, group)
	data, err := pages.LoadRaw(slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

type IndexHandlerFunc func(http.ResponseWriter, *http.Request, Index)

func (server *Server) withIndex(fn IndexHandlerFunc) http.HandlerFunc {
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

func (server *Server) indexAll(w http.ResponseWriter, r *http.Request, index Index) {
	entries, err := index.List()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, r, &kb.Page{
		Owner: "",
		Slug:  "index/all",
		Title: "All",
		Story: kb.StoryFromEntries(entries),
	})
}

func (server *Server) indexTags(w http.ResponseWriter, r *http.Request, index Index) {
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
			story.Append(kb.Entry(entry.Name, strconv.Itoa(entry.Count)+" pages", kb.Slugify("index/tag/"+entry.Name)))
		}
	}

	WriteJSON(w, r, &kb.Page{
		Owner: "",
		Slug:  "index/tags",
		Title: "Tags",
		Story: story,
	})
}

func (server *Server) indexTag(w http.ResponseWriter, r *http.Request, index Index) {
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
	WriteJSON(w, r, &kb.Page{
		Owner: "",
		Slug:  kb.Slugify("index/tag/" + tag),
		Title: tag,
		Story: kb.StoryFromEntries(entries),
	})
}

func (server *Server) indexGroups(w http.ResponseWriter, r *http.Request, index Index) {
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
			story.Append(kb.Entry(entry.Name, entry.Description, "index/group/"+entry.ID))
		}
	}

	WriteJSON(w, r, &kb.Page{
		Owner: "",
		Slug:  "index/groups",
		Title: "Groups",
		Story: story,
	})
}

func (server *Server) indexGroup(w http.ResponseWriter, r *http.Request, index Index) {
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

	WriteJSON(w, r, &kb.Page{
		Owner:    "",
		Slug:     "index/group/" + group,
		Title:    info.Name,
		Synopsis: info.Description,
		Story:    kb.StoryFromEntries(entries),
	})
}

func (server *Server) indexRecent(w http.ResponseWriter, r *http.Request, index Index) {
	entries, err := index.RecentChanges(30)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	WriteJSON(w, r, &kb.Page{
		Owner: "",
		Slug:  "index/recent-changes",
		Title: "Recent Changes",
		Story: kb.StoryFromEntries(entries),
	})
}

func WriteJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}
