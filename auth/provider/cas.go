package provider

import (
	"errors"
	"html/template"
	"net/url"
	"os"
	"strings"

	"github.com/raintreeinc/knowledgebase/auth/provider/trust"
	"github.com/raintreeinc/knowledgebase/kb"
)

type CAS struct {
	Provider string
	Key      []byte
}

func (conf *CAS) Boot() template.HTML { return "" }

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
	company, companyid := getCompanyIdRedirect(p.Get("company"), p.Get("companyid"))

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

// use different ID for a company if there is a redirect set up in the environment variables
func getCompanyIdRedirect(company, companyid string) (string, string) {
	key := strings.ToUpper(strings.Replace(company, " ", "_", -1))

	idRedirect, exists := os.LookupEnv(key)
	if exists {
		return company, idRedirect
	}

	return company, companyid
}
