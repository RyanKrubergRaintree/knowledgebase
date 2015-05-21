// stubdb implements a Database with a single page "/home" for each group.
package stubdb

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

type User struct {
	Name   string
	Groups []string
}

func (user *User) BelongsTo(groupname string) bool {
	for _, x := range user.Groups {
		if x == groupname {
			return true
		}
	}
	return false
}

type Database struct {
	Users map[string]*User
}

// params is a string "username:group1,group2;username2:group2,group3"
func New(params string) *Database {
	db := &Database{
		Users: make(map[string]*User),
	}

	for _, usergroups := range strings.Split(params, ";") {
		tokens := strings.Split(usergroups, ":")
		user := &User{tokens[0], strings.Split(tokens[1], ",")}
		db.Users[user.Name] = user
	}

	return db
}

func (db *Database) User(username string) (*User, error) {
	user, ok := db.Users[username]
	if !ok {
		return nil, kbserver.ErrInvalidUser
	}
	return user, nil
}

func (db *Database) Access(username, groupname string) (*User, error) {
	user, ok := db.Users[username]
	if !ok {
		return nil, kbserver.ErrInvalidUser
	}

	if !user.BelongsTo(groupname) {
		return nil, kbserver.ErrUserNotAllowed
	}
	return user, nil
}

func (db *Database) PagesByOwner(username, groupname string) (kbserver.Pages, error) {
	_, err := db.Access(username, groupname)
	if err != nil {
		return nil, err
	}
	return &Pages{groupname}, nil
}

func (db *Database) IndexByUser(username string) (kbserver.Index, error) {
	user, err := db.User(username)
	if err != nil {
		return nil, err
	}
	return &Index{user}, nil
}

type Pages struct {
	Owner string
}

func (pages *Pages) All() ([]kb.PageEntry, error) {
	return []kb.PageEntry{{
		Owner:    pages.Owner,
		Slug:     "/home",
		Title:    "Home",
		Synopsis: "Simple home page.",
		Tags:     []string{"home", "lorem"},
		Modified: time.Now(),
	}}, nil
}

func (pages *Pages) Exists(slug kb.Slug) bool {
	return slug == "/home"
}

func (pages *Pages) Create(slug kb.Slug, page *kb.Page) error {
	return kbserver.ErrUserNotAllowed
}

func (pages *Pages) Load(slug kb.Slug) (*kb.Page, error) {
	if slug != "/home" {
		return nil, kbserver.ErrPageMissing
	}

	return &kb.Page{
		Owner:    pages.Owner,
		Slug:     "/home",
		Title:    "Home",
		Synopsis: "Simple home page.",
		Story: kb.Story{
			kb.Tags("home", "lorem"),
			kb.Paragraph(loremipsum),
			kb.Paragraph(loremipsum),
			kb.Paragraph(loremipsum),
		},
	}, nil
}

func (pages *Pages) LoadRaw(slug kb.Slug) ([]byte, error) {
	page, err := pages.Load(slug)
	if err != nil {
		return nil, err
	}

	return json.Marshal(page)
}

func (pages *Pages) Save(slug kb.Slug, page *kb.Page) error {
	return kbserver.ErrUserNotAllowed
}

type Index struct {
	User *User
}

func (index *Index) All() ([]kb.PageEntry, error) {
	r := []kb.PageEntry{}
	for _, group := range index.User.Groups {
		pages, _ := (&Pages{group}).All()
		r = append(r, pages...)
	}
	return r, nil
}

func (index *Index) Search(text string) ([]kb.PageEntry, error) {
	if strings.Contains(loremipsum, text) {
		return index.All()
	}
	return []kb.PageEntry{}, nil
}

func (index *Index) Tags() ([]kb.TagEntry, error) {
	return []kb.TagEntry{
		{"home", rand.Intn(10) + 1},
		{"lorem", rand.Intn(10) + 1},
	}, nil
}

func (index *Index) ByTag(tag string) ([]kb.PageEntry, error) {
	if tag == "lorem" || tag == "home" {
		return index.All()
	}
	return nil, nil
}

func (index *Index) RecentChanges(n int) ([]kb.PageEntry, error) {
	return index.All()
}

const loremipsum = `Lorem ipsum [[dolor]] sit amet, consectetur adipisicing elit. 
Cum, ex, accusantium. Maiores magnam nostrum, illum [[inventore]], esse odio eveniet
ipsum architecto impedit fugit sit [[eaque]], aut! Fuga dolorum sunt nisi.`
