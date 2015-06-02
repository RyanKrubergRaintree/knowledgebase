package main

import (
	"flag"
	"log"
	"os"

	"github.com/raintreeinc/knowledgebase/kbserver/pgdb"

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

	// Load database
	db, err := pgdb.New(*database)
	if err != nil {
		log.Fatal(err)
	}

	users, err := db.Users().List()
	log.Println(err)
	for _, user := range users {
		log.Printf("%+v\n", user)
	}
}
