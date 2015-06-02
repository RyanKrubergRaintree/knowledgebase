package user

import (
	"fmt"
	"html"
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
		ID:          "user",
		Name:        "User",
		Public:      true,
		Description: "Displays user information.",
	}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/user:current", mod.current).Methods("GET")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{
		{
			Owner:    "user",
			Slug:     "user:current",
			Title:    "Current",
			Synopsis: "Information about the current user.",
		},
	}
}

var esc = html.EscapeString

func (mod *Module) current(w http.ResponseWriter, r *http.Request) {
	auth, user, ok := mod.server.AccessUserInfo(w, r)
	if !ok {
		return
	}

	story := kb.Story{}

	story.Append(kb.HTML(fmt.Sprintf(`
		<p><b>Info:</b></p>
		<table>
			<tr><td>ID</td><td>%v</td></tr>
			<tr><td>Name</td><td>%v</td></tr>
			<tr><td>Email</td><td>%v</td></tr>
		</table>
	`, user.ID, esc(user.Name), esc(user.Email))))

	el := "<p><b>Member of:</b></p><ul>"
	for _, group := range user.Groups {
		el += "<li><a href='group:" + esc(group) + "'>" + esc(group) + "</a></li>"
	}
	el += "</ul>"

	story.Append(kb.HTML(el))
	story.Append(kb.HTML(fmt.Sprintf(`
		<p><b>Authentication:</b></p>
		<table>
			<tr><td>AuthID</td><td>%s</td></tr>
			<tr><td>ID</td><td>%s</td></tr>
			<tr><td>Email</td><td>%s</td></tr>
			<tr><td>Provider</td><td>%s</td></tr>
		</table>
	`, esc(auth.AuthID), auth.ID, esc(auth.Email), esc(auth.Provider))))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner:    "user",
		Slug:     "user:current",
		Title:    "Current",
		Synopsis: "Information about the current user.",
		Story:    story,
	})
}
