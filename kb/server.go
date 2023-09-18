package kb

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type Auth interface {
	Verify(w http.ResponseWriter, r *http.Request) (User, error)
}

type Module interface {
	Info() Group
	Pages() []PageEntry
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	Auth Auth
	Database
	Modules map[Slug]Module
}

func NewServer(auth Auth, database Database) *Server {
	return &Server{
		Auth:     auth,
		Database: database,
		Modules:  make(map[Slug]Module),
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

func (server *Server) login(w http.ResponseWriter, r *http.Request) (User, bool) {
	user, err := server.Auth.Verify(w, r)
	if err != nil {
		w.Header().Add("WWW-Authenticate", "X-Auth-Token")
		//TODO: serve correct error message
		http.Error(w, "Session expired!", http.StatusUnauthorized)
		return User{}, false
	}
	return user, true
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	AddCommonHeaders(w)

	// handle pre-flight request for LMS
	// todo: configuration file for whitelisting origins?
	if (*r).Method == "OPTIONS" && strings.HasPrefix((*r).RequestURI, "/lms=/uploadImage/") {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "X-Auth-Token")
		return
	}

	user, ok := server.login(w, r)
	if !ok {
		return
	}

	groupID, pageID := TokenizeLink(r.URL.Path)
	if groupID == "" {
		http.Error(w, "No page owner specified:\n"+
			"page links should have format owner=page-name.",
			http.StatusBadRequest)
		return
	}

	// disable caching of pages
	w.Header().Add("Cache-Control", "no-cache")
	w.Header().Add("Pragma", "no-cache")

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
		http.Error(w, "Not enough rights to view this content.", http.StatusForbidden)
		return
	case Reader:
		allowedMethods = []string{"GET"}
	case Editor:
		allowedMethods = []string{"GET", "POST", "PUT"}
	case Moderator:
		allowedMethods = []string{"GET", "POST", "PUT", "OVERWRITE", "DELETE"}
	default:
		log.Println("Invalid rights returned for user %s got %d.", user.ID, rights)
		http.Error(w, "Invalid rights.", http.StatusInternalServerError)
		return
	}

	requestedVersionStr := r.URL.Query().Get("history")
	versionedRequest := false
	requestedVersion := -1
	if requestedVersionStr != "" {
		versionedRequest = true
		if v, err := strconv.Atoi(requestedVersionStr); err == nil {
			requestedVersion = v
		}

		if rights == Moderator {
			allowedMethods = []string{"GET"}
		} else {
			allowedMethods = []string{}
		}
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
		if versionedRequest {
			if requestedVersionStr == "all" {
				entries, err := pages.History(pageID)
				if err != nil {
					WriteResult(w, err)
					return
				}

				_, title, _ := TokenizeLink3(string(pageID))
				page := &Page{
					Slug:  pageID,
					Title: SlugToTitle(title) + " (history)",
					Story: StoryFromEntries(entries),
				}
				page.WriteResponse(w)
			} else {
				data, err := pages.LoadRawVersion(pageID, requestedVersion)
				if err != nil {
					WriteResult(w, err)
					return
				}
				// TODO: modify header

				w.Header().Set("Content-Type", "application/json")
				w.Write(data)
			}
		} else {
			data, err := pages.LoadRaw(pageID)
			if err != nil {
				WriteResult(w, err)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(data)
		}

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
	case "POST":
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

func AddCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Frame-Options", "SAMEORIGIN")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
}

func getNonce() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Failed to generate random bytes: %v", err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

func AddCSPHeader(w http.ResponseWriter) string {
	// sha256 is for https://apis.google.com/js/platform.js
	// 'unsafe-inline' is present for back-comp reasons (ignored by browsers supporting nonces/hashes).
    CSPTemplate := `
		default-src
			'self';
		script-src
			'self'
			'unsafe-inline'
			'nonce-%s'
			'sha256-0LjTTmOvpWMJbo1V4agDu9F+Lhv28WhMGI6o7CJMsVI='
			https://*.gstatic.com
			*.google-analytics.com
			*.google.com
			*.googleapis.com
			*.apis.google.com;
		connect-src
			'self'
			*.google-analytics.com
			*.google.com
			*.googleapis.com;
		frame-src
			'self'
			*.google.com;
		font-src
			'self'
			fonts.gstatic.com;
		style-src
			'self'
			'unsafe-inline'
			fonts.googleapis.com;
		img-src
			'self'
			'unsafe-inline'
			data:
			https://*;
		media-src
			* blob:;
		base-uri
			'none';
		object-src
			'none';
		`
	nonce := getNonce()
	w.Header().Set("Content-Security-Policy", fmt.Sprintf(CSPTemplate, nonce))

	return nonce
}


