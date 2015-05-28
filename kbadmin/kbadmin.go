package kbadmin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ kbserver.System = &System{}

type System struct {
	server *kbserver.Server
	router *mux.Router
}

func New(server *kbserver.Server) *System {
	sys := &System{
		server: server,
		router: mux.NewRouter(),
	}
	sys.init()
	return sys
}

func (sys *System) Info() kbserver.Group {
	return kbserver.Group{
		ID:          "admin",
		Name:        "Admin",
		Public:      false,
		Description: "Administrative pages.",
	}
}

//TODO
func (sys *System) Pages() []kb.PageEntry { return nil }

func (sys *System) init() {
	sys.router.HandleFunc("/admin:upload-help", sys.uploadhelp).Methods("GET")
	sys.router.HandleFunc("/admin:upload-help", sys.loadhelp).Methods("POST")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.router.ServeHTTP(w, r)
}

func (sys *System) uploadhelp(w http.ResponseWriter, r *http.Request) {
	_, ok := sys.server.AccessAdmin(w, r)
	if !ok {
		return
	}

	story := kb.Story{}
	story.Append(kb.HTML(`
	<from>
		<textarea></textarea>
	</form>
	`))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "admin",
		Slug:  "admin:upload-help",
		Title: "Upload Help",
		Story: story,
	})
}

func (sys *System) loadhelp(w http.ResponseWriter, r *http.Request) {

}
