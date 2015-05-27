package kbserver

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Sources interface {
	Include() string
}

type System interface {
	Info() Group
	Pages() []kb.PageEntry
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	Domain    string
	Templates string
	Database
	Sources
	Context

	Systems map[kb.Slug]System
}

func New(domain, templates string, db Database, sources Sources, context Context) *Server {
	return &Server{
		Domain:    domain,
		Templates: templates,
		Database:  db,
		Sources:   sources,
		Context:   context,

		Systems: make(map[kb.Slug]System),
	}
}

func (server *Server) AddSystem(system System) {
	slug := system.Info().ID
	_, exists := server.Systems[slug]
	if exists {
		panic("System " + system.Info().Name + " already exists.")
	}
	server.Systems[slug] = system
}

func tokenizeLink(link string) (owner kb.Slug, page kb.Slug) {
	if strings.HasPrefix(link, "/") {
		link = link[1:]
	}
	slug := kb.Slugify(link)

	i := strings.Index(string(slug), ":")
	if i < 0 {
		return "", slug
	}
	return slug[:i], slug
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		err := server.Present(w, r, "index.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	path := strings.TrimSuffix(r.URL.Path, ".json")

	group, slug := tokenizeLink(path)
	if group == "" {
		http.Error(w, "No owner specified", http.StatusBadRequest)
		return
	}

	if sys, ok := server.Systems[group]; ok {
		sys.ServeHTTP(w, r)
		return
	}

	user, err := server.CurrentUser(w, r)
	if err != nil || !server.CanRead(user.ID, group) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if _, err := server.Users().ByID(user.ID); err != nil {
		server.Users().Create(User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Admin: false,
		})

		server.Groups().AddMember("community", user.ID)
		server.Groups().AddMember("engineering", user.ID)
	}

	pages := server.PagesByGroup(user.ID, group)
	switch r.Method {
	case "GET":
		data, err := pages.LoadRaw(slug)
		if err != nil {
			HandleError(w, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(data)
	case "PUT":
		page, err := kb.ReadJSONPage(r.Body)
		r.Body.Close()
		if err != nil {
			HandleError(w, err)
			return
		}

		owner, _ := tokenizeLink(string(page.Slug))
		if !server.CanWrite(user.ID, owner) {
			http.Error(w, "not a member of the group", http.StatusUnauthorized)
			return
		}
		page.Owner = owner

		if page.Title == "" {
			http.Error(w, "title missing", http.StatusBadRequest)
			return
		}

		if err := pages.Create(page); err != nil {
			HandleError(w, err)
			return
		}

		WriteJSON(w, r, page)
	case "PATCH":
		action, err := kb.ReadJSONAction(r.Body)
		r.Body.Close()
		if err != nil {
			HandleError(w, err)
			return
		}

		data, _ := json.Marshal(action)
		log.Println("PATCH", string(data))

		if !server.CanWrite(user.ID, group) {
			http.Error(w, "not a member of the group", http.StatusUnauthorized)
			return
		}

		page, err := pages.Load(slug)
		if err != nil {
			HandleError(w, err)
			return
		}

		if err := page.Apply(action); err != nil {
			HandleError(w, err)
			return
		}

		if err := pages.Save(slug, page); err != nil {
			HandleError(w, err)
		}
	case "DELETE":
		http.Error(w, "not implemented", http.StatusNotImplemented)
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func WriteJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(v)
	if err != nil {
		HandleError(w, err)
		return
	}
	w.Write(data)
}

func HandleError(w http.ResponseWriter, err error) {
	switch err {
	case ErrPageExists:
		http.Error(w, err.Error(), http.StatusForbidden)
	case ErrPageNotExist:
		http.Error(w, err.Error(), http.StatusNotFound)
	case ErrUserNotAllowed:
		http.Error(w, err.Error(), http.StatusUnauthorized)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (server *Server) Present(w http.ResponseWriter, r *http.Request, tname string, data interface{}) error {
	ts, err := template.New("").Funcs(
		template.FuncMap{
			"Site": func() interface{} {
				return map[string]string{
					"ShortTitle": "KB",
					"Title":      "Knowledge Base",
					"Company":    "Raintree Systems Inc.",
				}
			},
			"User": func() kb.User {
				user, _ := server.CurrentUser(w, r)
				return user
			},
			"UserGroups": func() []string {
				user, _ := server.CurrentUser(w, r)
				info, _ := server.Users().ByID(user.ID)
				return info.Groups
			},
			"Include": func() template.HTML {
				return template.HTML(server.Sources.Include())
			},
		},
	).ParseGlob(filepath.Join(server.Templates, "*.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if err := ts.ExecuteTemplate(w, tname, data); err != nil {
		return err
	}
	return nil
}
