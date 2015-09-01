package dispatch

import (
	"net/http"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"
)

var _ kb.Module = &Module{}

type Module struct {
	group  kb.Group
	server *kb.Server
}

func New(group kb.Group, server *kb.Server) *Module {
	mod := &Module{
		group:  group,
		server: server,
	}
	return mod
}

func (mod *Module) Info() kb.Group {
	return mod.group
}

func (mod *Module) Pages() []kb.PageEntry { return []kb.PageEntry{} }

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	groupID, titleID, pageID := kb.TokenizeLink3(r.URL.Path)

	if groupID != mod.group.ID {
		http.Error(w, "Invalid owner specified:\nexpected "+string(mod.group.ID)+".",
			http.StatusBadRequest)
		return
	}

	_, index, ok := mod.server.IndexContext(w, r)
	if !ok {
		return
	}

	temp, err := index.ByTitle(titleID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	entries := []kb.PageEntry{}
	prefix := string(mod.group.ID + "/")
	for _, entry := range temp {
		if !strings.HasPrefix(string(entry.Slug), prefix) {
			continue
		}
		groupID, _ := kb.TokenizeLink(string(entry.Slug))
		title := strings.TrimPrefix(string(groupID), prefix)

		entries = append(entries, kb.PageEntry{
			Title: title,
			Slug:  entry.Slug,
		})
	}

	kb.SortPageEntriesByTitle(entries)
	kb.ReversePageEntries(entries)

	page := &kb.Page{Slug: pageID}
	if len(entries) > 0 {
		page.Title = entries[0].Title
		page.Story.Append(kb.Tags(entries[0].Tags...))
		page.Story.Append(kb.Paragraph(entries[0].Synopsis))
	} else {
		page.Title = kb.SlugToTitle(titleID)
	}

	if len(entries) == 0 {
		page.Story.Append(kb.Paragraph("No pages."))
	} else {
		page.Story.Append(kb.HTML("<h2>Versions</h2>"))
		for _, entry := range entries {
			page.Story.Append(kb.Entry(entry.Title, "", entry.Slug))
		}

	}

	page.WriteResponse(w)
}
