package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/raintreeinc/livebundle"

	"github.com/raintreeinc/knowledgebase/auth"
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kb/pgdb"

	"github.com/raintreeinc/knowledgebase/module/admin"
	"github.com/raintreeinc/knowledgebase/module/dita"
	"github.com/raintreeinc/knowledgebase/module/group"
	"github.com/raintreeinc/knowledgebase/module/page"
	"github.com/raintreeinc/knowledgebase/module/tag"
	"github.com/raintreeinc/knowledgebase/module/user"

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

	rules = flag.String("rules", "rules.json", "different rules for server")

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
	if rds := RDS(); rds != "" {
		*database = rds
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
	if os.Getenv("RULES") != "" {
		*rules = os.Getenv("RULES")
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

	log.Println("Initializing DB")
	if err := db.Initialize(); err != nil {
		log.Fatal(err)
	}
	log.Println("DB Initialization complete.")

	// protect server with authentication
	url := "http://" + *domain
	auth.Register(os.Getenv("APPKEY"), url, "/system/auth", auth.ClientsFromEnv())

	// create server
	server := kb.NewServer(kb.ServerInfo{
		Domain:     *domain,
		ShortTitle: "KB",
		Title:      "Knowledge Base",
		Company:    "Raintree Systems Inc.",
	}, auth.New(), client, db)

	ruleset := MustLoadRules(*rules)
	server.Rules = ruleset

	// add systems
	server.AddModule(admin.New(server))
	server.AddModule(group.New(server))
	server.AddModule(page.New(server))
	server.AddModule(tag.New(server))
	server.AddModule(user.New(server))

	if *ditamap != "" {
		server.AddModule(dita.New("Dita", *ditamap, server))
	}

	http.Handle("/", server)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type RuleSet struct {
	Admins []kb.Slug
	Email  map[string][]kb.Slug `json:"email"`
}

func (rs *RuleSet) Login(user kb.User, db kb.Database) error {
	context := db.Context("admin")
	_, err := context.Users().ByID(user.ID)
	created := err == nil

	createUserIfNeeded := func() {
		if !created {
			err := context.Users().Create(user)
			if err != nil {
				log.Println(err)
				return
			}
			created = true
		}
	}

	for _, admin := range rs.Admins {
		if user.ID == admin {
			createUserIfNeeded()
			context.Access().SetAdmin(user.ID, true)
		}
	}

	for rule, groups := range rs.Email {
		matched, err := regexp.MatchString(rule, user.Email)
		if err != nil {
			log.Println(err)
		}
		if matched && err == nil {
			createUserIfNeeded()
			for _, group := range groups {
				context.Access().AddUser(group, user.ID)
			}
		}
	}
	return nil
}

func MustLoadRules(filename string) *RuleSet {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	rs := &RuleSet{}
	if err := json.Unmarshal(data, rs); err != nil {
		log.Fatal(err)
	}

	return rs
}

func RDS() string {
	user := os.Getenv("RDS_USERNAME")
	pass := os.Getenv("RDS_PASSWORD")

	dbname := os.Getenv("RDS_DB_NAME")
	host := os.Getenv("RDS_HOSTNAME")
	port := os.Getenv("RDS_PORT")

	if user == "" || pass == "" || dbname == "" || host == "" || port == "" {
		return ""
	}

	return fmt.Sprintf("user='%s' password='%s' dbname='%s' host='%s' port='%s'", user, pass, dbname, host, port)
}
