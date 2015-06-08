package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/kb"
)

var _ kb.Module = &Module{}

type Module struct {
	server *kb.Server
	router *mux.Router
}

func New(server *kb.Server) *Module {
	mod := &Module{
		server: server,
		router: mux.NewRouter(),
	}
	mod.init()
	return mod
}

func (mod *Module) Info() kb.Group {
	return kb.Group{
		ID:          "admin",
		Name:        "Admin",
		Public:      false,
		Description: "Administrative pages.",
	}
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{{
		Slug:     "admin:upload-help",
		Title:    "Upload Help",
		Synopsis: "Page for updating help.",
	}}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/admin:upload-help", mod.uploadhelp).Methods("GET")
	mod.router.HandleFunc("/admin:upload-help", mod.loadhelp).Methods("POST")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

func (mod *Module) uploadhelp(w http.ResponseWriter, r *http.Request) {
	_, ok := mod.server.AdminContext(w, r)
	if !ok {
		return
	}

	page := &kb.Page{
		Slug:  "admin:upload-help",
		Title: "Upload Help",
		Story: kb.Story{kb.HTML(`<from><textarea></textarea></form>`)},
	}

	page.WriteResponse(w)
}

func (mod *Module) loadhelp(w http.ResponseWriter, r *http.Request) {

}
