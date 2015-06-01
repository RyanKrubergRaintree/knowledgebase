package kbserver

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"

	"github.com/gorilla/mux"
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
	w.Header().Set("Allow", "GET")
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
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if _, err := server.Users().ByID(user.ID); err == ErrUserNotExist {
		err = server.Users().Create(User{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
		if err != nil {
			log.Println("Failed creating user:", err)
		}

		err = server.Groups().AddMember("community", user.ID)
		if err != nil {
			log.Println("Failed adding to community:", err)
		}
	}

	if !server.CanRead(user.ID, group) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if server.CanWrite(user.ID, group) {
		w.Header().Set("Allow", "GET, PUT, PATCH, DELETE")
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

type HandleUser func(userID kb.Slug, w http.ResponseWriter, r *http.Request)
type HandleGroup func(userID, groupID kb.Slug, w http.ResponseWriter, r *http.Request)

type Router struct {
	*Server
	*mux.Router
}

func NewRouter(server *Server) Router {
	return Router{server, mux.NewRouter()}
}

func (s *Server) AccessAuth(w http.ResponseWriter, r *http.Request) (userID kb.Slug, ok bool) {
	auth, err := s.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	return auth.ID, true
}

func (s *Server) AccessIndex(w http.ResponseWriter, r *http.Request) (index Index, ok bool) {
	var userID kb.Slug
	userID, ok = s.AccessAuth(w, r)
	if !ok {
		return
	}

	return s.Database.IndexByUser(userID), true
}

func (s *Server) AccessUserInfo(w http.ResponseWriter, r *http.Request) (auth kb.User, user User, ok bool) {
	auth, err := s.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err = s.Users().ByID(auth.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return auth, user, true
}

func (s *Server) AccessGroup(w http.ResponseWriter, r *http.Request) (userID, groupID kb.Slug, ok bool) {
	userID, ok = s.AccessAuth(w, r)
	if !ok {
		return
	}

	groupID = SlugParam(r, "group-id")
	ok = groupID != ""
	if !ok {
		http.Error(w, "group-id missing", http.StatusBadRequest)
		return
	}

	return
}

func (s *Server) AccessGroupRead(w http.ResponseWriter, r *http.Request) (userID, groupID kb.Slug, ok bool) {
	userID, groupID, ok = s.AccessGroup(w, r)
	if !ok {
		return
	}
	ok = s.CanRead(userID, groupID)
	if !ok {
		http.Error(w, "Not allowed to see this group.", http.StatusUnauthorized)
	}
	return
}

func (s *Server) AccessGroupWrite(w http.ResponseWriter, r *http.Request) (userID, groupID kb.Slug, ok bool) {
	userID, groupID, ok = s.AccessGroup(w, r)
	if !ok {
		return
	}
	ok = s.CanWrite(userID, groupID)
	if !ok {
		http.Error(w, "Not a member of this group.", http.StatusUnauthorized)
	}
	return
}

func (s *Server) AccessAdmin(w http.ResponseWriter, r *http.Request) (userID kb.Slug, ok bool) {
	userID, ok = s.AccessAuth(w, r)
	if !ok {
		return
	}

	ok = s.CanWrite(userID, "admin")
	if !ok {
		http.Error(w, "Not an administrator.", http.StatusUnauthorized)
	}
	return
}

func SlugParam(r *http.Request, name string) kb.Slug {
	v := mux.Vars(r)[name]
	if v == "" {
		return ""
	}
	return kb.Slugify(v)
}
