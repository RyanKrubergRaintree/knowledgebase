package tag

import (
	"html"
	"net/http"
	"strconv"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"
)

var _ kb.Module = &Module{}

type Module struct {
	server *kb.Server
}

func New(server *kb.Server) *Module {
	return &Module{server}
}

func (mod *Module) Info() kb.Group {
	return kb.Group{
		ID:          "tag",
		Name:        "Tag",
		Public:      true,
		Description: "Displays tag index.",
	}
}

func (mod *Module) Pages() []kb.PageEntry {
	return []kb.PageEntry{{
		Slug:     "tag=tags",
		Title:    "Tags",
		Synopsis: "Listing of all tags.",
	}}
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/tag=tags" {
		mod.tags(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/tag=pages/") {
		id := strings.TrimPrefix(r.URL.Path, "/tag=pages/")
		if id == "" {
			http.Error(w, "tag-id missing", http.StatusBadRequest)
			return
		}
		mod.pages(w, r, kb.Slugify(id))
	} else if strings.HasPrefix(r.URL.Path, "/tag=first/") {
		id := strings.TrimPrefix(r.URL.Path, "/tag=first/")
		if id == "" {
			http.Error(w, "id missing", http.StatusBadRequest)
			return
		}
		mod.first(w, r, kb.Slugify(id))
	} else {
		http.NotFound(w, r)
	}
}

func (mod *Module) pages(w http.ResponseWriter, r *http.Request, tag kb.Slug) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	entries, err := index.ByTag(tag)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &kb.Page{
		Slug:  "tag=" + tag,
		Title: kb.SlugToTitle(tag),
		Story: kb.StoryFromEntries(entries),
	}

	page.WriteResponse(w)
}

func (mod *Module) tags(w http.ResponseWriter, r *http.Request) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	entries, err := index.Tags()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	page := &kb.Page{
		Slug:  "tag=tags",
		Title: "Tags",
	}

	if len(entries) == 0 {
		page.Story.Append(kb.Paragraph("No results."))
	} else {
		for _, entry := range entries {
			page.Story.Append(kb.Entry(
				html.EscapeString(entry.Name),
				strconv.Itoa(entry.Count)+" pages",
				kb.Slugify("tag="+entry.Name)))
		}
	}

	page.WriteResponse(w)
}

func (mod *Module) first(w http.ResponseWriter, r *http.Request, tag kb.Slug) {
	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	filter := r.FormValue("filter")

	var entries []kb.PageEntry
	var err error
	if filter == "" {
		entries, err = index.ByTag(tag)
	} else {
		filter = string(kb.Slugify(filter))
		entries, err = index.ByTagFilter(tag, "help-", "help-"+filter)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(entries) == 0 {
		http.Error(w, "No entries.", http.StatusResetContent)
		return
	}

	first := entries[0]
	http.Redirect(w, r, string("/"+first.Slug), http.StatusSeeOther)
}
