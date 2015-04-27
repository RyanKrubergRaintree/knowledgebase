package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/raintreeinc/knowledgebase/farm"
	"github.com/raintreeinc/knowledgebase/farm/admin"
	"github.com/raintreeinc/knowledgebase/farm/auth"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
)

// TODO: add
//  https://github.com/unrolled/secure
//  https://github.com/justinas/nosurf

var (
	defaultDir = filepath.Join("assets", "templates", "**")

	addr      = flag.String("listen", ":80", "http server `address`")
	clientdir = flag.String("client", "client", "client `directory`")
	templates = flag.String("templates", defaultDir, "templates `glob`")
	conffile  = flag.String("config", "knowledgebase.toml", "farm configuration")
)

func main() {
	flag.Parse()

	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	if host != "" || port != "" {
		*addr = host + ":" + port
	}

	conf := farm.Config{}
	if _, err := toml.DecodeFile(*conffile, &conf); err != nil {
		log.Fatal(err)
	}

	if conf.ClientDir == "" {
		conf.ClientDir = *clientdir
	}

	if os.Getenv("CLIENTDIR") != "" {
		conf.ClientDir = os.Getenv("CLIENTDIR")
	}
	if os.Getenv("DATABASE") != "" {
		conf.Database = os.Getenv("DATABASE")
	}
	if os.Getenv("DOMAIN") != "" {
		conf.Domain = os.Getenv("DOMAIN")
	}

	log.Printf("Starting with database %s\n", conf.Database)
	log.Printf("Starting with domain %s\n", conf.Domain)

	renderer, err := farm.NewRenderer(*templates)
	if err != nil {
		log.Fatal(err)
	}

	//TODO: move domain initialization inside farm
	auth := &auth.Context{
		Renderer:    renderer,
		Domain:      conf.Domain,
		LoginURL:    "http://auth." + conf.Domain + "/login",
		CallbackURL: "http://auth." + conf.Domain + "/callback",
		Sessions:    sessions.NewFilesystemStore("", []byte("some secret")),
	}
	auth.RegisterProviders()

	admin := &admin.Server{Database: conf.Database}

	farm, err := farm.New(conf, auth, admin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting %s on %s", conf.Domain, *addr)
	log.Fatal(http.ListenAndServe(*addr, farm))
}
