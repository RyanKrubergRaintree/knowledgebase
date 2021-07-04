package client

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/raintreeinc/livepkg"

	"github.com/raintreeinc/knowledgebase/auth"
)

var (
	errTimeSkewedMessage = "\n\n" +
		"Knowledge Base authentication failed due to incorrect computer time.\n" +
		"Accuracy within to 2 minutes is required.\n\n" +
		"Tip: Enable automatic time adjustment in your computer's settings."
)

type Info struct {
	Domain string

	ShortTitle string
	Title      string
	Company    string

	TrackingID string
	Version    string
}

type Server struct {
	Info
	Login *auth.Server

	development bool
	bootstrap   string
	dir         string
	assets      http.Handler
	client      http.Handler
}

func NewServer(info Info, login *auth.Server, dir string, development bool) *Server {
	return &Server{
		Info:  info,
		Login: login,

		development: development,
		bootstrap:   filepath.Join(dir, "index.html"),
		dir:         dir,
		assets: http.StripPrefix("/assets/",
			http.FileServer(http.Dir(filepath.Join(dir, "assets")))),
		client: livepkg.NewServer(
			http.Dir(dir),
			development,
			"/boot.js",
		),
	}
}

func (server *Server) index(w http.ResponseWriter, r *http.Request) {
	initialSession, initialSessionErr := server.Login.SessionFromHeader(r)
	if initialSessionErr == auth.ErrTimeSkewed {
		http.Error(w, errTimeSkewedMessage, http.StatusUnauthorized)
		return
	}

	ts, err := template.New("").Funcs(
		template.FuncMap{
			"Development": func() bool { return server.development },
			"Site":        func() Info { return server.Info },
			"InitialSession": func() template.JS {
				if initialSessionErr != nil {
					return "null"
				}

				data, _ := json.Marshal(initialSession)
				return template.JS(data)
			},
			"LoginProviders": func() interface{} {
				return server.Login.Provider
			},
		},
	).ParseGlob(server.bootstrap)

	if err != nil {
		log.Printf("Error parsing templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-UA-Compatible", "IE=edge")

	if err := ts.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Error executing template: %s", err)
		return
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/favicon.ico":
		http.ServeFile(w, r, filepath.Join(server.dir, "assets", "ico", "favicon.ico"))
	case r.URL.Path == "/":
		server.index(w, r)
	case r.URL.Path == "/apilogin":
		server.apiLogin(w, r)
	case strings.HasPrefix(r.URL.Path, "/assets/"):
		server.assets.ServeHTTP(w, r)
	default:
		server.client.ServeHTTP(w, r)
	}
}

func (server *Server) apiLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	session, err := server.Login.SessionFromHeader(r)
	if err == auth.ErrTimeSkewed {
		http.Error(w, errTimeSkewedMessage, http.StatusUnauthorized)
		return
	}

	if session == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(session.Token)
}
