package group

import (
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
		ID:          "group",
		Name:        "Group",
		Public:      true,
		Description: "Displays group information.",
	}
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{{
		Slug:     "group:groups",
		Title:    "Groups",
		Synopsis: "List of all groups.",
	}}
}

func (mod *Module) init() {
	mod.router.HandleFunc("/group:groups", mod.groups).Methods("GET")
	mod.router.HandleFunc("/group:modules", mod.modules).Methods("GET")
	mod.router.HandleFunc("/group:module-{module-id}", mod.modulePages).Methods("GET")
	mod.router.HandleFunc("/group:{group-id}-details", mod.details).Methods("GET")
	mod.router.HandleFunc("/group:{group-id}-members", mod.members).Methods("GET")
	mod.router.HandleFunc("/group:{group-id}", mod.pages).Methods("GET")
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mod.router.ServeHTTP(w, r)
}

var esc = html.EscapeString

func (mod *Module) details(w http.ResponseWriter, r *http.Request) {
	context, groupID, ok := mod.server.GroupContext(w, r, kb.Reader)
	if !ok {
		return
	}

	group, err := context.Groups().ByID(groupID)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	page := &kb.Page{
		Slug:  "group:" + groupID + "-details",
		Title: esc(group.Name) + " Details",
	}
	page.Story.Append(kb.HTML(fmt.Sprintf(`
		<p><b>Info:</b></p>
		<table>
			<tr><td>ID</td><td>%s</td></tr>
			<tr><td>Name</td><td>%s</td></tr>
			<tr><td>Public</td><td>%v</td></tr>
			<tr><td>Description</td><td>%s</td></tr>
		</table>
	`, group.ID, esc(group.Name), group.Public, esc(group.Description))))

	page.WriteResponse(w)
}

func (mod *Module) modulePages(w http.ResponseWriter, r *http.Request) {
	moduleID := kb.SlugParam(r, "module-id")
	if moduleID == "" {
		http.Error(w, "module-id missing", http.StatusBadRequest)
		return
	}

	module, ok := mod.server.Modules[moduleID]
	if !ok {
		http.NotFound(w, r)
		return
	}

	info := module.Info()
	page := &kb.Page{
		Slug:     "group:module-" + info.ID,
		Title:    "Module " + info.Name,
		Synopsis: info.Description,
	}

	entries := module.Pages()
	kb.SortPageEntriesBySlug(entries)
	page.Story = kb.StoryFromEntries(entries)
	page.Story.Prepend(kb.Paragraph(info.Description))

	page.WriteResponse(w)
}

func (mod *Module) pages(w http.ResponseWriter, r *http.Request) {
	context, groupID, ok := mod.server.GroupContext(w, r, kb.Reader)
	if !ok {
		return
	}

	info, err := context.Groups().ByID(groupID)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	entries, err := context.Index(context.ActiveUserID()).ByGroup(info.ID)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	page := &kb.Page{
		Slug:     "group:" + info.ID,
		Title:    info.Name,
		Synopsis: info.Description,
	}

	page.Story = kb.StoryFromEntries(entries)
	page.Story.Prepend(kb.Paragraph(info.Description))

	page.WriteResponse(w)
}

func (mod *Module) groups(w http.ResponseWriter, r *http.Request) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	page := &kb.Page{
		Slug:  "group:groups",
		Title: "Groups",
	}

	entries, err := index.Groups(kb.Reader)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	if len(entries) == 0 {
		page.Story.Append(kb.Paragraph("No results."))
	} else {
		for _, entry := range entries {
			page.Story.Append(kb.Entry(entry.Name, entry.Description, "group:"+entry.ID))
		}
	}

	page.WriteResponse(w)
}

func (mod *Module) modules(w http.ResponseWriter, r *http.Request) {
	_, ok := mod.server.UserContext(w, r)
	if !ok {
		return
	}

	page := &kb.Page{
		Slug:  "group:modules",
		Title: "Modules",
	}

	if len(mod.server.Modules) == 0 {
		page.Story.Append(kb.Paragraph("No results."))
	} else {
		for _, Module := range mod.server.Modules {
			entry := Module.Info()
			page.Story.Append(kb.Entry(
				esc("Module "+entry.Name),
				esc(entry.Description),
				"group:module-"+entry.ID,
			))
		}
	}

	page.WriteResponse(w)
}

func (mod *Module) members(w http.ResponseWriter, r *http.Request) {
	context, groupID, ok := mod.server.GroupContext(w, r, kb.Moderator)
	if !ok {
		return
	}

	group, err := context.Groups().ByID(groupID)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	members, err := context.Access().List(groupID)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	page := &kb.Page{
		Slug:  "group:" + groupID + "-members",
		Title: group.Name + " Members",
	}

	page.Story.Append(kb.HTML("<h2>Moderators</h2>"))
	el := "<ul>"
	for _, member := range members {
		if !member.IsGroup {
			el += "<li>" + html.EscapeString(member.Name) + "</li>"
		}
	}
	el += "</ul>"
	page.Story.Append(kb.HTML(el))

	page.Story.Append(kb.HTML("<h2>Community</h2>"))
	el = "<ul>"
	for _, member := range members {
		if member.IsGroup {
			el += "<li>" + html.EscapeString(member.Name) + " = " + string(member.Access) + "</li>"
		}
	}
	el += "</ul>"
	page.Story.Append(kb.HTML(el))

	page.WriteResponse(w)
}
