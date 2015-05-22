package testdata

import (
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ = &kb.Page{}

func SetupDatabase(db kbserver.Database) {
	db.Users().Add("Egon Elbre")
}
