package kbgroup

import (
	"fmt"
	"html"
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
		ID:          "group",
		Name:        "Group",
		Public:      true,
		Description: "Displays group information.",
	}
}

func (sys *System) Pages() []kb.PageEntry {
	return []kb.PageEntry{
		{
			Owner:    "group",
			Slug:     "group:groups",
			Title:    "Groups",
			Synopsis: "List of all groups.",
		},
	}
}

func (sys *System) init() {
	sys.router.HandleFunc("/group:groups", sys.groups).Methods("GET")
	sys.router.HandleFunc("/group:{group-id}-details", sys.details).Methods("GET")
	sys.router.HandleFunc("/group:{group-id}-members", sys.members).Methods("GET")
	sys.router.HandleFunc("/group:{group-id}", sys.pages).Methods("GET")
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	sys.router.ServeHTTP(w, r)
}

var esc = html.EscapeString

func (sys *System) details(w http.ResponseWriter, r *http.Request) {
	_, groupID, ok := sys.server.AccessGroupRead(w, r)
	if !ok {
		return
	}

	group, err := sys.server.Groups().ByID(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	story := kb.Story{}

	story.Append(kb.HTML(fmt.Sprintf(`
		<p><b>Info:</b></p>
		<table>
			<tr><td>ID</td><td>%s</td></tr>
			<tr><td>Name</td><td>%s</td></tr>
			<tr><td>Public</td><td>%v</td></tr>
			<tr><td>Description</td><td>%s</td></tr>
		</table>
	`, group.ID, esc(group.Name), group.Public, esc(group.Description))))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "group",
		Slug:  "group:" + groupID + "-details",
		Title: esc(group.Name) + " Details",
		Story: story,
	})
}

func (sys *System) systemPages(sysId kb.Slug, w http.ResponseWriter, r *http.Request) {
	sysgroup := sys.server.Systems[sysId]
	group := sysgroup.Info()

	entries := sysgroup.Pages()
	kb.SortPageEntriesBySlug(entries)
	story := kb.StoryFromEntries(entries)
	story.Prepend(kb.Paragraph(group.Description))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner:    "group",
		Slug:     "group:" + sysId,
		Title:    group.Name,
		Synopsis: group.Description,
		Story:    story,
	})
}

func (sys *System) pages(w http.ResponseWriter, r *http.Request) {
	userID, groupID, ok := sys.server.AccessGroup(w, r)
	if !ok {
		return
	}
	if _, isSystem := sys.server.Systems[groupID]; isSystem {
		sys.systemPages(groupID, w, r)
		return
	}

	index := sys.server.IndexByUser(userID)
	entries, err := index.ByGroup(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, err := sys.server.Groups().ByID(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	story := kb.StoryFromEntries(entries)
	story.Prepend(kb.Paragraph(group.Description))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner:    "group",
		Slug:     "group:" + groupID,
		Title:    group.Name,
		Synopsis: group.Description,
		Story:    story,
	})
}

func (sys *System) groups(w http.ResponseWriter, r *http.Request) {
	index, ok := sys.server.AccessIndex(w, r)
	if !ok {
		return
	}

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

	if len(sys.server.Systems) > 0 {
		story.Append(kb.HTML("<h2>System Groups:</h2>"))
		for _, system := range sys.server.Systems {
			entry := system.Info()
			story.Append(kb.Entry(
				esc(entry.Name),
				esc(entry.Description),
				"group:"+entry.ID,
			))
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
	_, groupID, ok := sys.server.AccessGroupWrite(w, r)
	if !ok {
		return
	}

	group, err := sys.server.Groups().ByID(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := sys.server.Groups().MembersOf(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	story := kb.Story{}

	el := "<ul>"
	for _, member := range members {
		el += "<li>" + html.EscapeString(member.Name) + "</li>"
	}
	el += "</ul>"
	story.Append(kb.HTML(el))

	kbserver.WriteJSON(w, r, &kb.Page{
		Owner: "group",
		Slug:  "group:" + groupID + "-members",
		Title: group.Name + " Members",
		Story: story,
	})
}
