package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/raintreeinc/livebundle"

	"github.com/raintreeinc/knowledgebase/auth"
	"github.com/raintreeinc/knowledgebase/kbserver"
	"github.com/raintreeinc/knowledgebase/kbserver/pgdb"

	"github.com/raintreeinc/knowledgebase/kbadmin"
	"github.com/raintreeinc/knowledgebase/kbdita"
	"github.com/raintreeinc/knowledgebase/kbgroup"
	"github.com/raintreeinc/knowledgebase/kbpage"
	"github.com/raintreeinc/knowledgebase/kbtag"
	"github.com/raintreeinc/knowledgebase/kbuser"

	"github.com/gorilla/sessions"

	_ "github.com/lib/pq"
)

// TODO: add
//  https://github.com/unrolled/secure
//  https://github.com/justinas/nosurf

var (
	addr     = flag.String("listen", ":80", "http server `address`")
	database = flag.String("database", "user=root dbname=knowledgebase sslmode=disable", "database `params`")
	domain   = flag.String("domain", "", "`domain`")
	conffile = flag.String("config", "knowledgebase.toml", "farm configuration")

	development = flag.Bool("development", true, "development mode")
	ditamap     = flag.String("dita", "", "ditamap file for showing live dita")

	templatesdir = flag.String("templates", "templates", "templates `directory`")
	assetsdir    = flag.String("assets", "assets", "assets `directory`")
	clientdir    = flag.String("client", "client", "client `directory`")
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
	if os.Getenv("DITAMAP") != "" {
		*ditamap = os.Getenv("DITAMAP")
	}

	log.Printf("Starting with database %s\n", *database)
	log.Printf("Starting with domain %s\n", *domain)

	log.Printf("Starting %s on %s", *domain, *addr)

	// Serve static files
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir(*assetsdir))))

	// Serve client
	client := livebundle.New(*clientdir, "/client/", *development)
	http.Handle("/client/", client)

	// Load database
	db, err := pgdb.New(*database)
	if err != nil {
		log.Fatal(err)
	}

	// create default groups
	db.Groups().Create(kbserver.Group{
		ID:          "community",
		Name:        "Community",
		Public:      true,
		Description: "All editing users",
	})

	db.Groups().Create(kbserver.Group{
		ID:          "engineering",
		Name:        "Engineering",
		Public:      true,
		Description: "Raintree Engineering",
	})

	db.Groups().Create(kbserver.Group{
		ID:          "help",
		Name:        "Help",
		Public:      true,
		Description: "Raintree Help",
	})

	// context
	store := sessions.NewFilesystemStore("", []byte("some secret"))
	context := kbserver.NewContext(*domain, store)

	// create KnowledgeBase server
	server := kbserver.New(*domain, *templatesdir, db, client, context)

	// add systems
	server.AddSystem(kbadmin.New(server))
	server.AddSystem(kbgroup.New(server))
	server.AddSystem(kbpage.New(server))
	server.AddSystem(kbtag.New(server))
	server.AddSystem(kbuser.New(server))

	if *ditamap != "" {
		server.AddSystem(kbdita.New("Dita", *ditamap, server))
	}

	// protect server with authentication
	url := "http://" + *domain
	auth.Register(os.Getenv("APPKEY"), url, auth.ClientsFromEnv())
	front := auth.New(server)

	http.Handle("/", front)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
