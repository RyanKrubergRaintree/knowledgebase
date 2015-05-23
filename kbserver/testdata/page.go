package testdata

import (
	"math/rand"

	lorem "github.com/drhodes/golorem"
	"github.com/raintreeinc/knowledgebase/kb"
)

func NewPage(owner string, title string) *kb.Page {
	page := &kb.Page{
		Owner:    kb.Slugify(owner),
		Slug:     kb.Slugify(owner + ":" + title),
		Title:    title,
		Synopsis: lorem.Paragraph(1, 1),
		Story: kb.Story{
			kb.Tags("Welcome", "Home", "Some Example"),
			kb.Paragraph("Link to self: [[Community:Welcome]]."),
			kb.Paragraph("External link: [[http://google.com google.com]]"),
			kb.Paragraph("kbpage: [[page:Pages]], [[page:Recent Changes]]"),
			kb.Paragraph("kbgroup: [[group:Groups]], [[group:Community]], [[group:Community Details]]"),
			kb.Paragraph("kbtag: [[tag:Tags]], [[tag:Home]]"),
			kb.Paragraph("kbuser: [[user:Egon Elbre]]"),
			kb.Paragraph(lorem.Paragraph(1, 1)),
			kb.Paragraph(lorem.Paragraph(1, 1)),
			kb.Paragraph(lorem.Paragraph(1, 1)),
		},
	}

	N := rand.Intn(5)
	for i := 0; i < N; i++ {
		page.Story.Append(kb.Paragraph(lorem.Paragraph(1, 1)))
	}

	return page
}
