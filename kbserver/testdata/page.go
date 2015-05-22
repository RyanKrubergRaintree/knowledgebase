package testdata

import (
	"math/rand"

	lorem "github.com/drhodes/golorem"
	"github.com/raintreeinc/knowledgebase/kb"
)

func NewPage(owner string, title string) *kb.Page {
	page := &kb.Page{
		Owner:    kb.Slugify(owner),
		Slug:     kb.Slugify(title),
		Title:    title,
		Synopsis: lorem.Paragraph(1, 1),
		Story: kb.Story{
			kb.Paragraph("Simple link: [[Community:Simple]]."),
			kb.Paragraph("Link to self: [[Community:Welcome]]."),
			kb.Paragraph("External link: [[http://neti.ee neti.ee]]"),
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
