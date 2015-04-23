package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"
)

func (ctx *Context) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.URL.Path == "/":
		if ctx.LoggedIn(w, r) {
			ctx.userPage(w, r)
		} else {
			http.Redirect(w, r, ctx.LoginURL, http.StatusTemporaryRedirect)
		}
	case r.URL.Path == "/logout":
		ctx.LogOut(w, r)
	case r.URL.Path == "/login":
		ctx.loginPage(w, r)
	case strings.HasPrefix(r.URL.Path, "/provider/"):
		ctx.provider(w, r)
	case strings.HasPrefix(r.URL.Path, "/callback/"):
		ctx.callback(w, r)
	}
}

func (ctx *Context) Authenticated(w http.ResponseWriter, r *http.Request, user kb.User) {
	if !strings.HasSuffix(user.Email, "@raintreeinc.com") {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "<h1>Access allowed only for @raintreeinc.com gmail accounts.</h1>")
		return
	}

	session, err := ctx.session(r)
	if err != nil {
		log.Println(err)
	}
	session.Values["user"] = user
	next, nextok := session.Values["next"]
	delete(session.Values, "next")
	session.Save(r, w)

	if nextok {
		if url, ok := next.(string); ok {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (ctx *Context) userPage(w http.ResponseWriter, r *http.Request) {
	user, err := ctx.GetUser(w, r)
	if err != nil {
		log.Println(err)
	}
	//TODO: better user page
	fmt.Fprintf(w, "%#v", user)
}
