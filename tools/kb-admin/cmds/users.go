package cmds

import (
	"flag"
	"fmt"
	"os"

	"github.com/raintreeinc/knowledgebase/kb"
)

func init() {
	Register(Command{
		Name: "list-users",
		Desc: "Lists users",
		Run:  ListUsers,
	})
}

func ListUsers(DB kb.Database, fs *flag.FlagSet, args []string) {
	//nolint:errcheck
	fs.Parse(args)

	users, err := DB.Context("admin").Users().List()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, user := range users {
		fmt.Printf("%+v\n", user)
	}
}
