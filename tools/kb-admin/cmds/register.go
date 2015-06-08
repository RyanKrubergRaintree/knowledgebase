package cmds

import (
	"flag"

	"github.com/raintreeinc/knowledgebase/kb"
)

var List []Command

type Command struct {
	Name string
	Desc string
	Run  func(DB kb.Database, fs *flag.FlagSet, args []string)
}

func Register(cmd Command) { List = append(List, cmd) }
