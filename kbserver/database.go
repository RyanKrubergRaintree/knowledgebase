package kbserver

import (
	"errors"

	"github.com/raintreeinc/knowledgebase/kb"
)

var (
	ErrUserNotExist   = errors.New("user does not exist")
	ErrGroupNotExist  = errors.New("group does not exist")
	ErrUserNotAllowed = errors.New("user does not have sufficient permissions")
	ErrPageExists     = errors.New("page already exists")
	ErrPageNotExist   = errors.New("page does not exist")
)

type Database interface {
	Users() Users
	Groups() Groups

	CanRead(user, group kb.Slug) bool
	CanWrite(user, group kb.Slug) bool

	PagesByGroup(user, group kb.Slug) Pages
	IndexByUser(user kb.Slug) Index
}

type Users interface {
	ByID(id kb.Slug) (User, error)
	Create(user User) error
	Delete(id kb.Slug) error
	List() ([]User, error)
}

type Groups interface {
	ByID(id kb.Slug) (Group, error)
	Create(group Group) error
	Delete(group kb.Slug) error
	List() ([]Group, error)

	AddMember(group, user kb.Slug) error
	RemoveMember(group, user kb.Slug) error
	MembersOf(group kb.Slug) ([]User, error)
}

type Pages interface {
	Create(page *kb.Page) error
	Load(slug kb.Slug) (*kb.Page, error)
	LoadRaw(slug kb.Slug) ([]byte, error)
	Save(slug kb.Slug, page *kb.Page) error
	// SaveRaw(slug kb.Slug, page []byte) error
	List() ([]kb.PageEntry, error)
}

type Index interface {
	List() ([]kb.PageEntry, error)
	Search(text string) ([]kb.PageEntry, error)

	Tags() ([]kb.TagEntry, error)
	ByTag(tag string) ([]kb.PageEntry, error)

	Groups() ([]Group, error)
	ByGroup(group kb.Slug) ([]kb.PageEntry, error)

	RecentChanges(n int) ([]kb.PageEntry, error)
}

type User struct {
	ID          kb.Slug
	Name        string
	Email       string
	Description string
	Admin       bool
	Groups      []string
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
	ID          kb.Slug
	Name        string
	Public      bool
	Description string
}