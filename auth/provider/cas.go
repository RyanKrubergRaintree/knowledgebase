package provider

import (
	"errors"
	"net/url"

	"github.com/raintreeinc/knowledgebase/auth/provider/trust"
	"github.com/raintreeinc/knowledgebase/kb"
)

type CAS struct {
	Provider string
	Key      []byte
}

func (conf *CAS) Verify(params, code string) (kb.User, error) {
	id, err := trust.Peer{conf.Key}.Verify(code)
	if err != nil {
		return kb.User{}, err
	}

	p, err := url.ParseQuery(params)
	if err != nil {
		return kb.User{}, err
	}

	user := p.Get("user")
	company := p.Get("company")
	companyid := p.Get("companyid")
	if company+"="+user != id {
		return kb.User{}, errors.New("invalid id provided")
	}

	return kb.User{
		AuthID:       string(kb.Slugify(companyid + "=" + user)),
		AuthProvider: conf.Provider,

		ID:      kb.Slugify(id),
		Email:   "",
		Name:    id,
		Company: company,

		MaxAccess: kb.Moderator,
	}, nil
}
