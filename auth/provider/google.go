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
	Host         string
}

func (conf *Google) Boot() template.HTML {
	client_id := template.JSEscapeString(conf.ClientID)
	login_uri := template.JSEscapeString(conf.Host)
	hd := template.JSEscapeString(conf.HostedDomain)

	if login_uri != "" {
		login_uri = "https://" + login_uri
	}
	login_uri += "/system/auth/google"

	var head string
	// https://developers.google.com/identity/gsi/web/guides/client-library
	head += `<script src="https://accounts.google.com/gsi/client"></script>`
	head += `
		<script id="gsi_information" type="application/json">
			{
				"hd": "` + hd + `",
				"client_id": "` + client_id + `",
				"login_uri": "` + login_uri + `"
			}
		</script>
	`
	return template.HTML(head)
}

func (conf *Google) Verify(user, code string) (kb.User, error) {
	//TODO: verify that token is valid
	token := code

	const tokinfo = "https://www.googleapis.com/oauth2/v3/tokeninfo"
	r, err := http.Get(tokinfo + "?id_token=" + url.QueryEscape(token))
	if err != nil {
		return kb.User{}, fmt.Errorf("failed to get response for user \"%v\": %v", user, err)
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
		return kb.User{}, fmt.Errorf("failed to read result for user \"%v\": %v", user, err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return kb.User{}, fmt.Errorf("invalid JSON for user \"%v\": %v", user, err)
	}

	if r.StatusCode != http.StatusOK {
		return kb.User{}, fmt.Errorf("failed to authenticate with Google user \"%v\": %s", user, r.Status)
	}

	if result.EmailVerified != "true" ||
		result.Audience != conf.ClientID ||
		result.Email != user {
		return kb.User{}, fmt.Errorf("invalid token for \"%v\"", user)
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
