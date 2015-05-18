package kbserver

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Presenter interface {
	Present(w http.ResponseWriter, r *http.Request, tname string, data interface{}) error
}

type Server struct {
	Domain string
	Presenter
	Context
}

func New(domain, database string, presenter Presenter, context Context) *Server {
	return &Server{
		Domain:    domain,
		Presenter: presenter,
		Context:   context,
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		err := server.Present(w, r, "index.html", nil)
		if err != nil {
			log.Println(err)
		}
		return
	}

	path := strings.TrimSuffix(r.URL.Path, ".json")
	owner, slug := kb.SplitOwner(path)
	if owner == "" && strings.HasSuffix(r.Host, "."+server.Domain) {
		owner = strings.TrimSuffix(r.Host, "."+server.Domain)
	}
	if owner == "" {
		http.NotFound(w, r)
	}

	fmt.Fprintf(w, "<h1>%s</h1><h2>%s</h2>", owner, slug)

	//http.NotFound(w, r)
}
