package provider

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Google struct {
	ClientID     string
	ClientSecret string
}

func (conf *Google) Info() map[string]string {
	return map[string]string{
		"kind":     "google",
		"clientID": conf.ClientID,
	}
}

func (conf *Google) Verify(user, code string) (kb.User, error) {
	//TODO: verify that token is valid
	token := code

	const tokinfo = "https://www.googleapis.com/oauth2/v3/tokeninfo"
	r, err := http.Get(tokinfo + "?id_token=" + url.QueryEscape(token))
	if err != nil {
		return kb.User{}, fmt.Errorf("Failed to get response: %v", err)
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
		return kb.User{}, fmt.Errorf("Failed to read result: %v", err)
	}

	err = json.Unmarshal(data, &result)
	if err != nil {
		return kb.User{}, fmt.Errorf("Invalid result: %v", err)
	}

	if r.StatusCode != http.StatusOK {
		return kb.User{}, fmt.Errorf("Failed to authenticate with Google: %s")
	}

	if result.EmailVerified != "true" ||
		result.Audience != conf.ClientID ||
		result.Email != user {
		return kb.User{}, errors.New("Invalid token")
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
