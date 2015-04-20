package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/oauth"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
)

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

	ctx := &oauth.Context{
		Domain:      conf.Domain,
		LoginURL:    "http://login." + conf.Domain + "/login",
		CallbackURL: "http://login." + conf.Domain + "/callback",
		Sessions:    sessions.NewFilesystemStore("", []byte("some secret")),
	}
	ctx.RegisterProviders()

	farm, err := kb.NewFarm(conf, ctx)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting %s on %s", conf.Domain, *addr)
	log.Fatal(http.ListenAndServe(*addr, farm))
}
