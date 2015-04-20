package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/raintreeinc/knowledgebase/admin"
	"github.com/raintreeinc/knowledgebase/auth"
	"github.com/raintreeinc/knowledgebase/kb"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
)

// TODO: add
//  https://github.com/unrolled/secure
//  https://github.com/justinas/nosurf

var (
	addr      = flag.String("listen", ":80", "http server `address`")
	clientdir = flag.String("client", "", "client `directory`")
	conffile  = flag.String("config", "knowledgebase.toml", "farm configuration")
)

func main() {
	flag.Parse()

	host, port := os.Getenv("HOST"), os.Getenv("PORT")
	if host != "" || port != "" {
		*addr = host + ":" + port
	}

	if os.Getenv("KB_CLIENTDIR") != "" {
		*clientdir = os.Getenv("KB_CLIENTDIR")
	}

	conf := kb.FarmConfig{}
	if _, err := toml.DecodeFile(*conffile, &conf); err != nil {
		log.Fatal(err)
	}

	if conf.ClientDir == "" {
		conf.ClientDir = *clientdir
	}

	//TODO: move domain initialization inside farm
	auth := &auth.Context{
		Domain:      conf.Domain,
		LoginURL:    "http://auth." + conf.Domain + "/login",
		CallbackURL: "http://auth." + conf.Domain + "/callback",
		Sessions:    sessions.NewFilesystemStore("", []byte("some secret")),
	}
	auth.RegisterProviders()

	admin := &admin.Server{
		Database: conf.Database,
	}

	farm, err := kb.NewFarm(conf, auth, admin)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting %s on %s", conf.Domain, *addr)
	log.Fatal(http.ListenAndServe(*addr, farm))
}
