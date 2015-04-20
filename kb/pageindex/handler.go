package pageindex

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/egonelbre/fedwiki"
	"github.com/egonelbre/fedwiki/item"
)

type Handler struct {
	Index
}

func (index Handler) entries(title string, pages []*fedwiki.PageHeader, err error) (code int, template string, data interface{}) {
	if err != nil {
		return fedwiki.ErrorResponse(http.StatusInternalServerError, err.Error())
	}

	return http.StatusOK, "", &fedwiki.Page{
		PageHeader: fedwiki.PageHeader{
			Title: title,
			Date:  fedwiki.NewDate(time.Now()),
		},
		Story: HeadersToStory(pages),
	}
}

func (index Handler) Handle(r *http.Request) (code int, template string, data interface{}) {
	slug := fedwiki.Slug(r.URL.Path)
	if err := fedwiki.ValidateSlug(slug); err != nil {
		return fedwiki.ErrorResponse(http.StatusBadRequest, err.Error())
	}

	tokens := strings.SplitN(string(slug), "/", 3)
	if len(tokens) < 2 {
		return fedwiki.ErrorResponse(http.StatusBadRequest, "Invalid path.")
	}

	switch tokens[1] {
	case "all":
		pages, err := index.All()
		//TODO: add information when last changed
		return index.entries("Pages", pages, err)
	case "recent-changes":
		n := 30
		if len(tokens) >= 3 {
			x, err := strconv.Atoi(tokens[2])
			if err != nil {
				n = x
			}
		}
		pages, err := index.RecentChanges(n)
		//TODO: add information about when changed
		return index.entries("Recent Changes", pages, err)
	case "search":
		query := r.URL.Query().Get("q")
		pages, err := index.Search(query)
		title := "Search \"" + query + "\""
		return index.entries(title, pages, err)
	case "tag":
		if len(tokens) < 3 {
			return fedwiki.ErrorResponse(http.StatusBadRequest, "Invalid request.")
		}
		pages, err := index.PagesByTag(tokens[2])
		title := "Pages tagged \"" + tokens[2] + "\""
		return index.entries(title, pages, err)
	case "tags":
		tags, err := index.Tags()
		if err != nil {
			return fedwiki.ErrorResponse(http.StatusInternalServerError, err.Error())
		}

		p := &fedwiki.Page{}
		p.Title = "Tags"
		p.Date = fedwiki.NewDate(time.Now())

		if len(tags) == 0 {
			p.Story.Append(item.Paragraph("No results found."))
		}
		for _, ti := range tags {
			p.Story.Append(Entry(
				ti.Name,
				"",
				//TODO: proper reverse map instead of hardcoding the path
				fedwiki.Slug("/index/tag/"+ti.Name),
			))
		}

		return http.StatusOK, "", p
	case "citations":
		if len(tokens) < 3 {
			return fedwiki.ErrorResponse(http.StatusBadRequest, "Invalid request.")
		}
		slug := fedwiki.Slug("/" + tokens[2])
		pages, err := index.PagesCiting(slug)
		title := "Pages citing \"" + string(slug) + "\""
		return index.entries(title, pages, err)
	}

	return fedwiki.ErrorResponse(http.StatusNotFound, `Index "%s" does not exist.`, tokens[1])
}
