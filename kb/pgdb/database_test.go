package pgdb_test

import (
	"testing"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kb/pgdb"
)

const dbparams = `
	user=integration
	password=integration
	dbname=knowledgebase_integration
	sslmode=disable
`

func TestIntegration(t *testing.T) {
	db, err := pgdb.New(dbparams)
	if err != nil {
		t.Fatal(err)
		return
	}

	log := func(txt string, err error) {
		if err != nil {
			t.Errorf(txt + ":" + err.Error())
		}
	}

	assert := func(txt string, ok bool) {
		if !ok {
			t.Errorf(txt)
		}
	}

	_, err = db.Exec(`
		DROP SCHEMA public CASCADE;
		CREATE SCHEMA public;
		GRANT ALL ON SCHEMA public TO integration;
		GRANT ALL ON SCHEMA public TO public;
		COMMENT ON SCHEMA public IS 'standard public schema';
	`)
	log("Destroying database", err)
	log("Initializing database", db.Initialize())

	context := db.Context("admin")

	// Managing groups:
	log("Creating a public group", context.Groups().Create(kb.Group{
		ID:      "public",
		OwnerID: "public",
		Name:    "Public",
		Public:  true,
	}))

	assert("Duplicate group", context.Groups().Create(kb.Group{ID: "public"}) == kb.ErrGroupExists)

	log("Creating a private group", context.Groups().Create(kb.Group{
		ID:      "private",
		OwnerID: "private",
		Name:    "Private",
		Public:  false,
	}))

	log("Creating a member group", context.Groups().Create(kb.Group{
		ID:      "member",
		OwnerID: "member",
		Name:    "Member",
		Public:  false,
	}))

	log("Creating a delete group", context.Groups().Create(kb.Group{
		ID:      "delete",
		OwnerID: "delete",
		Name:    "Delete",
		Public:  true,
	}))

	log("Deleting a group", context.Groups().Delete("delete"))

	group, err := context.Groups().ByID("private")
	log("Getting private group", err)
	assert("Getting private group", group.ID == "private")

	_, err = context.Groups().ByID("delete")
	if err != kb.ErrGroupNotExist {
		log("Getting deleted group", err)
	}

	groups, err := context.Groups().List()
	log("Listing groups", err)
	assert("Listing 3 groups", len(groups) == 3)

	// Managing users:
	log("Creating a reader user", context.Users().Create(kb.User{
		ID:        "reader",
		Name:      "Reader",
		MaxAccess: kb.Moderator,
	}))

	log("Creating an editor user", context.Users().Create(kb.User{
		ID:        "editor",
		Name:      "Editor",
		Email:     "editor@example.com",
		MaxAccess: kb.Moderator,
	}))
	assert("Duplicate user", context.Users().Create(kb.User{ID: "editor", MaxAccess: kb.Moderator}) == kb.ErrUserExists)

	log("Creating an moderator user", context.Users().Create(kb.User{
		ID:        "moderator",
		Name:      "Moderator",
		MaxAccess: kb.Moderator,
	}))

	log("Creating an admin user", context.Users().Create(kb.User{
		ID:        "admin",
		Name:      "Admin",
		Admin:     true,
		MaxAccess: kb.Moderator,
	}))

	log("Creating a dynamic user", context.Users().Create(kb.User{
		ID:        "dynamic",
		Name:      "Dynamic",
		MaxAccess: kb.Moderator,
	}))

	log("Creating a delete user", context.Users().Create(kb.User{
		ID:        "delete",
		Name:      "Delete",
		MaxAccess: kb.Moderator,
	}))

	log("Deleting an user", context.Users().Delete("delete"))

	user, err := context.Users().ByID("editor")
	log("Getting editor user", err)
	assert("Getting editor user", user.ID == "editor" && user.Email == "editor@example.com")

	_, err = context.Users().ByID("delete")
	if err != kb.ErrUserNotExist {
		log("Getting deleted user", err)
	}

	users, err := context.Users().List()
	log("Listing users", err)
	assert("Listing 5 users", len(users) == 5)

	// Managing accesses:
	rights := func(txt string, exp, got kb.Rights) {
		if exp != got {
			t.Errorf(txt+": exp %v got %v", exp, got)
		}
	}

	assert("reader should not be an admin", !context.Access().IsAdmin("reader"))
	assert("editor should not be an admin", !context.Access().IsAdmin("editor"))
	assert("admin should be an admin", context.Access().IsAdmin("admin"))

	rights("non-member accessing public page", kb.Reader, context.Access().Rights("public", "reader"))
	rights("non-member accessing private page", kb.Blocked, context.Access().Rights("private", "reader"))

	log("add dynamic rights", context.Access().AddUser("private", "dynamic"))
	rights("dynamic rights", kb.Moderator, context.Access().Rights("private", "dynamic"))
	log("remove dynamic rights", context.Access().RemoveUser("private", "dynamic"))
	rights("dynamic rights", kb.Blocked, context.Access().Rights("private", "dynamic"))

	log("Creating a readers group", context.Groups().Create(kb.Group{ID: "readers", OwnerID: "readers", Name: "Readers"}))
	log("Creating a moderators group", context.Groups().Create(kb.Group{ID: "moderators", OwnerID: "moderators", Name: "Moderator"}))

	log("add reader to group", context.Access().AddUser("readers", "reader"))
	log("add moderator to group", context.Access().AddUser("moderators", "moderator"))

	log("Adding private <- readers", context.Access().CommunityAdd("private", "readers", kb.Reader))
	log("Adding private <- moderators", context.Access().CommunityAdd("private", "moderators", kb.Moderator))

	rights("Reader rights", kb.Reader, context.Access().Rights("private", "reader"))
	rights("Moderator rights", kb.Moderator, context.Access().Rights("private", "moderator"))

	log("Updating private <- moderators", context.Access().CommunityAdd("private", "moderators", kb.Editor))
	rights("Updating moderator rights", kb.Editor, context.Access().Rights("private", "moderator"))
	log("Removing private <- moderators", context.Access().CommunityRemove("private", "moderators"))
	rights("Removed rights", kb.Blocked, context.Access().Rights("private", "moderator"))

	members, err := context.Access().List("public")
	log("Access items for public", err)
	assert("Access items for public", len(members) == 0)

	log("add reader to private", context.Access().AddUser("private", "reader"))
	log("add moderator to private", context.Access().AddUser("private", "moderator"))

	members, err = context.Access().List("private")
	log("Access items for moderators", err)
	assert("Access items for moderators", len(members) == 3) // reader user, moderator user, readers group

	// handling of pages
	log("Creating page", context.Pages("private").Create(welcomePage))
	assert("Duplicate page creation", context.Pages("private").Create(welcomePage) == kb.ErrPageExists)

	page, err := context.Pages("private").Load("private:welcome")
	log("Loading page", err)
	assert("Correct page", samePage(page, welcomePage))

	log("Overwrite page", context.Pages("private").Overwrite("private:welcome", 1, welcomePage2))
	assert("Concurrent edit", context.Pages("private").Overwrite("private:welcome", 1, welcomePage2) == kb.ErrConcurrentEdit)

	log("Add paragraph", context.Pages("private").Edit("private:welcome", 4, kb.Action{
		"type": "add",
		"item": kb.Paragraph("Hello World..."),
	}))

	pages, err := context.Pages("private").List()
	log("List pages", err)
	assert("Must have 1 entry", len(pages) == 1)

	assert("Concurrent delete page", context.Pages("private").Delete("private:welcome", 1) == kb.ErrConcurrentEdit)
	log("Delete page", context.Pages("private").Delete("private:welcome", 5))
	_, err = context.Pages("private").Load("private:welcome")
	assert("Loading deleted page", err == kb.ErrPageNotExist)

	// index handling
	pages, err = context.Index("reader").List()
	log("List index", err)
	assert("No pages ATM", len(pages) == 0)

	log("Create page", context.Pages("private").Create(welcomePage))
	pages, err = context.Index("reader").List()
	log("List index", err)
	assert("List single page", len(pages) == 1)

	pages, err = context.Index("reader").Search("lorem")
	log("Search index", err)
	assert("Search single page", len(pages) == 1)

	tags, err := context.Index("reader").Tags()
	log("List tags", err)
	assert("Two tags", len(tags) == 2 && tags[0].Name == "lorem" && tags[1].Name == "welcome")

	pages, err = context.Index("reader").ByTag("welcome")
	log("Search by tag", err)
	assert("Tag single page", len(pages) == 1)

}

var welcomePage = &kb.Page{
	Slug:     "private:welcome",
	Title:    "Welcome",
	Version:  1,
	Synopsis: "Lorem ipsum dolor sit amet, consectetur adipisicing elit.",
	Story: kb.Story{
		kb.Tags("welcome", "lorem"),
		kb.Paragraph("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Deserunt officiis illum similique eos, quia cum perspiciatis corporis sint illo, magni placeat recusandae sit veritatis veniam, reprehenderit quibusdam, quidem minus error."),
	},
}

var welcomePage2 = &kb.Page{
	Slug:     "private:welcome",
	Title:    "Welcome",
	Version:  4,
	Synopsis: "Lorem ipsum dolor sit amet, consectetur adipisicing elit.",
	Story: kb.Story{
		kb.Tags("welcome", "lorem"),
		kb.Paragraph("Lorem ipsum dolor sit amet, consectetur adipisicing elit. Deserunt officiis illum similique eos, quia cum perspiciatis corporis sint illo, magni placeat recusandae sit veritatis veniam, reprehenderit quibusdam, quidem minus error."),
	},
}

func samePage(a, b *kb.Page) bool {
	return a.Slug == b.Slug &&
		a.Title == b.Title &&
		a.Version == b.Version &&
		a.Synopsis == b.Synopsis
}
