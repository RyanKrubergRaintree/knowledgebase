package kbgroup

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

func (sys *System) Name() string { return "Group" }

func (sys *System) init() {
	m := sys.router
	m.HandleFunc("/group:groups", sys.groups).Methods("GET")
	m.HandleFunc("/group:{groupid}-details", sys.info).Methods("GET")
	m.HandleFunc("/group:{groupid}-members", sys.members).Methods("GET")
	m.HandleFunc("/group:{groupid}", sys.pages).Methods("GET")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.router.ServeHTTP(w, r)
}

func (sys *System) info(w http.ResponseWriter, r *http.Request) {
	groupval := mux.Vars(r)["groupid"]
	if groupval == "" {
		http.Error(w, "user id is missing", http.StatusBadRequest)
		return
	}
	groupid := kb.Slugify(groupval)

	user, err := sys.server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !sys.server.Database.CanRead(user.ID, groupid) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	group, err := sys.server.Groups().ByID(groupid)
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
			<tr><td>Public</td><td>%v</td></tr>
			<tr><td>Description</td><td>%s</td></tr>
		</table>
	`, group.ID, group.Name, group.Public, group.Description)))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "group",
		Slug:  "group:" + groupid + "-details",
		Title: group.Name + " Details",
		Story: story,
	})
}

func (sys *System) pages(w http.ResponseWriter, r *http.Request) {
	user, err := sys.server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	index := sys.server.IndexByUser(user.ID)

	groupval := mux.Vars(r)["groupid"]
	if groupval == "" {
		http.Error(w, "Group param is missing", http.StatusBadRequest)
		return
	}
	groupid := kb.Slugify(groupval)

	entries, err := index.ByGroup(groupid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, err := sys.server.Groups().ByID(groupid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	story := kb.StoryFromEntries(entries)
	story.Prepend(kb.Paragraph(group.Description))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner:    "group",
		Slug:     "group:" + groupid,
		Title:    group.Name,
		Synopsis: group.Description,
		Story:    story,
	})
}

func (sys *System) groups(w http.ResponseWriter, r *http.Request) {
	user, err := sys.server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	index := sys.server.IndexByUser(user.ID)

	entries, err := index.Groups()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	story := kb.Story{}
	if len(entries) == 0 {
		story.Append(kb.Paragraph("No results."))
	} else {
		for _, entry := range entries {
			story.Append(kb.Entry(entry.Name, entry.Description, "group:"+entry.ID))
		}
	}

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "group",
		Slug:  "group:groups",
		Title: "Groups",
		Story: story,
	})
}

func (sys *System) members(w http.ResponseWriter, r *http.Request) {
	user, err := sys.server.CurrentUser(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	groups := sys.server.Groups()

	groupval := mux.Vars(r)["groupid"]
	if groupval == "" {
		http.Error(w, "Group param is missing", http.StatusBadRequest)
		return
	}
	groupid := kb.Slugify(groupval)

	if !sys.server.CanWrite(user.ID, groupid) {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	group, err := sys.server.Groups().ByID(groupid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := groups.MembersOf(groupid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	el := "<ul>"
	for _, member := range members {
		el += "<li>" + member.Name + "</li>"
	}
	el += "</ul>"

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "group",
		Slug:  "group:" + groupid + "-members",
		Title: group.Name + " Members",
		Story: kb.Story{kb.HTML(el)},
	})
}
