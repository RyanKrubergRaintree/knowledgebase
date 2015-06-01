package main

import (
	"flag"
	"log"
	"os"

	"github.com/raintreeinc/knowledgebase/kbserver/pgdb"

	_ "github.com/lib/pq"
)

var (
	database = flag.String("database", "user=root password='root' dbname=knowledgebase sslmode=disable", "database `params`")
)

func main() {
	flag.Parse()
	if os.Getenv("DATABASE") != "" {
		*database = os.Getenv("DATABASE")
	}

	// Load database
	db, err := pgdb.New(*database)
	if err != nil {
		log.Fatal(err)
	}

	index := db.IndexByUser("egon-elbre")
	pages, err := index.ByTag("tutorials")
	log.Println(err)
	for _, page := range pages {
		log.Println(page)
	}

	//	rows, err := db.Query(`SELECT Owner, Slug, Tags, NormTags FROM Pages`)
	//	for rows.Next() {
	//		var owner, slug kb.Slug
	//		var tags, ntags stringSlice
	//		rows.Scan(&owner, &slug, &tags, &ntags)
	//		log.Println(owner, "\t", slug, tags, ntags)
	//	}
}
