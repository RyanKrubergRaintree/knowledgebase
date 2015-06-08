package auth

import (
	"crypto/rand"
	"fmt"
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

var authPath = ""

var _ kb.Auth = &Auth{}

type Auth struct{}

func New() *Auth { return &Auth{} }

func setProvider(r *http.Request, provider string) {
	q := r.URL.Query()
	q.Add(":provider", provider)
	r.URL.RawQuery = q.Encode()
}

func (auth *Auth) Start(w http.ResponseWriter, r *http.Request) {
	setProvider(r, path.Base(r.URL.Path))
	gothic.BeginAuthHandler(w, r)
}

func (auth *Auth) Finish(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	provider := path.Base(r.URL.Path)
	setProvider(r, provider)

	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		return kb.User{}, err
	}

	return kb.User{
		AuthID:       user.UserID,
		AuthProvider: provider,

		ID:    kb.Slugify(user.Name),
		Email: user.Email,
		Name:  user.Name,
	}, nil
}

func (auth *Auth) Logins() (logins []kb.AuthLogin) {
	for _, provider := range goth.GetProviders() {
		name := provider.Name()
		logins = append(logins, kb.AuthLogin{
			URL:  path.Join(authPath, "provider", name),
			Name: displayName[name],
			Icon: iconName[name],
		})
	}
	return
}

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

var iconName = map[string]string{
	"twitter":  "mdi mdi-twitter",
	"facebook": "mdi mdi-facebook",
	"gplus":    "mdi mdi-google",
	"github":   "mdi mdi-github-circle",
	"linkedin": "mdi mdi-linkedIn",
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
func Register(appkey string, url string, authPrefix string, clients map[string]Client) {
	if authPath != "" {
		panic("authentication path already set")
	}
	authPath = authPrefix

	if appkey != "" {
		gothic.AppKey = appkey
	} else {
		var key [32]byte
		_, err := rand.Read(key[:])
		if err != nil {
			panic(err)
		}
		gothic.AppKey = fmt.Sprintf("%64x", key)
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
