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
	PagesByGroup(user, group string) Pages
	IndexByUser(user string) Index
}

type Users interface {
	Create(user User) error
	Delete(name string) error
	ByName(name string) (User, error)
	List() ([]User, error)
}

type Groups interface {
	Create(group Group) error
	Delete(name string) error
	List() ([]Group, error)

	AddMember(group string, user string) error
	RemoveMember(group string, user string) error
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

	RecentChanges(n int) ([]kb.PageEntry, error)
}

type User struct {
	Name        string
	Email       string
	Description string
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
	Name        string
	Public      bool
	Description string
}
