package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"sync/atomic"
	"time"

	"github.com/egonelbre/fedwiki"
	"github.com/egonelbre/fedwiki/item"
	"github.com/egonelbre/fedwiki/template"

	"github.com/raintreeinc/knowledgebase/dita/ditaconv"
	"github.com/raintreeinc/knowledgebase/dita/ditaindex"
	"github.com/raintreeinc/knowledgebase/dita/xmlconv"
)

const defaultAddr = ":8001"

var (
	addr      = flag.String("listen", os.Getenv("KBDITA"), "listening `address`")
	ditamap   = flag.String("ditamap", os.Getenv("KBDITAMAP"), "main `ditamap`")
	clientdir = flag.String("client", os.Getenv("KBCLIENT"), "kbclient `directory`")
	pagesdir  = flag.String("pagesdir", "_pages", "output `directory`")

	dirviews = flag.String("views", "views", "`directory` for page templates")

	global        atomic.Value
	triggerReload = make(chan struct{}, 1e6)
)

type convertErr struct {
	Slug   fedwiki.Slug
	Fatal  error
	Errors []error
}

type rendered struct {
	Page *fedwiki.Page
	JSON []byte
}

type Store struct {
	Pages   map[fedwiki.Slug]*rendered
	Headers []*fedwiki.PageHeader
	Slugs   []fedwiki.Slug

	Created time.Time

	LoadErrs    []error
	MappingErrs []error
	ConvertErrs []convertErr
}

func (s *Store) ErrorPage() *fedwiki.Page {
	page := &fedwiki.Page{}
	page.Slug = "/system/errors"
	page.Title = "Errors"
	page.Date = fedwiki.NewDate(s.Created)

	page.Story.Append(item.HTML(`<form action='/system/reload' target="_blank" method='POST'><input type='submit' value='Reload'></form>`))

	page.Story.Append(item.HTML("<h3>Loading</h3>"))
	for _, err := range s.LoadErrs {
		page.Story.Append(item.Paragraph(err.Error()))
	}

	page.Story.Append(item.HTML("<h3>Mapping</h3>"))
	for _, err := range s.MappingErrs {
		page.Story.Append(item.Paragraph(err.Error()))
	}

	page.Story.Append(item.HTML("<h3>Converting</h3>"))
	for _, errs := range s.ConvertErrs {
		text := "<h4>[" + string(errs.Slug) + "]</h4>"
		for _, err := range errs.Errors {
			text += "<p>" + err.Error() + "</p>"
		}
		page.Story.Append(item.HTML(text))
	}

	return page
}

func (s *Store) HomePage() *fedwiki.Page {
	page := &fedwiki.Page{}
	page.Slug = "/home"
	page.Title = "Home"
	page.Date = fedwiki.NewDate(s.Created)

	content := "<h3>Pages:</h3>"
	content += "<ul>"
	for _, h := range s.Headers {
		content += fmt.Sprintf("<li><a href=\"%s\">%s</a></li>", h.Slug, h.Title)
	}
	content += "</ul>"
	page.Story.Append(item.HTML(content))

	return page
}

func reload() (*Store, error) {
	if *pagesdir != "" {
		os.RemoveAll(*pagesdir)
		os.MkdirAll(*pagesdir, 0755)
	}

	store := &Store{
		Pages:   make(map[fedwiki.Slug]*rendered),
		Created: time.Now(),
	}

	index, errs := ditaindex.Load(*ditamap)
	store.LoadErrs = errs

	mapping, errs := ditaconv.CreateMapping(index)
	store.MappingErrs = errs

	mapping.Rules.Merge(CustomRules())

	for slug, topic := range mapping.BySlug {
		page, fatal, errs := mapping.Convert(topic)
		if fatal != nil {
			store.ConvertErrs = append(store.ConvertErrs, convertErr{Slug: slug, Fatal: fatal})
			continue
		} else if len(errs) > 0 {
			store.ConvertErrs = append(store.ConvertErrs, convertErr{Slug: slug, Errors: errs})
		}

		rendered := &rendered{}
		rendered.Page = page
		rendered.JSON, _ = json.Marshal(page)

		if *pagesdir != "" {
			filename := string(page.Slug) + ".json"
			err := ioutil.WriteFile(filepath.Join(*pagesdir, filename), rendered.JSON, 0755)
			if err != nil {
				log.Println(err)
			}
		}

		store.Pages[slug] = rendered
		store.Headers = append(store.Headers, &page.PageHeader)
		store.Slugs = append(store.Slugs, slug)
	}

	sort.Sort(ByTitle(store.Headers))
	sort.Sort(BySlug(store.Slugs))

	return store, nil
}

type ByTitle []*fedwiki.PageHeader

func (s ByTitle) Len() int           { return len(s) }
func (s ByTitle) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByTitle) Less(i, j int) bool { return s[i].Title < s[j].Title }

type BySlug []fedwiki.Slug

func (s BySlug) Len() int           { return len(s) }
func (s BySlug) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySlug) Less(i, j int) bool { return s[i] < s[j] }

func reloader() {
	for {
		log.Println("Reloading")
		store, err := reload()
		log.Println("Reloaded")

		if err == nil {
			global.Store(store)
		} else {
			log.Println(err)
		}

		// clear reload channel
	clearReload:
		for {
			select {
			case <-triggerReload:
			default:
				break clearReload
			}
		}

		// wait for reload
		select {
		case <-time.After(1 * time.Hour):
		case <-triggerReload:
		}
	}
}

func main() {
	flag.Parse()
	if *ditamap == "" {
		fmt.Printf("specifying a ditamap is required")
		flag.Usage()
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if *addr == "" && port != "" {
		*addr = "localhost:" + port
	} else if *addr == "" {
		*addr = defaultAddr
	}
	if *clientdir == "" {
		*clientdir = filepath.Join("..", "kbclient")
	}

	go reloader()

	templates := template.New(filepath.Join(*dirviews, "*"))

	server := &fedwiki.Server{
		Handler:  fedwiki.HandlerFunc(serve),
		Template: templates,
	}

	http.Handle("/client/", http.StripPrefix("/client/", http.FileServer(http.Dir(*clientdir))))
	http.Handle("/",
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Access-Control-Allow-Methods", "GET")
			rw.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
			if r.URL.Path == "" || r.URL.Path == "/" {
				http.ServeFile(rw, r, filepath.Join(*clientdir, "index.html"))
				return
			}

			if r.Method == "POST" && r.URL.Path == "/system/reload" {
				triggerReload <- struct{}{}
				http.Error(rw, "", http.StatusOK)
				return
			}

			server.ServeHTTP(rw, r)
		}))

	log.Println("listening on", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func serve(r *http.Request) (code int, template string, data interface{}) {
	if r.Method != "GET" {
		return fedwiki.ErrorResponse(http.StatusForbidden, "Method %s is not allowed", r.Method)
	}

	store, ok := global.Load().(*Store)
	if !ok {
		return fedwiki.ErrorResponse(http.StatusNotFound, "Page %s not loaded", r.URL.Path)
	}

	switch r.URL.Path {
	case "/system/sitemap":
		return http.StatusOK, "sitemap", store.Headers
	case "/system/slugs":
		return http.StatusOK, "slugs", store.Slugs
	case "/system/errors":
		return http.StatusOK, "", store.ErrorPage()
	case "/home":
		return http.StatusOK, "", store.HomePage()
	}

	slug := fedwiki.Slugify(r.URL.Path[1:])
	page, found := store.Pages[slug]
	if found {
		return http.StatusOK, "", page.Page
	} else {
		return fedwiki.ErrorResponse(http.StatusNotFound, "Page %s not found", r.URL.Path)
	}
}

func CustomRules() *xmlconv.Rules {
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
		},
		Unwrap: map[string]bool{
			"ui-item":      true,
			"faq-item":     true,
			"setup-option": true,
		},
	}
}
