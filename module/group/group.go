package group

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
		ID:          "group",
		Name:        "Group",
		Public:      true,
		Description: "Displays group information.",
	}
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{
		{
			Owner:    "group",
			Slug:     "group:groups",
			Title:    "Groups",
			Synopsis: "List of all groups.",
		},
	}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/group:groups", mod.groups).Methods("GET")
	mod.router.HandleFunc("/group:{group-id}-details", mod.details).Methods("GET")
	mod.router.HandleFunc("/group:{group-id}-members", mod.members).Methods("GET")
	mod.router.HandleFunc("/group:{group-id}", mod.pages).Methods("GET")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

var esc = html.EscapeString

func (mod *Module) details(w http.ResponseWriter, r *http.Request) {
	_, groupID, ok := mod.server.AccessGroupRead(w, r)
	if !ok {
		return
	}

	group, err := mod.server.Groups().ByID(groupID)
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

func (mod *Module) modulePages(sysId kb.Slug, w http.ResponseWriter, r *http.Request) {
	sysgroup := mod.server.Modules[sysId]
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

func (mod *Module) pages(w http.ResponseWriter, r *http.Request) {
	userID, groupID, ok := mod.server.AccessGroup(w, r)
	if !ok {
		return
	}
	if _, isSystem := mod.server.Modules[groupID]; isSystem {
		mod.modulePages(groupID, w, r)
		return
	}

	index := mod.server.IndexByUser(userID)
	entries, err := index.ByGroup(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	group, err := mod.server.Groups().ByID(groupID)
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

func (mod *Module) groups(w http.ResponseWriter, r *http.Request) {
	index, ok := mod.server.AccessIndex(w, r)
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

	if len(mod.server.Modules) > 0 {
		story.Append(kb.HTML("<h2>Modules:</h2>"))
		for _, Module := range mod.server.Modules {
			entry := Module.Info()
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

func (mod *Module) members(w http.ResponseWriter, r *http.Request) {
	_, groupID, ok := mod.server.AccessGroupWrite(w, r)
	if !ok {
		return
	}

	group, err := mod.server.Groups().ByID(groupID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	members, err := mod.server.Groups().MembersOf(groupID)
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
