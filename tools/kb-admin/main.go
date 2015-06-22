package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kb/pgdb"
	"github.com/raintreeinc/knowledgebase/tools/kb-admin/cmds"

	_ "github.com/lib/pq"
)

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

func OpenDB() kb.Database {
	params := RDS()
	if params == "" {
		params = os.Getenv("DATABASE")
	}
	db, err := pgdb.New(params)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	args := os.Args[1:]
	if len(args) <= 0 {
		fmt.Println("Available commands:")
		for _, cmd := range cmds.List {
			fmt.Printf("%16s : %s\n", cmd.Name, cmd.Desc)
		}
		return
	}

	cmdname := args[0]
	for _, cmd := range cmds.List {
		if cmd.Name == cmdname {
			fs := flag.NewFlagSet(cmdname, flag.ExitOnError)
			fs.Usage = fs.PrintDefaults
			cmd.Run(OpenDB(), fs, args[1:])
			return
		}
	}

	fmt.Println("Unknown command " + cmdname)
}
