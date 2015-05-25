package kbdita

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/bradfitz/slice"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"

	"github.com/raintreeinc/knowledgebase/ditaconv"
	"github.com/raintreeinc/knowledgebase/ditaconv/xmlconv"
)

var _ *kb.Page
var _ kbserver.System = &System{}

type System struct {
	name    string
	ditamap string
	server  *kbserver.Server

	store atomic.Value
}

func New(name, ditamap string, server *kbserver.Server) *System {
	sys := &System{
		name:    name,
		ditamap: ditamap,
		server:  server,
	}
	sys.init()
	return sys
}

func (sys *System) Name() string { return sys.name }

func (sys *System) init() {
	sys.store.Store(newstore())
	go sys.monitor()
}

func (sys *System) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	store := sys.store.Load().(*store)
	path := strings.TrimPrefix(r.URL.Path, "/")
	slug := kb.Slugify(path)
	if data, ok := store.raw[slug]; ok {
		w.Write(data)
		return
	}

	name := kb.Slugify(sys.name)
	switch slug {
	case name + ":conversion-errors":
	case name + ":all-pages":
		page := &kb.Page{
			Slug:     name + ":all-pages",
			Title:    "All Pages",
			Modified: time.Now(),
		}

		content := "<ul>"
		for _, slug := range store.slugs {
			page := store.pages[slug]
			content += fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", slug, page.Title)
		}
		content += "</ul>"

		page.Story.Append(kb.HTML(content))
		kbserver.WriteJSON(w, r, page)
		return
	}
	http.NotFound(w, r)
}

func (sys *System) reload() {
	start := time.Now()
	sys.store.Store(load(sys.name, sys.ditamap))
	log.Println("DITA reloaded (", time.Since(start), ")")
}

func (sys *System) monitor() {
	modified := time.Now()
	sys.reload()
	for range time.Tick(1 * time.Second) {
		filepath.Walk(filepath.Dir(sys.ditamap),
			func(_ string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if info.ModTime().After(modified) {
					modified = time.Now()
					sys.reload()
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
		ownerslug := kb.Slugify(prefix+":") + slug
		mapping.ByTopic[topic] = ownerslug
		delete(mapping.BySlug, slug)
		mapping.BySlug[ownerslug] = topic
	}

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
