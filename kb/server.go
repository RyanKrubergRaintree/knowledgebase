package kb

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Client interface {
	Include(version string) string
}

type Module interface {
	Info() Group
	Pages() []PageEntry
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Rules interface {
	Login(user User, db Database) error
}

type ServerInfo struct {
	Domain string

	ShortTitle string
	Title      string
	Company    string

	Version string
}

type AuthLogin struct{ URL, Name, Icon string }
type Auth interface {
	Logins() []AuthLogin

	Start(w http.ResponseWriter, r *http.Request)
	Finish(w http.ResponseWriter, r *http.Request) (User, error)
}

type Server struct {
	ServerInfo
	Templates string

	Auth Auth
	Database
	Client
	Rules Rules

	Modules map[Slug]Module
}

func NewServer(info ServerInfo, auth Auth, client Client, database Database) *Server {
	return &Server{
		ServerInfo: info,
		Templates:  "templates",

		Auth:     auth,
		Database: database,
		Client:   client,

		Modules: make(map[Slug]Module),
	}
}

func (server *Server) AddModule(module Module) {
	info := module.Info()
	slug := info.ID
	_, exists := server.Modules[slug]
	if exists {
		panic("Module " + info.Name + " already exists.")
	}
	server.Modules[slug] = module
}

func (server *Server) finishLogin(w http.ResponseWriter, r *http.Request) {
	user, err := server.Auth.Finish(w, r)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/system/auth/forbidden", http.StatusFound)
		return
	}

	if server.Rules != nil {
		err = server.Rules.Login(user, server.Database)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/system/auth/forbidden", http.StatusFound)
			return
		}
	}

	server.Sessions().SaveUser(w, r, user)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (server *Server) login(w http.ResponseWriter, r *http.Request) (User, bool) {
	if strings.HasPrefix(r.URL.Path, "/system/auth/") {
		switch r.URL.Path {
		case "/system/auth/login":
			server.Sessions().ClearUser(w, r)
			server.Present(w, r, "login.html", map[string]interface{}{
				"Logins": server.Auth.Logins(),
			})
		case "/system/auth/forbidden":
			server.Sessions().ClearUser(w, r)
			server.Present(w, r, "forbidden.html", nil)
		case "/system/auth/logout":
			server.Sessions().ClearUser(w, r)
			http.Redirect(w, r, "/", http.StatusFound)
		default:
			if strings.HasPrefix(r.URL.Path, "/system/auth/provider/") {
				server.Auth.Start(w, r)
			} else if strings.HasPrefix(r.URL.Path, "/system/auth/callback/") {
				server.finishLogin(w, r)
			} else {
				http.NotFound(w, r)
			}
		}
		return User{}, false
	}

	user, err := server.Sessions().GetUser(w, r)
	if err != nil {
		server.Sessions().ClearUser(w, r)
		http.Redirect(w, r, "/system/auth/login", http.StatusFound)
		return User{}, false
	}

	return user, true
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user, ok := server.login(w, r)
	if !ok {
		return
	}

	if r.URL.Path == "/" {
		server.Present(w, r, "index.html", nil)
		return
	}

	groupID, pageID := TokenizeLink(r.URL.Path)
	if groupID == "" {
		http.Error(w, "No page owner specified:\n"+
			"page links should have format owner:page-name.",
			http.StatusBadRequest)
		return
	}

	// modules must handle everything by themselves
	if module, ok := server.Modules[groupID]; ok {
		module.ServeHTTP(w, r)
		return
	}

	context := server.Context(user.ID)
	rights := context.Access().Rights(groupID, user.ID)
	var allowedMethods []string

	switch rights {
	case Blocked:
		http.Error(w, "Not enough rights to view this content.", http.StatusUnauthorized)
		return
	case Reader:
		allowedMethods = []string{"GET"}
	case Editor:
		allowedMethods = []string{"GET", "PATCH", "PUT"}
	case Moderator:
		allowedMethods = []string{"GET", "PATCH", "PUT", "OVERWRITE", "DELETE"}
	default:
		log.Println("Invalid rights returned for user %s got %d.", user.ID, rights)
		http.Error(w, "Invalid rights.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Allow", strings.Join(allowedMethods, ","))
	if !allowed(r.Method, allowedMethods) {
		http.Error(w, "Method "+r.Method+" not allowed.", http.StatusMethodNotAllowed)
		return
	}

	pages := context.Pages(groupID)
	switch r.Method {

	// reading a page
	case "GET":
		data, err := pages.LoadRaw(pageID)
		if err != nil {
			WriteResult(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(data)

	// creating/overwrite a page
	case "PUT", "OVERWRITE":
		version, err := getExpectedVersion(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		page, err := ReadJSONPage(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON content: %s", err), http.StatusBadRequest)
			return
		}

		pageOwner, _ := TokenizeLink(string(page.Slug))
		if pageOwner != groupID {
			http.Error(w, "Invalid parameters specified.", http.StatusBadRequest)
			return
		}

		if page.Title == "" {
			http.Error(w, "Page title is missing.", http.StatusBadRequest)
			return
		}

		if r.Method == "PUT" {
			WriteResult(w, pages.Create(page))
		} else if r.Method == "OVERWRITE" {
			WriteResult(w, pages.Overwrite(pageID, version, page))
		} else {
			panic("Invalid method")
		}

	// updating a page
	case "PATCH":
		version, err := getExpectedVersion(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		action, err := ReadJSONAction(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, fmt.Sprintf("Invalid JSON content: %s", err), http.StatusBadRequest)
			return
		}

		WriteResult(w, pages.Edit(pageID, version, action))

	// deleting a page
	case "DELETE":
		version, err := getExpectedVersion(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		WriteResult(w, pages.Delete(pageID, version))
	default:
		panic("Invalid method " + r.Method)
	}
}

func getExpectedVersion(r *http.Request) (int, error) {
	clientExpects := r.Header.Get("If-Match")
	if clientExpects != "" {
		ver, err := strconv.Atoi(clientExpects)
		if err != nil {
			return -1, errors.New("Invalid version specified")
		}
		return ver, nil
	}
	return -1, nil
}

func allowed(method string, allowedMethods []string) bool {
	for _, m := range allowedMethods {
		if m == method {
			return true
		}
	}
	return false
}

func (server *Server) UserContext(w http.ResponseWriter, r *http.Request) (Context, bool) {
	user, ok := server.login(w, r)
	if !ok {
		return nil, false
	}
	return server.Context(user.ID), true
}

func (server *Server) AdminContext(w http.ResponseWriter, r *http.Request) (Context, bool) {
	user, ok := server.login(w, r)
	if !ok {
		return nil, false
	}

	context := server.Context(user.ID)
	if !context.Access().IsAdmin(user.ID) {
		http.Error(w, "Not an administrative user.", http.StatusUnauthorized)
		return nil, false
	}
	return context, true
}

func SlugParam(r *http.Request, name string) Slug {
	v := mux.Vars(r)[name]
	if v == "" {
		return ""
	}
	return Slugify(v)
}

func (server *Server) GroupContext(w http.ResponseWriter, r *http.Request, min Rights) (Context, Slug, bool) {
	user, ok := server.login(w, r)
	if !ok {
		return nil, "", false
	}

	groupID := SlugParam(r, "group-id")
	ok = groupID != ""
	if !ok {
		http.Error(w, "group-id missing", http.StatusBadRequest)
		return nil, "", false
	}

	context := server.Context(user.ID)
	rights := context.Access().Rights(groupID, user.ID)
	if rights.Level() < min.Level() {
		http.Error(w, "Not an enough rights. You are "+string(rights)+", but need to be "+string(min)+".", http.StatusUnauthorized)
		return nil, groupID, false
	}
	return context, groupID, true
}

func (server *Server) IndexContext(w http.ResponseWriter, r *http.Request) (Context, Index, bool) {
	context, ok := server.UserContext(w, r)
	if !ok {
		return nil, nil, false
	}
	return context, context.Index(context.ActiveUserID()), true
}

func (server *Server) Present(w http.ResponseWriter, r *http.Request, tname string, data interface{}) {
	//TODO: this can be cached
	ts, err := template.New("").Funcs(
		template.FuncMap{
			"Site": func() ServerInfo { return server.ServerInfo },
			"User": func() User {
				user, _ := server.Sessions().GetUser(w, r)
				return user
			},
			"Client": func() template.HTML {
				return template.HTML(server.Client.Include(server.Version))
			},
		},
	).ParseGlob(filepath.Join(server.Templates, "*.html"))

	if err != nil {
		log.Printf("Error parsing templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, tname, data); err != nil {
		log.Printf("Error executing template: %s", err)
		return
	}
}

func WriteResult(w http.ResponseWriter, err error) {
	switch err {
	case nil:
		w.WriteHeader(http.StatusOK)
	case ErrPageExists:
		http.Error(w, err.Error(), http.StatusForbidden)
	case ErrPageNotExist:
		http.Error(w, err.Error(), http.StatusNotFound)
	default:
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (p *Page) WriteResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	return p.Write(w)
}
