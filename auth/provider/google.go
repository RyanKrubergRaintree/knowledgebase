package provider

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Google struct {
	ClientID     string
	ClientSecret string
	HostedDomain string
}

func (conf *Google) Boot() template.HTML {
	var head string
	head += `<script src="https://apis.google.com/js/platform.js"></script>`
	head += `<meta name="google-signin-client_id" content="` +
		template.JSEscapeString(conf.ClientID) + `">`
	head += `<script>var GoogleHostedDomain="` + template.JSEscapeString(conf.HostedDomain) + `"</script>`
	return template.HTML(head)
}

func (conf *Google) Verify(user, code string) (kb.User, error) {
	//TODO: verify that token is valid
	token := code

	const tokinfo = "https://www.googleapis.com/oauth2/v3/tokeninfo"
	r, err := http.Get(tokinfo + "?id_token=" + url.QueryEscape(token))
	if err != nil {
		return kb.User{}, fmt.Errorf("Failed to get response for user \"%v\": %v", user, err)
	}

	var result struct {
		Email         string `json:"email,omitempty"`
		EmailVerified string `json:"email_verified"`

		Name           string `json:"name"`
		Issuer         string `json:"iss"`
		Audience       string `json:"aud"`
		Expiration     string `json:"exp"`
		IssuedAt       string `json:"iat"`
		UserIdentifier string `json:"sub"`
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return kb.User{}, fmt.Errorf("Failed to read result for user \"%v\": %v", user, err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return kb.User{}, fmt.Errorf("Invalid JSON for user \"%v\": %v", user, err)
	}

	if r.StatusCode != http.StatusOK {
		return kb.User{}, fmt.Errorf("Failed to authenticate with Google user \"%v\": %s", user, r.Status)
	}

	if result.EmailVerified != "true" ||
		result.Audience != conf.ClientID ||
		result.Email != user {
		return kb.User{}, fmt.Errorf("Invalid token for \"%v\"", user)
	}

	return kb.User{
		AuthID:       result.UserIdentifier,
		AuthProvider: "gplus",

		ID:    kb.Slugify(result.Name),
		Email: result.Email,
		Name:  result.Name,

		MaxAccess: kb.Moderator,
	}, nil
}
