package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/raintreeinc/knowledgebase/assets"
	"github.com/raintreeinc/knowledgebase/auth"
	"github.com/raintreeinc/knowledgebase/kb"

	"github.com/gorilla/sessions"
)

// TODO: add
//  https://github.com/unrolled/secure
//  https://github.com/justinas/nosurf

var (
	addr      = flag.String("listen", ":80", "http server `address`")
	assetsdir = flag.String("assets", "assets", "assets `directory`")
	database  = flag.String("database", "", "database `params`")
	domain    = flag.String("domain", "", "`domain`")
	conffile  = flag.String("config", "knowledgebase.toml", "farm configuration")
)

func main() {
	flag.Parse()

	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	if host != "" || port != "" {
		*addr = host + ":" + port
	}

	if os.Getenv("ASSETSDIR") != "" {
		*assetsdir = os.Getenv("ASSETSDIR")
	}
	if os.Getenv("DATABASE") != "" {
		*database = os.Getenv("DATABASE")
	}
	if os.Getenv("DOMAIN") != "" {
		*domain = os.Getenv("DOMAIN")
	}

	log.Printf("Starting with database %s\n", *database)
	log.Printf("Starting with domain %s\n", *domain)

	log.Printf("Starting %s on %s", *domain, *addr)

	// Serve static files
	files := assets.NewFiles(*assetsdir, []string{".css", ".png", ".ico", ".jpg", ".js"})
	http.Handle("/static/", files)

	presenter := assets.NewPresenter(*assetsdir, "*.html", map[string]string{
		"ShortTitle": "KB",
		"Title":      "Knowledge Base",
		"Company":    "Raintree Systems Inc.",
	})

	// context
	store := sessions.NewFilesystemStore("", []byte("some secret"))
	context := kb.NewContext(store)

	// create KnowledgeBase server
	server := kb.NewServer(*domain, *database, presenter)

	// protect server with authentication
	url := "http://" + *domain
	auth.Register(os.Getenv("APPKEY"), url, auth.ClientsFromEnv())
	front := auth.New(server, context, presenter)

	http.Handle("/", front)

	log.Fatal(http.ListenAndServe(*addr, nil))

	/*
		//TODO: move domain initialization inside farm
		auth := &auth.Context{
			Renderer:    renderer,
			Domain:      conf.Domain,
			LoginURL:    "http://auth." + conf.Domain + "/login",
			CallbackURL: "http://auth." + conf.Domain + "/callback",
			Sessions:    sessions.NewFilesystemStore("", []byte("some secret")),
		}
		auth.RegisterProviders()

		admin := &admin.Server{Renderer: renderer, Database: conf.Database}

		farm, err := farm.New(conf, auth, admin)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Starting %s on %s", conf.Domain, *addr)
		log.Fatal(http.ListenAndServe(*addr, farm))
	*/
}
