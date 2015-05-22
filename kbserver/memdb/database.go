// memdb implements a in-memory database
package memdb

import (
	"encoding/json"

	"github.com/raintreeinc/knowledgebase/kb"
	"github.com/raintreeinc/knowledgebase/kbserver"
)

type User struct {
	Name   string
	Info   kb.User
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

type Group struct {
	Public bool
	Pages  map[kb.Slug]*kb.Page
}

type Database struct {
	Users  map[string]*User
	Groups map[string]*Group
}

func New(params string) *Database {
	return &Database{
		Users:  make(map[string]*User),
		Groups: make(map[string]*Group),
	}
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

	group, exists := db.Groups[groupname]
	if !exists {
		return nil, kbserver.ErrGroupNotExist
	}

	return Pages{group}, nil
}

func (db *Database) IndexByUser(username string) (kbserver.Index, error) {
	user, err := db.User(username)
	if err != nil {
		return nil, err
	}
	return &Index{db, user}, nil
}

type Pages struct {
	*Group
}

func (group Pages) All() ([]kb.PageEntry, error) {
	entries := make([]kb.PageEntry, 0, len(group.Pages))
	for slug, page := range group.Pages {
		entries = append(entries, kb.PageEntryFrom(page))
	}
	return entries, nil
}

func (group Pages) Exists(slug kb.Slug) bool {
	_, exists := group.Pages[slug]
	return exists
}

func (group Pages) Create(slug kb.Slug, page *kb.Page) error {
	if group.Exists(slug) {
		return kbserver.ErrPageExists
	}

	group.Pages[slug] = page
	return nil
}

func (group Pages) Load(slug kb.Slug) (*kb.Page, error) {
	if group.Exists(slug) {
		return kbserver.ErrPageExists
	}

	page, exists := group.Pages[slug]
	if !exists {
		return nil, kbserver.ErrPageNotExist
	}
	return page, nil
}

func (group Pages) LoadRaw(slug kb.Slug) ([]byte, error) {
	page, err := group.Load(slug)
	if err != nil {
		return nil, err
	}

	return json.Marshal(page)
}

func (group Pages) Save(slug kb.Slug, page *kb.Page) error {
	if !group.Exists(slug) {
		return kbserver.ErrPageNotExist
	}

	group.Pages[slug] = page
	return nil
}

type Index struct {
	*Database
	User *User
}

func (index *Index) All() ([]kb.PageEntry, error) {
	entries := []kb.PageEntry{}

	for gname, group := range index.Groups {
		if !index.User.BelongsTo(gname) && !group.Public {
			continue
		}

		all, _ := (Pages{group}).All()
		entries = append(entries, all...)
	}

	kb.SortPageEntriesBySlug(entries)
	return entries, nil
}

func (index *Index) Search(text string) ([]kb.PageEntry, error) {
	//TODO:
	return []kb.PageEntry{}, nil
}

func (index *Index) Tags() ([]kb.TagEntry, error) {
	tags := make(map[string]int)
	pages, _ := index.All()

	for _, page := range pages {
		for _, tag := range page.Tags {
			tags[tag]++
		}
	}

	entries := make([]kb.TagEntry, 0, len(tags))
	for name, count := range tags {
		entries = append(entries, kb.TagEntry{
			Name:  name,
			Count: count,
		})
	}

	kb.SortTagEntriesByName(entries)
	return entries
}

func (index *Index) ByTag(tag string) ([]kb.PageEntry, error) {
	all, err := index.All()
	if err != nil {
		return nil, err
	}

	entries := make([]kb.PageEntry, 0, len(all))
	for _, entry := range all {
		if entry.HasTag(tag) {
			entries = append(entries, entry)
		}
	}

	kb.SortPageEntriesBySlug(entries)
	return entries, nil
}

func (index *Index) RecentChanges(n int) ([]kb.PageEntry, error) {
	entries, err := index.All()
	if err != nil {
		return nil, err
	}

	kb.SortPageEntriesByDate(entries)

	if n < len(entries) {
		entries = entries[:n]
	}
	return entries, nil
}
