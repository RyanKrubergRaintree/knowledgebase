package admin

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ kbserver.Module = &Module{}

type Module struct {
	server *kbserver.Server
	router *mux.Router
}

func New(server *kbserver.Server) *Module {
	mod := &Module{
		server: server,
		router: mux.NewRouter(),
	}
	mod.init()
	return mod
}

func (mod *Module) Info() kbserver.Group {
	return kbserver.Group{
		ID:          "admin",
		Name:        "Admin",
		Public:      false,
		Description: "Administrative pages.",
	}
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{
		{
			Owner:    "admin",
			Slug:     "admin:upload-help",
			Title:    "Upload Help",
			Synopsis: "Page for updating help.",
		},
	}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/admin:upload-help", mod.uploadhelp).Methods("GET")
	mod.router.HandleFunc("/admin:upload-help", mod.loadhelp).Methods("POST")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

func (mod *Module) uploadhelp(w http.ResponseWriter, r *http.Request) {
	_, ok := mod.server.AccessAdmin(w, r)
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

func (mod *Module) loadhelp(w http.ResponseWriter, r *http.Request) {

}
