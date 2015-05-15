package auth

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/twitter"

	"github.com/raintreeinc/knowledgebase/kb"
)

type register func(key, secret, callback string)

var providers = map[string]register{
	"twitter":  func(key, secret, cb string) { goth.UseProviders(twitter.New(key, secret, cb)) },
	"facebook": func(key, secret, cb string) { goth.UseProviders(facebook.New(key, secret, cb)) },
	"gplus":    func(key, secret, cb string) { goth.UseProviders(gplus.New(key, secret, cb)) },
	"github":   func(key, secret, cb string) { goth.UseProviders(github.New(key, secret, cb)) },
	"linkedin": func(key, secret, cb string) { goth.UseProviders(linkedin.New(key, secret, cb)) },
}

var displayName = map[string]string{
	"twitter":  "Twitter",
	"facebook": "Facebook",
	"gplus":    "Google",
	"github":   "Github",
	"linkedin": "LinkedIn",
}

type Client struct{ Key, Secret string }

func ClientsFromEnv() map[string]Client {
	clients := make(map[string]Client)
	for name := range providers {
		prefix := strings.ToUpper(name)
		key, secret := os.Getenv(prefix+"_KEY"), os.Getenv(prefix+"_SECRET")
		if key != "" && secret != "" {
			clients[name] = Client{key, secret}
		}
	}
	return clients
}

// Register registers all oauth providers
//
// Supported providers:
//   twitter
//   facebook
//   gplus
//   github
//   linkedin
func Register(appkey string, url string, clients map[string]Client) {
	if appkey != "" {
		gothic.AppKey = appkey
	}

	cb := func(provider string) string {
		return url + path.Join(authPath, "callback", provider)
	}

	for name, client := range clients {
		provider := strings.ToLower(name)
		register, ok := providers[provider]
		if !ok {
			panic(fmt.Sprintf("Unimplemented provider %s", provider))
		}
		register(client.Key, client.Secret, cb(provider))
	}
}

type loginInfo struct {
	URL  string
	Name string
}

func getLogins() []loginInfo {
	logins := []loginInfo{}
	for _, provider := range goth.GetProviders() {
		name := provider.Name()
		logins = append(logins, loginInfo{
			URL:  path.Join(authPath, "provider", name),
			Name: displayName[name],
		})
	}
	return logins
}

func setProvider(r *http.Request, name string) {
	q := r.URL.Query()
	q.Add(":provider", name)
	r.URL.RawQuery = q.Encode()
}

func init() {
	gothic.GetState = func(r *http.Request) string {
		return r.URL.Query().Get("state")
	}
}

func (front *Front) redirect(w http.ResponseWriter, r *http.Request) {
	setProvider(r, path.Base(r.URL.Path))
	gothic.BeginAuthHandler(w, r)
}

func (front *Front) callback(w http.ResponseWriter, r *http.Request) {
	provider := path.Base(r.URL.Path)
	setProvider(r, provider)
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		log.Println(err)
		return
	}

	err = front.Context.Login(w, r, kb.User{
		ID:       user.UserID,
		Email:    user.Email,
		Name:     user.Name,
		Provider: provider,
	})

	if err != nil {
		front.Context.Add(w, r, "error", err.Error())
		http.Redirect(w, r, path.Join(authPath, "forbidden"), http.StatusFound)
		return
	}

	backto := front.Context.Once(w, r, "after-login")
	if backto == "" {
		backto = "/"
	}

	http.Redirect(w, r, backto, http.StatusFound)
}
