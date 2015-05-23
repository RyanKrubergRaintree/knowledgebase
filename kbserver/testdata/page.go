package testdata

import (
	"github.com/raintreeinc/knowledgebase/kb"
)

func NewPage(owner string, title string) *kb.Page {
	page := &kb.Page{
		Owner:    kb.Slugify(owner),
		Slug:     kb.Slugify(owner + ":" + title),
		Title:    title,
		Synopsis: "Welcome to a test page.",
		Story: kb.Story{
			kb.Tags("Welcome", "Home", "Some Example"),
			kb.Paragraph("Link to self: [[Community:Welcome]]."),
			kb.Paragraph("External link: [[http://google.com google.com]]"),
			kb.Paragraph("kbpage: [[page:Pages]], [[page:Recent Changes]]"),
			kb.Paragraph("kbgroup: [[group:Groups]], [[group:Community]], [[group:Community Details]]"),
			kb.Paragraph("kbtag: [[tag:Tags]], [[tag:Home]]"),
			kb.Paragraph("kbuser: [[user:Egon Elbre]]"),
		},
	}

	return page
}
