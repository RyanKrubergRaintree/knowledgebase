package kbuser

import (
	"fmt"
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
		ID:          "user",
		Name:        "User",
		Public:      true,
		Description: "Displays user information.",
	}
}

func (sys *System) init() {
	sys.router.HandleFunc("/user:current", sys.current).Methods("GET")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.router.ServeHTTP(w, r)
}

func (sys *System) Pages() []kb.PageEntry {
	return []kb.PageEntry{
		{
			Owner:    "user",
			Slug:     "user:current",
			Title:    "Current",
			Synopsis: "Information about the current user.",
		},
	}
}

func (sys *System) current(w http.ResponseWriter, r *http.Request) {
	auth, user, ok := sys.server.AccessUserInfo(w, r)
	if !ok {
		return
	}

	story := kb.Story{}

	//TODO: use sanitiziation
	story.Append(kb.HTML(fmt.Sprintf(`
		<p><b>Info:</b></p>
		<table>
			<tr><td>ID</td><td>%v</td></tr>
			<tr><td>Name</td><td>%v</td></tr>
			<tr><td>Email</td><td>%v</td></tr>
			<tr><td>Admin</td><td>%v</td></tr>
		</table>
	`, user.ID, user.Name, user.Email, user.Admin)))

	el := "<p><b>Member of:</b></p><ul>"
	for _, group := range user.Groups {
		el += "<li><a href='group:" + group + "'>" + group + "</a></li>"
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
	`, auth.AuthID, auth.ID, auth.Email, auth.Provider)))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner:    "user",
		Slug:     "user:current",
		Title:    "Current",
		Synopsis: "Information about the current user.",
		Story:    story,
	})
}
