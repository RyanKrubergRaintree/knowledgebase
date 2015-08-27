package user

import (
	"encoding/json"
	"fmt"
	"html"
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
		ID:          "user",
		Name:        "User",
		Public:      true,
		Description: "Displays user information.",
	}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/user:current", mod.current).Methods("GET")
	mod.router.HandleFunc("/user:editor-groups", mod.groups).Methods("GET")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{{
		Slug:     "user:current",
		Title:    "Current",
		Synopsis: "Information about the current user.",
	}}
}

var esc = html.EscapeString

func (mod *Module) current(w http.ResponseWriter, r *http.Request) {
	context, ok := mod.server.UserContext(w, r)
	if !ok {
		return
	}

	user, err := context.Users().ByID(context.ActiveUserID())
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	page := &kb.Page{
		Slug:     "user:current",
		Title:    "Current",
		Synopsis: "Information about the current user.",
	}

	page.Story.Append(kb.HTML(fmt.Sprintf(`
		<p><b>Info:</b></p>
		<table>
			<tr><td>ID</td><td>%v</td></tr>
			<tr><td>Name</td><td>%v</td></tr>
			<tr><td>Email</td><td>%v</td></tr>
		</table>
	`, user.ID, esc(user.Name), esc(user.Email))))

	page.WriteResponse(w)
}

func (mod *Module) groups(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "application/json" {
		http.Error(w, "Accept header must be application/json", http.StatusNotAcceptable)
		return
	}

	context, ok := mod.server.UserContext(w, r)
	if !ok {
		return
	}

	groups, err := context.Index(context.ActiveUserID()).Groups(kb.Editor)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	var result struct {
		Groups []string `json:"groups"`
	}

	for _, group := range groups {
		result.Groups = append(result.Groups, group.Name)
	}

	data, err := json.Marshal(result)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
