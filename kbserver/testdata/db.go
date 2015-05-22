package testdata

import (
	"log"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ = &kb.Page{}

func check(errors ...error) {
	failed := false
	for i, err := range errors {
		if err != nil {
			log.Println(i, err)
			failed = true
		}
	}
	if failed {
		log.Fatal("DB Setup failed!")
	}
}

func SetupDatabase(db kbserver.Database) {
	check(
		db.Groups().Create(kbserver.Group{
			Name:        "Community",
			Public:      true,
			Description: "All editing users",
		}),
		db.Groups().Create(kbserver.Group{
			Name:        "Engineering",
			Public:      true,
			Description: "Raintree Engineering",
		}),
		db.Groups().Create(kbserver.Group{
			Name:        "Help",
			Public:      true,
			Description: "Raintree Help",
		}),

		db.Users().Create(kbserver.User{
			Name:        "Admin",
			Email:       "",
			Description: "",
		}),

		db.Groups().AddMember("Engineering", "Admin"),
		db.Groups().AddMember("Community", "Admin"),

		db.PagesByGroup("Admin", "Community").Create(NewPage("Community", "Welcome")),

		db.Users().Create(kbserver.User{
			Name:        "Egon Elbre",
			Email:       "egonelbre@gmail.com",
			Description: "Raintree Engineering",
		}),
		db.Groups().AddMember("Engineering", "Egon Elbre"),
		db.Groups().AddMember("Community", "Egon Elbre"),
	)
}
