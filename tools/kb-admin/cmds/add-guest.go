package cmds

import (
	"flag"
	"fmt"
	"os"

	"github.com/raintreeinc/knowledgebase/kb"
)

func init() {
	Register(Command{
		Name: "add-guest",
		Desc: "Add guest user",
		Run:  AddGuest,
	})
}

func AddGuest(DB kb.Database, fs *flag.FlagSet, args []string) {
	name := fs.String("name", "", "name for the login")
	email := fs.String("email", "", "email for the login")
	password := fs.String("password", "", "password for the login")
	fs.Parse(args)

	if *name == "" || *password == "" {
		fmt.Println("username and password must be specified")
		fs.Usage()
		os.Exit(1)
	}

	guest := DB.Context("admin").GuestLogin()
	err := guest.Add(*name, *email, *password)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
