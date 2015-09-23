package admin

import (
	"bytes"
	"html"
	"html/template"
	"net/http"
	"strings"

	"github.com/raintreeinc/knowledgebase/extra/simpleform"
	"github.com/raintreeinc/knowledgebase/kb"
)

var esc = html.EscapeString

func (mod *Module) groups(w http.ResponseWriter, r *http.Request) {
	context, ok := mod.server.AdminContext(w, r)
	if !ok {
		return
	}

	if r.Method == "POST" {
		action := r.Header.Get("action")
		switch action {
		case "create-group":
			name := strings.TrimSpace(r.FormValue("name"))
			if name == "" {
				http.Error(w, "Name not specified.", http.StatusBadRequest)
				return
			}
			owner := strings.TrimSpace(r.FormValue("owner"))
			description := strings.TrimSpace(r.FormValue("description"))
			public := strings.TrimSpace(r.FormValue("visibility")) == "public"

			if owner == "" {
				owner = name
			}

			var err error
			group := kb.Group{
				ID:          kb.Slugify(name),
				Name:        name,
				OwnerID:     kb.Slugify(owner),
				Public:      public,
				Description: description,
			}

			err = context.Groups().Create(group)
			if err != nil {
				kb.WriteResult(w, err)
				return
			}
			w.Write([]byte("group created"))
			return
		case "add-user":
			user := strings.TrimSpace(r.FormValue("user"))
			group := strings.TrimSpace(r.FormValue("group"))
			if user == "" || group == "" {
				http.Error(w, "User/Group not specified.", http.StatusBadRequest)
				return
			}

			userid := kb.Slugify(user)
			groupid := kb.Slugify(group)
			err := context.Access().AddUser(groupid, userid)

			if err != nil {
				kb.WriteResult(w, err)
				return
			}
			w.Write([]byte("user added"))
			return
		default:
			http.Error(w, "Invalid action "+action+" specified", http.StatusBadRequest)
			return
		}
	}

	page := &kb.Page{
		Slug:  "admin=groups",
		Title: "Groups",
	}

	page.Story.Append(simpleform.New(
		"/"+string(page.Slug), "",
		simpleform.Field("name", "Name"),
		simpleform.Field("owner", "Owner"),
		simpleform.Field("description", "Description"),
		simpleform.Option("visibility", []string{"public", "private"}),
		simpleform.Button("create-group", "Create"),
	))

	groups, err := context.Groups().List()
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	var buf bytes.Buffer
	err = templGroups.Execute(&buf, groups)
	if err != nil {
		kb.WriteResult(w, err)
		return
	}

	page.Story.Append(kb.HTML(buf.String()))

	page.Story.Append(kb.HTML("<h2>Add user</h2>"))
	page.Story.Append(simpleform.New(
		"/"+string(page.Slug), "",
		simpleform.Field("group", "Group"),
		simpleform.Field("user", "User"),
		simpleform.Button("add-user", "Add"),
	))

	page.WriteResponse(w)
}

var templGroups = template.Must(template.New("").Parse(`
	<table class="tight">
		<thead><td>Name</td><td>Owner</td><td>Public</td><td>Description</td><td></td></thead>
		{{ range . }}
		<tr>
			<td>{{.Name}}</td>
			<td>{{if (ne .ID .OwnerID)}}{{.OwnerID}}{{end}}</td>
			<td>{{if .Public}}{{ else }}private{{end}}</td>
			<td>{{.Description}}</td>
			<td><a class="mdi mdi-pencil" style="text-decoration:none;" href="/group=moderate-{{.ID}}"></a></td>
		</tr>
		{{ end }}
	</table>
`))
