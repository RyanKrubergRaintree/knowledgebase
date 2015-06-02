package main

import (
	"flag"
	"log"
	"os"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbdita"
	"github.com/raintreeinc/knowledgebase/kbserver"
	"github.com/raintreeinc/knowledgebase/kbserver/pgdb"

	"github.com/raintreeinc/knowledgebase/ditaconv"

	_ "github.com/lib/pq"
)

var (
	database = flag.String("database", "user=root dbname=knowledgebase sslmode=disable", "database `params`")
)

func main() {
	flag.Parse()
	if os.Getenv("DATABASE") != "" {
		*database = os.Getenv("DATABASE")
	}

	ditamap := os.Getenv("DITAMAP")

	// Load database
	db, err := pgdb.New(*database)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec(`DELETE * FROM Pages WHERE Owner = "help"`)

	//users, _ := db.Users().List()
	//log.Println(users)

	db.Users().Create(kbserver.User{
		ID:   "help-uploader",
		Name: "Help Uploader",
	})
	db.Groups().AddMember("help", "help-uploader")

	index, errs := ditaconv.LoadIndex(ditamap)
	log.Println(errs)

	mapping, errs := ditaconv.CreateMapping(index)
	log.Println(errs)

	owner := kb.Slug("help")
	for topic, slug := range mapping.ByTopic {
		ownerslug := owner + ":" + slug
		mapping.ByTopic[topic] = ownerslug
		delete(mapping.BySlug, slug)
		mapping.BySlug[ownerslug] = topic
	}

	group := db.PagesByGroup("help-uploader", "help")
	mapping.Rules.Merge(kbdita.RaintreeDITA())

	for _, topic := range mapping.BySlug {
		page, fatal, errs := mapping.Convert(topic)
		if fatal != nil {
			log.Println(fatal)
			continue
		} else if len(errs) > 0 {
			log.Println(errs)
		}

		if page.Slug == "" {
			log.Println("EMPTY SLUG")
			continue
		}

		page.Owner = owner
		if page.Slug[0] == '/' {
			page.Slug = page.Slug[1:]
		}

		err = group.Create(page)
		if err != nil {
			log.Println("ERROR", page.Slug, ":", err)
		} else {
			log.Println(page.Slug)
		}
	}
}
