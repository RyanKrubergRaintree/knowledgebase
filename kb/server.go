package kb

import (
	"log"
	"net/http"
)

type Presenter interface {
	Present(w http.ResponseWriter, r *http.Request, tname string, data interface{}) error
}

type Server struct {
	Domain string
	Presenter
}

func NewServer(domain, database string, presenter Presenter) *Server {
	return &Server{
		Domain:    domain,
		Presenter: presenter,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		err := s.Present(w, r, "index.html", map[string]interface{}{
			"SiteTitle": "KB",
		})
		if err != nil {
			log.Println(err)
		}
		return
	}

	http.NotFound(w, r)
}
