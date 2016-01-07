package dita

import (
	"encoding/json"
	"encoding/xml"
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

	"github.com/bradfitz/slice"

	"github.com/raintreeinc/knowledgebase/extra/ditaindex"
	"github.com/raintreeinc/knowledgebase/kb"

	"github.com/raintreeinc/knowledgebase/ditaconv"
	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
)

var _ *kb.Page
var _ kb.Module = &Module{}

type Module struct {
	name    string
	ditamap string
	server  *kb.Server

	store atomic.Value
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
	mod.store.Store(newstore())
	go mod.monitor()
}

func (mod *Module) Pages() (r []kb.PageEntry) {
	store := mod.store.Load().(*store)
	for _, page := range store.pages {
		r = append(r, kb.PageEntryFrom(page))
	}
	return
}

func (mod *Module) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store := mod.store.Load().(*store)
	path := strings.TrimPrefix(r.URL.Path, "/")
	slug := kb.Slugify(path)
	if data, ok := store.raw[slug]; ok {
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
		for _, err := range store.errLoad {
			page.Story.Append(kb.Paragraph(err.Error()))
		}

		page.Story.Append(kb.HTML("<h3>Mapping</h3>"))
		for _, err := range store.errMapping {
			page.Story.Append(kb.Paragraph(err.Error()))
		}

		page.Story.Append(kb.HTML("<h3>Converting</h3>"))
		for _, errs := range store.errConvert {
			text := "<h4>[[" + string(errs.slug) + "]]</h4>"
			for _, err := range errs.errors {
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
		for _, slug := range store.slugs {
			page := store.pages[slug]
			content += fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", slug, html.EscapeString(page.Title))
		}
		content += "</ul>"

		page.Story.Append(kb.HTML(content))
		page.WriteResponse(w)
		return

	case name + "=index":
		page := &kb.Page{
			Slug:     name + "=index",
			Title:    "Index",
			Modified: time.Now(),
		}

		page.Story.Append(ditaindex.New("index", store.index))
		page.WriteResponse(w)
		return
	}
	http.NotFound(w, r)
}

func (mod *Module) reload() {
	start := time.Now()
	mod.store.Store(load(mod.name, mod.ditamap))
	log.Println("DITA reloaded (", time.Since(start), ")")
}

func (mod *Module) monitor() {
	modified := time.Now()
	mod.reload()
	for range time.Tick(1 * time.Second) {
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

func newstore() *store {
	return &store{
		pages: make(map[kb.Slug]*kb.Page),
		raw:   make(map[kb.Slug][]byte),
	}
}

type store struct {
	pages map[kb.Slug]*kb.Page
	raw   map[kb.Slug][]byte
	slugs []kb.Slug
	index *ditaindex.Item

	errLoad    []error
	errMapping []error
	errConvert []convertError
}

type convertError struct {
	slug   kb.Slug
	fatal  error
	errors []error
}

func load(prefix, ditamap string) *store {
	store := newstore()

	index, errs := ditaconv.LoadIndex(ditamap)
	store.errLoad = errs

	mapping, errs := ditaconv.CreateMapping(index)
	store.errMapping = errs

	for topic, slug := range mapping.ByTopic {
		ownerslug := kb.Slugify(prefix+"=") + slug
		mapping.ByTopic[topic] = ownerslug
		delete(mapping.BySlug, slug)
		mapping.BySlug[ownerslug] = topic
	}

	store.index = ditaindex.EntryToItem(mapping, index.Nav)

	mapping.Rules.Merge(RaintreeDITA())
	for slug, topic := range mapping.BySlug {
		page, fatal, errs := mapping.Convert(topic)
		if fatal != nil {
			store.errConvert = append(store.errConvert, convertError{slug: slug, fatal: fatal})
		} else if len(errs) > 0 {
			store.errConvert = append(store.errConvert, convertError{slug: slug, errors: errs})
		}

		data, err := json.Marshal(page)
		if err != nil {
			log.Println(err)
		}

		store.pages[slug] = page
		store.raw[slug] = data
		store.slugs = append(store.slugs, slug)
	}

	slice.Sort(store.slugs, func(i, j int) bool {
		return store.slugs[i] < store.slugs[j]
	})

	return store
}

func RaintreeDITA() *xmlconv.Rules {
	return &xmlconv.Rules{
		Translate: map[string]string{
			"keystroke": "span",
			"secright":  "span",

			// faq
			"faq":          "dl",
			"faq-question": "dt",
			"faq-answer":   "dd",

			//UI items
			"ui-item-list": "dl",

			"ui-item-name":        "dt",
			"ui-item-description": "dd",

			// setup options
			"setup-options": "dl",

			"setup-option-name":        "dt",
			"setup-option-description": "dd",

			"settingdesc": "div",
			"settingname": "h3",
		},
		Remove: map[string]bool{
			"settinghead": true,
		},
		Unwrap: map[string]bool{
			"ui-item":      true,
			"faq-item":     true,
			"setup-option": true,

			"settings": true,
			"setting":  true,
		},
		Callback: map[string]xmlconv.Callback{
			"settingdefault": func(enc xmlconv.Encoder, dec *xml.Decoder, start *xml.StartElement) error {
				val, _ := xmlconv.Text(dec, start)
				if val != "" {
					err := enc.WriteRaw("<p>Default value: " + val + "</p>")
					if err != nil {
						return err
					}
				}
				return nil
			},
			"settinglevels": func(enc xmlconv.Encoder, dec *xml.Decoder, start *xml.StartElement) error {
				err := enc.WriteRaw("<p>Levels where it can be defined:</p>")
				if err != nil {
					return err
				}
				if err := enc.Rules().ConvertChildren(enc, dec, start); err != nil {
					return err
				}
				return nil
			},
			"settingsample": func(enc xmlconv.Encoder, dec *xml.Decoder, start *xml.StartElement) error {
				err := enc.WriteRaw("<p>Example:</p>")
				if err != nil {
					return err
				}
				if err := enc.Rules().ConvertChildren(enc, dec, start); err != nil {
					return err
				}
				return nil
			},
		},
	}
}

/*

type Mapping struct {
	Index   *Index
	Topics  map[string]*Topic
	BySlug  map[kb.Slug]*Topic
	ByTopic map[*Topic]kb.Slug
	Rules   *xmlconv.Rules
}

func (m *Mapping) TopicsSorted() (r []*Topic) {
	for _, topic := range m.Topics {
		r = append(r, topic)
	}
	sort.Sort(byfilename(r))
	return r
}

type byfilename []*Topic

func (a byfilename) Len() int           { return len(a) }
func (a byfilename) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byfilename) Less(i, j int) bool { return a[i].Filename < a[j].Filename }

func CreateMapping(index *Index) (*Mapping, []error) {
	topics := index.Topics

	var errors []error
	byslug := make(map[kb.Slug]*Topic, len(topics))
	bytopic := make(map[*Topic]kb.Slug, len(topics))

	// assign slugs to the topics
	for _, topic := range topics {
		topic.Title = topic.Title
		topic.ShortTitle = topic.ShortTitle
		slug := kb.Slugify(topic.Title)

		if other, clash := byslug[slug]; clash {
			errors = append(errors, fmt.Errorf("clashing title \"%v\" in \"%v\" and \"%v\"", topic.Title, topic.Filename, other.Filename))
			continue
		}

		if topic.Title == "" {
			errors = append(errors, fmt.Errorf("title missing in \"%v\"", topic.Filename))
			continue
		}

		byslug[slug] = topic
		bytopic[topic] = slug
	}

	// promote to shorter titles, if possible
	for prev, topic := range byslug {
		if topic.ShortTitle == "" || len(topic.Title) <= len(topic.ShortTitle) {
			continue
		}

		slug := kb.Slugify(topic.ShortTitle)
		if _, exists := byslug[slug]; exists {
			continue
		}
		topic.Title = topic.ShortTitle
		topic.ShortTitle = ""

		delete(byslug, prev)
		byslug[slug] = topic
		bytopic[topic] = slug
	}

	m := &Mapping{
		Rules:   NewHTMLRules(),
		Index:   index,
		Topics:  topics,
		BySlug:  byslug,
		ByTopic: bytopic,
	}

	return m, errors
}


*/
