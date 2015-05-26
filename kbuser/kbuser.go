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
	m := sys.router
	m.HandleFunc("/user:{userid}", sys.userinfo).Methods("GET")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.router.ServeHTTP(w, r)
}

func (sys *System) userinfo(w http.ResponseWriter, r *http.Request) {
	userval := mux.Vars(r)["userid"]
	if userval == "" {
		http.Error(w, "user id is missing", http.StatusBadRequest)
		return
	}
	userid := kb.Slugify(userval)

	user, err := sys.server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userid != user.ID {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	info, err := sys.server.Users().ByID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	story := kb.Story{}

	//TODO: use sanitiziation
	story.Append(kb.HTML(fmt.Sprintf(`
		<p><b>Info:</b></p>
		<table>
			<tr><td>ID</td><td>%s</td></tr>
			<tr><td>Name</td><td>%s</td></tr>
			<tr><td>Email</td><td>%s</td></tr>
		</table>
	`, info.ID, info.Name, info.Email)))

	el := "<p><b>Member of:</b></p><ul>"
	for _, group := range info.Groups {
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
	`, user.AuthID, user.ID, user.Email, user.Provider)))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "user",
		Slug:  "user:" + userid,
		Title: user.Name,
		Story: story,
	})
}
