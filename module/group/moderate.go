package group

import (
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/raintreeinc/knowledgebase/extra/simpleform"
	"github.com/raintreeinc/knowledgebase/kb"
)

var esc = html.EscapeString

func (mod *Module) moderate(w http.ResponseWriter, r *http.Request) {
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
		Slug:  "group:moderate-" + groupID + "",
		Title: "Moderate " + group.Name,
	}

	page.Story.Append(kb.HTML(fmt.Sprintf(`
		<p><b>Info:</b></p>
		<table class="tight">
			<tr><td>ID</td><td>%s</td></tr>
			<tr><td>Owner</td><td>%s</td></tr>
			<tr><td>Name</td><td>%s</td></tr>
			<tr><td>Public</td><td>%v</td></tr>
			<tr><td>Description</td><td>%s</td></tr>
		</table>
	`, group.ID, group.OwnerID, esc(group.Name), group.Public, esc(group.Description))))

	page.Story.Append(kb.HTML("<p><b>Moderators:</b></p>"))

	page.Story.Append(simpleform.New(
		"/"+string(page.Slug), "",
		simpleform.Field("name", "Name"),
		simpleform.Button("ADD-USER", "Add"),
		simpleform.Button("REMOVE-USER", "Remove"),
	))

	el := `<ul class="tight">`
	for _, member := range members {
		if !member.IsGroup {
			el += "<li>" + html.EscapeString(member.Name) + "</li>"
		}
	}
	el += "</ul>"
	page.Story.Append(kb.HTML(el))

	page.Story.Append(kb.HTML("<p><b>Community:</b></p>"))

	page.Story.Append(simpleform.New(
		"/"+string(page.Slug),
		"",
		simpleform.Field("name", "Name"),
		simpleform.Option("rights", []string{string(kb.Reader), string(kb.Editor), string(kb.Moderator)}),
		simpleform.Button("ADD-COMMUNITY", "Add"),
		simpleform.Button("REMOVE-COMMUNITY", "Remove"),
	))

	el = `<ul class="tight">`
	for _, member := range members {
		if member.IsGroup {
			el += "<li>" + html.EscapeString(member.Name) + " = " + string(member.Access) + "</li>"
		}
	}
	el += "</ul>"
	page.Story.Append(kb.HTML(el))

	page.WriteResponse(w)
}
