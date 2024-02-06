package cmds

import (
	"flag"
	"fmt"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kb/pgdb"
)

func init() {
	Register(Command{
		Name: "destroy-database-permanently",
		Desc: "Destroys all database content",
		Run:  DestroyDatabase,
	})
}

func DestroyDatabase(DB kb.Database, fs *flag.FlagSet, args []string) {
	//nolint:errcheck
	fs.Parse(args)

	db := DB.(*pgdb.Database)

	_, err := db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO public;
		COMMENT ON SCHEMA public IS 'standard public schema';
	`)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Database has been destroyed.")
	}
}
