package admin

import (
	"net/http"

	"github.com/raintreeinc/knowledgebase/farm"
)

type Server struct {
	Renderer farm.Renderer

	Database string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		s.Renderer.Render(w, "admin.html", nil)
	case "/updatehelp":
		s.updateHelp(w, r)
	default:
		http.NotFound(w, r)
	}
}
