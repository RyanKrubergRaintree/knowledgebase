package kb

import "encoding/gob"

type User struct {
	AuthID   string
	ID       Slug
	Email    string
	Name     string
	Provider string
}

func init() {
	gob.Register(User{})
}
