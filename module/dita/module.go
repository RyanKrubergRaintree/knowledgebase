package dita

import (
	"errors"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kb/items/index"
)

var _ *kb.Page
var _ kb.Module = &Module{}

type Module struct {
	name    string
	ditamap string
	server  *kb.Server

	cache atomic.Value
}

func New(name, ditamap string, server *kb.Server) *Module {
	mod := &Module{
		name:    name,
		ditamap: ditamap,
		server:  server,
	}
	mod.init()
	return mod
}

func (mod *Module) Info() kb.Group {
	return kb.Group{
		ID:          kb.Slugify(mod.name),
		Name:        mod.name,
		Public:      true,
		Description: "Displays \"" + mod.ditamap + "\" content.",
	}
}

func (mod *Module) init() {
	mod.cache.Store(NewConversion("", ""))
	go mod.monitor()
}

func (mod *Module) Pages() (r []kb.PageEntry) {
	cache := mod.cache.Load().(*Conversion)
	for _, page := range cache.Pages {
		r = append(r, kb.PageEntryFrom(page))
	}
	return
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cache := mod.cache.Load().(*Conversion)
	path := strings.TrimPrefix(r.URL.Path, "/")
	slug := kb.Slugify(path)
	if data, ok := cache.Raw[slug]; ok {
		//nolint:errcheck
		w.Write(data)
		return
	}

	name := kb.Slugify(mod.name)
	switch slug {
	case name + "=errors":
		page := &kb.Page{}
		page.Slug = name + "=errors"
		page.Title = "Errors"
		page.Modified = time.Now()

		page.Story.Append(kb.HTML("<h3>Loading</h3>"))
		for _, err := range cache.LoadErrors {
			page.Story.Append(kb.Paragraph(err.Error()))
		}

		page.Story.Append(kb.HTML("<h3>Mapping</h3>"))
		for _, err := range cache.MappingErrors {
			page.Story.Append(kb.Paragraph(err.Error()))
		}

		page.Story.Append(kb.HTML("<h3>Converting</h3>"))
		for _, errs := range cache.Errors {
			text := "<h4>[[" + string(errs.Slug) + "]]</h4>"
			if errs.Fatal != nil {
				text += "<p style='background:#f88;'>" + errs.Fatal.Error() + "</p>"
			}
			for _, err := range errs.Errors {
				text += "<p>" + err.Error() + "</p>"
			}
			page.Story.Append(kb.HTML(text))
		}

		w.Header().Set("Content-Type", "application/json")
		if err := page.Write(w); err != nil {
			log.Println(err)
		}
		return

	case name + "=all-pages":
		page := &kb.Page{
			Slug:     name + "=all-pages",
			Title:    "All Pages",
			Modified: time.Now(),
		}

		content := "<ul>"
		for _, slug := range cache.Slugs {
			page := cache.Pages[slug]
			content += fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", slug, html.EscapeString(page.Title))
		}
		content += "</ul>"

		page.Story.Append(kb.HTML(content))
		//nolint:errcheck
		page.WriteResponse(w)
		return

	case name + "=index":
		page := &kb.Page{
			Slug:     name + "=index",
			Title:    "Index",
			Modified: time.Now(),
		}

		page.Story.Append(index.New("index", cache.Nav))
		//nolint:errcheck
		page.WriteResponse(w)
		return
	}
	http.NotFound(w, r)
}

func (mod *Module) reload() {
	start := time.Now()

	context := NewConversion(kb.Slugify(mod.name), mod.ditamap)
	context.Run()
	mod.cache.Store(context)

	log.Println("DITA reloaded (", time.Since(start), ")")
}

func (mod *Module) monitor() {
	modified := time.Now()
	mod.reload()
	for range time.Tick(3 * time.Second) {
		//nolint:errcheck
		filepath.Walk(filepath.Dir(mod.ditamap),
			func(_ string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.ModTime().After(modified) {
					modified = time.Now()
					mod.reload()
					return errors.New("stop iterate")
				}
				return nil
			})
	}
}
