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

func OpenDB() kb.Database {
	db, err := pgdb.New(os.Getenv("DATABASE"), os.Getenv("DOMAIN"))
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
			cmd.Run(OpenDB(), fs, args[1:])
			return
		}
	}

	fmt.Println("Unknown command " + cmdname)
}
