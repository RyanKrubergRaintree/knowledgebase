package kb

import "encoding/gob"

type User struct {
	ID       string
	Email    string
	Name     string
	Provider string
}

func init() {
	gob.Register(User{})
}
