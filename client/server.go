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

	bootstrap string
	assets    http.Handler
	client    http.Handler
}

func NewServer(info Info, login *auth.Server, dir string, development bool) *Server {
	return &Server{
		Info:  info,
		Login: login,

		bootstrap: filepath.Join(dir, "index.html"),
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
	ts, err := template.New("").Funcs(
		template.FuncMap{
			"Site": func() Info { return server.Info },
			"LoginProviders": func() template.JS {
				info := map[string]interface{}{}
				for name, provider := range server.Login.Provider {
					info[name] = provider.Info()
				}

				data, _ := json.Marshal(info)
				return template.JS(data)
			},
		},
	).ParseGlob(server.bootstrap)

	if err != nil {
		log.Printf("Error parsing templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := ts.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Error executing template: %s", err)
		return
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/favicon.ico":
		http.Redirect(w, r, "/assets/ico/favicon.ico", http.StatusMovedPermanently)
	case r.URL.Path == "/":
		server.index(w, r)
	case strings.HasPrefix(r.URL.Path, "/assets/"):
		server.assets.ServeHTTP(w, r)
	default:
		server.client.ServeHTTP(w, r)
	}
}
