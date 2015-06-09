package group

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/raintreeinc/knowledgebase/extra/simpleform"
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
	mod.router.HandleFunc("/group:{group-id}-members", mod.members).Methods(
		"GET",
		"ADD-USER", "REMOVE-USER",
		"ADD-COMMUNITY", "REMOVE-COMMUNITY")
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
			<tr><td>Owner</td><td>%s</td></tr>
			<tr><td>Name</td><td>%s</td></tr>
			<tr><td>Public</td><td>%v</td></tr>
			<tr><td>Description</td><td>%s</td></tr>
		</table>
	`, group.ID, group.OwnerID, esc(group.Name), group.Public, esc(group.Description))))

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

	switch r.Method {
	case "ADD-USER", "REMOVE-USER",
		"ADD-COMMUNITY", "REMOVE-COMMUNITY":

		name := strings.TrimSpace(r.FormValue("name"))
		if name == "" {
			http.Error(w, "Name not specified.", http.StatusBadRequest)
			return
		}

		id := kb.Slugify(name)
		var err error
		switch r.Method {
		case "ADD-USER":
			err = context.Access().AddUser(groupID, id)
		case "REMOVE-USER":
			err = context.Access().RemoveUser(groupID, id)
		case "ADD-COMMUNITY":
			rights := strings.TrimSpace(r.FormValue("rights"))
			if rights == "" {
				http.Error(w, "Rights not specified.", http.StatusBadRequest)
				return
			}
			err = context.Access().CommunityAdd(groupID, id, kb.Rights(rights))
		case "REMOVE-COMMUNITY":
			err = context.Access().CommunityRemove(groupID, id)
		}
		if err != nil {
			kb.WriteResult(w, err)
			return
		}

		switch r.Method {
		case "ADD-USER":
			w.Write([]byte("user added"))
		case "REMOVE-USER":
			w.Write([]byte("user removed"))
		case "ADD-COMMUNITY":
			w.Write([]byte("community added"))
		case "REMOVE-COMMUNITY":
			w.Write([]byte("community removed"))
		}

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

	page.Story.Append(simpleform.New(
		"/"+string(page.Slug), "",
		simpleform.Field("name", "Name"),
		simpleform.Button("ADD-USER", "Add"),
		simpleform.Button("REMOVE-USER", "Remove"),
	))

	el := "<ul>"
	for _, member := range members {
		if !member.IsGroup {
			el += "<li>" + html.EscapeString(member.Name) + "</li>"
		}
	}
	el += "</ul>"
	page.Story.Append(kb.HTML(el))

	page.Story.Append(kb.HTML("<h2>Community</h2>"))

	page.Story.Append(simpleform.New(
		"/"+string(page.Slug),
		"",
		simpleform.Field("name", "Name"),
		simpleform.Option("rights", []string{string(kb.Reader), string(kb.Editor), string(kb.Moderator)}),
		simpleform.Button("ADD-COMMUNITY", "Add"),
		simpleform.Button("REMOVE-COMMUNITY", "Remove"),
	))

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
