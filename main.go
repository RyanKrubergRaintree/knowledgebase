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
	files := assets.NewFiles(*assetsdir)
	http.Handle("/static/", files)

	// context
	store := sessions.NewFilesystemStore("", []byte("some secret"))
	context := kb.NewContext(store)

	// presenter
	presenter := assets.NewPresenter(*assetsdir, "*.html", map[string]string{
		"ShortTitle": "KB",
		"Title":      "Knowledge Base",
		"Company":    "Raintree Systems Inc.",
	}, context)

	// create KnowledgeBase server
	server := kb.NewServer(*domain, *database, presenter)

	// protect server with authentication
	url := "http://" + *domain
	auth.Register(os.Getenv("APPKEY"), url, auth.ClientsFromEnv())
	front := auth.New(server, context, presenter)

	http.Handle("/", front)

	log.Fatal(http.ListenAndServe(*addr, nil))
}
