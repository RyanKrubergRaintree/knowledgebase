package oauth

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/raintreeinc/knowledgebase/kb"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"github.com/markbates/goth/providers/linkedin"
	"github.com/markbates/goth/providers/twitter"
)

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

func (ctx *Context) RegisterProviders() {
	gothic.AppKey = os.Getenv("APPKEY")

	prefix := ctx.CallbackURL + "/"
	if os.Getenv("TWITTER_KEY") != "" {
		goth.UseProviders(twitter.New(os.Getenv("TWITTER_KEY"), os.Getenv("TWITTER_SECRET"), prefix+"twitter"))
	}
	if os.Getenv("FACEBOOK_KEY") != "" {
		goth.UseProviders(facebook.New(os.Getenv("FACEBOOK_KEY"), os.Getenv("FACEBOOK_SECRET"), prefix+"facebook"))
	}
	if os.Getenv("GPLUS_KEY") != "" {
		goth.UseProviders(gplus.New(os.Getenv("GPLUS_KEY"), os.Getenv("GPLUS_SECRET"), prefix+"gplus"))
	}
	if os.Getenv("GITHUB_KEY") != "" {
		goth.UseProviders(github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), prefix+"github"))
	}
	if os.Getenv("LINKEDIN_KEY") != "" {
		goth.UseProviders(linkedin.New(os.Getenv("LINKEDIN_KEY"), os.Getenv("LINKEDIN_SECRET"), prefix+"linkedin"))
	}
}

func (ctx *Context) RequestCredentials(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	r.URL.Host = r.Host
	r.URL.Fragment = ""
	q.Add("next", r.URL.String())

	u, err := url.Parse(ctx.LoginURL)
	if err != nil {
		http.Error(w, "Invalid Login URL", http.StatusInternalServerError)
		return
	}
	u.RawQuery = q.Encode()

	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func (ctx *Context) loginPage(w http.ResponseWriter, r *http.Request) {
	//TODO: better login page
	next := r.URL.Query().Get("next")
	if next != "" {
		session, _ := ctx.session(r)
		session.Values["next"] = next
		session.Save(r, w)
	}
	fmt.Fprint(w, `<a href="/provider/gplus">Google</a>`)
}

func (ctx *Context) provider(w http.ResponseWriter, r *http.Request) {
	setProvider(r, path.Base(r.URL.Path))
	gothic.BeginAuthHandler(w, r)
}

func (ctx *Context) callback(w http.ResponseWriter, r *http.Request) {
	provider := path.Base(r.URL.Path)
	setProvider(r, provider)
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		//TODO: proper error page
		fmt.Fprintln(w, err)
		return
	}

	ctx.Authenticated(w, r, kb.User{
		ID:       user.UserID,
		Email:    user.Email,
		Name:     user.Name,
		Provider: provider,
	})
}
