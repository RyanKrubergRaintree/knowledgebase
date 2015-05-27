package testdata

import (
	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

var _ = &kb.Page{}

func SetupDatabase(db kbserver.Database) {
	db.Groups().Create(kbserver.Group{
		ID:          "community",
		Name:        "Community",
		Public:      true,
		Description: "All editing users",
	})
	db.Groups().Create(kbserver.Group{
		ID:          "engineering",
		Name:        "Engineering",
		Public:      true,
		Description: "Raintree Engineering",
	})
	db.Groups().Create(kbserver.Group{
		ID:          "help",
		Name:        "Help",
		Public:      true,
		Description: "Raintree Help",
	})

	db.Users().Create(kbserver.User{
		ID:          "admin",
		Name:        "Admin",
		Email:       "",
		Description: "",
	})

	db.Groups().AddMember("engineering", "admin")
	db.Groups().AddMember("community", "admin")

	db.PagesByGroup("admin", "community").Create(NewPage("Community", "Welcome"))

	db.Users().Create(kbserver.User{
		ID:          "egon-elbre",
		Name:        "Egon Elbre",
		Email:       "egonelbre@gmail.com",
		Admin:       true,
		Description: "Raintree Engineering",
	})
	db.Groups().AddMember("engineering", "egon-elbre")
	db.Groups().AddMember("community", "egon-elbre")
}
