package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/raintreeinc/knowledgebase/auth"
	"github.com/raintreeinc/knowledgebase/kbserver"
	"github.com/raintreeinc/knowledgebase/kbserver/memdb"

	"github.com/gorilla/sessions"
)

// TODO: add
//  https://github.com/unrolled/secure
//  https://github.com/justinas/nosurf

var (
	addr     = flag.String("listen", ":80", "http server `address`")
	database = flag.String("database", "", "database `params`")
	domain   = flag.String("domain", "", "`domain`")
	conffile = flag.String("config", "knowledgebase.toml", "farm configuration")

	templatesdir = flag.String("templates", "templates", "templates `directory`")
	assetsdir    = flag.String("assets", "assets", "assets `directory`")
	//todo replace clientx with client
	clientdir = flag.String("client", "client", "client `directory`")
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
	if os.Getenv("CLIENTDIR") != "" {
		*clientdir = os.Getenv("CLIENTDIR")
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
	assets := kbserver.NewFiles(*assetsdir)
	http.Handle("/assets/", http.StripPrefix("/assets/", assets))

	// Serve javascript source
	source := kbserver.NewSource(*clientdir, true)
	http.Handle("/lib/", http.StripPrefix("/lib/", source))

	// Load database
	//TODO: replace with real database
	db := memdb.New(*database)

	// context
	store := sessions.NewFilesystemStore("", []byte("some secret"))
	context := kbserver.NewContext(*domain, store)

	// presenter
	presenter := kbserver.NewPresenter(*templatesdir, "*.html", map[string]string{
		"ShortTitle": "KB",
		"Title":      "Knowledge Base",
		"Company":    "Raintree Systems Inc.",
	}, source, context)

	// create KnowledgeBase server
	server := kbserver.New(*domain, db, presenter, context)

	// protect server with authentication
	url := "http://" + *domain
	auth.Register(os.Getenv("APPKEY"), url, auth.ClientsFromEnv())
	front := auth.New(server, context, presenter)

	// allow cross origin requests on sub-domains
	cors := kbserver.AllowSubdomainCORS(*domain, front)

	http.Handle("/", cors)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
