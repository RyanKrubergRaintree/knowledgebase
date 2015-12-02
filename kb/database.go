package kb

import (
	"encoding/gob"
	"errors"
	"html/template"
)

var (
	ErrUserExists    = errors.New("User already exists.")
	ErrUserNotExist  = errors.New("User does not exist.")
	ErrGroupExists   = errors.New("Group already exists.")
	ErrGroupNotExist = errors.New("Group does not exist.")
	ErrPageExists    = errors.New("Page already exists.")
	ErrPageNotExist  = errors.New("Page does not exist.")

	ErrConcurrentEdit = errors.New("Concurrent modification of page.")

	ErrInvalidSlug = errors.New("Invalid slug.")
)

type Database interface {
	Context(user Slug) Context
}

type Params map[string]interface{}

type Context interface {
	ActiveUserID() Slug
	Access() Access
	Users() Users
	Groups() Groups
	Index(user Slug) Index
	Pages(group Slug) Pages

	GuestLogin() GuestLogin
}

type Rights string

const (
	Blocked   Rights = "blocked"
	Reader           = "reader"
	Editor           = "editor"
	Moderator        = "moderator"
)

func (r Rights) Level() int {
	switch r {
	case Blocked:
		return 0
	case Reader:
		return 1
	case Editor:
		return 2
	case Moderator:
		return 3
	}
	return -1
}

type Access interface {
	VerifyUser(user User) error

	IsAdmin(user Slug) bool
	SetAdmin(user Slug, isAdmin bool) error

	Rights(group, user Slug) Rights

	// member is either a User or a Group
	AddUser(group, user Slug) error
	RemoveUser(group, user Slug) error

	CommunityAdd(group, member Slug, rights Rights) error
	CommunityRemove(group, member Slug) error

	List(group Slug) ([]Member, error)
}

type GuestLogin interface {
	Add(name, email, password string) error
	// implement auth.Provider
	Boot() template.HTML
	Verify(name, password string) (User, error)
}

type Users interface {
	ByID(id Slug) (User, error)
	Create(user User) error
	Delete(id Slug) error
	List() ([]User, error)
}

type Groups interface {
	ByID(id Slug) (Group, error)
	Create(group Group) error
	Delete(id Slug) error
	List() ([]Group, error)
}

type Pages interface {
	Create(page *Page) error

	Load(id Slug) (*Page, error)
	LoadRaw(id Slug) ([]byte, error)

	Overwrite(id Slug, version int, page *Page) error
	Edit(id Slug, version int, action Action) error
	Delete(id Slug, version int) error

	BatchReplace(pages map[Slug]*Page, complete func(Slug)) error

	List() ([]PageEntry, error)
	Journal(id Slug) ([]Action, error)
}

type Index interface {
	List() ([]PageEntry, error)

	Search(text string) ([]PageEntry, error)
	SearchFilter(text, exclude, include string) ([]PageEntry, error)

	Tags() ([]TagEntry, error)
	ByTag(tag Slug) ([]PageEntry, error)
	ByTagFilter(tag Slug, exclude, include string) ([]PageEntry, error)

	Groups(min Rights) ([]Group, error)
	ByGroup(groupID Slug) ([]PageEntry, error)

	ByTitle(title Slug) ([]PageEntry, error)

	RecentChanges() ([]PageEntry, error)
}

func init() { gob.Register(User{}) }

type User struct {
	ID        Slug   `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Company   string `json:"company"`
	Admin     bool   `json:"admin"`
	MaxAccess Rights `json:"-"`

	AuthID       string `json:"-"`
	AuthProvider string `json:"-"`
}

type Group struct {
	ID      Slug
	OwnerID Slug
	Name    string
	Public  bool

	Description string
}

type Member struct {
	ID      Slug
	Name    string
	IsGroup bool
	Access  Rights
}

func (g *Group) IsCommunity() bool { return g.ID == g.OwnerID }
