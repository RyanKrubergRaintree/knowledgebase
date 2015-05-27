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
			kb.Paragraph("Page: [[Page:Pages]], [[Page:Recent Changes]]"),
			kb.Paragraph("Group: [[Group:Groups]], [[Group:Community]], [[Group:Community Details]]"),
			kb.Paragraph("Tag: [[Tag:Tags]], [[Tag:Home]]"),
			kb.Paragraph("User: [[User:Egon Elbre]]"),
		},
	}

	return page
}
