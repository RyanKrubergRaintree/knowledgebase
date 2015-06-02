package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/raintreeinc/livebundle"

	"github.com/raintreeinc/knowledgebase/auth"
	"github.com/raintreeinc/knowledgebase/kbserver"
	"github.com/raintreeinc/knowledgebase/kbserver/pgdb"

	"github.com/raintreeinc/knowledgebase/module/admin"
	"github.com/raintreeinc/knowledgebase/module/dita"
	"github.com/raintreeinc/knowledgebase/module/group"
	"github.com/raintreeinc/knowledgebase/module/page"
	"github.com/raintreeinc/knowledgebase/module/tag"
	"github.com/raintreeinc/knowledgebase/module/user"

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

	if os.Getenv("DEVELOPMENT") != "" {
		v, err := strconv.ParseBool(os.Getenv("DEVELOPMENT"))
		if err == nil {
			*development = v
		}
	}

	log.Printf("Development %v\n", *development)
	log.Printf("Starting with database %s\n", *database)
	log.Printf("Starting %s on %s\n", *domain, *addr)

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

	db.Groups().Create(kbserver.Group{
		ID:          "admin",
		Name:        "Admin",
		Public:      false,
		Description: "Administrators",
	})

	// context
	store := sessions.NewFilesystemStore("", []byte("some secret"))
	context := kbserver.NewContext(*domain, store)

	// create KnowledgeBase server
	server := kbserver.New(*domain, *templatesdir, db, client, context)

	// add systems
	server.AddModule(admin.New(server))
	server.AddModule(group.New(server))
	server.AddModule(page.New(server))
	server.AddModule(tag.New(server))
	server.AddModule(user.New(server))

	if *ditamap != "" {
		server.AddModule(dita.New("Dita", *ditamap, server))
	}

	// protect server with authentication
	url := "http://" + *domain
	auth.Register(os.Getenv("APPKEY"), url, auth.ClientsFromEnv())
	front := auth.New(server)

	http.Handle("/", front)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
