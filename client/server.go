package client

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/raintreeinc/livepkg"

	// "github.com/raintreeinc/knowledgebase/auth"
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

	bootstrap string
	assets    http.Handler
	client    http.Handler
}

func NewServer(info Info, dir string, development bool) *Server {
	return &Server{
		Info: info,

		bootstrap: filepath.Join(dir, "index.html"),
		assets: http.StripPrefix("/client/assets/",
			http.FileServer(http.Dir(filepath.Join(dir, "assets")))),
		//TODO fix this
		client: livepkg.NewServer(
			http.Dir(filepath.Join(dir, "..")),
			development,
			"/client/boot.js",
		),
	}
}

func (server *Server) index(w http.ResponseWriter, r *http.Request) {
	ts, err := template.New("").Funcs(
		template.FuncMap{
			"Site": func() Info { return server.Info },
		},
	).ParseGlob(filepath.Join("client", "index.html"))

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
		http.Redirect(w, r, "/client/assets/ico/favicon.ico", http.StatusMovedPermanently)
	case r.URL.Path == "/":
		server.index(w, r)
	case strings.HasPrefix(r.URL.Path, "/client/assets/"):
		server.assets.ServeHTTP(w, r)
	case strings.HasPrefix(r.URL.Path, "/client/"):
		server.client.ServeHTTP(w, r)
	default:
		http.Error(w, "Page not found.", http.StatusNotFound)
	}
}
