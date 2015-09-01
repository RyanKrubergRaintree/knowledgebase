package auth

import (
	"errors"
	"net/http"

	"github.com/raintreeinc/knowledgebase/auth/trust"
	"github.com/raintreeinc/knowledgebase/kb"
)

type CAS struct {
	Provider string
	Key      []byte
}

func NewCAS(provider string, key []byte) *CAS {
	return &CAS{provider, key}
}

func (cas *CAS) Start(w http.ResponseWriter, r *http.Request) { panic("unimplemented") }
func (cas *CAS) Logins() (logins []kb.AuthLogin)              { panic("unimplemented") }

func (cas *CAS) Finish(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	id, err := trust.Peer{cas.Key}.Verify(r)
	if err != nil {
		return kb.User{}, err
	}

	// verify id
	user := r.FormValue("user")
	company := r.FormValue("company")

	if user+":"+company != id {
		return kb.User{}, errors.New("invalid id provided")
	}

	return kb.User{
		AuthID:       string(kb.Slugify(id)),
		AuthProvider: cas.Provider,

		ID:    kb.Slugify(id),
		Email: "",
		Name:  id,

		MaxAccess: kb.Editor,
	}, nil
}
