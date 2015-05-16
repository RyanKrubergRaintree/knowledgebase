package auth

import (
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/raintreeinc/knowledgebase/kbserver"
)

const authPath = "/system/auth"

type Front struct {
	kbserver.Presenter
	Context kbserver.Context
	Server  http.Handler
}

func New(server http.Handler, context kbserver.Context, presenter kbserver.Presenter) *Front {
	return &Front{
		Presenter: presenter,
		Context:   context,
		Server:    server,
	}
}

func (front *Front) login(w http.ResponseWriter, r *http.Request) {
	err := front.Present(w, r, "login.html", map[string]interface{}{
		"Logins": getLogins(),
	})
	if err != nil {
		log.Println(err)
	}
}

func (front *Front) forbidden(w http.ResponseWriter, r *http.Request) {
	message := front.Context.Once(w, r, "error")
	err := front.Present(w, r, "forbidden.html", map[string]string{
		"Message": message,
	})
	if err != nil {
		log.Println(err)
	}
}

func (front *Front) logout(w http.ResponseWriter, r *http.Request) {
	front.Context.Logout(w, r)
	http.Redirect(w, r, "/", http.StatusFound)
}

func (front *Front) Serve(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, authPath)
	switch {
	case path == "/login":
		front.login(w, r)
		return
	case path == "/forbidden":
		front.forbidden(w, r)
		return
	case path == "/logout":
		front.logout(w, r)
		return
	case strings.HasPrefix(path, "/provider/"):
		front.redirect(w, r)
		return
	case strings.HasPrefix(path, "/callback/"):
		front.callback(w, r)
		return
	}

	http.NotFound(w, r)
}

func (front *Front) RequireLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, path.Join(authPath, "login"), http.StatusFound)
}

func (front *Front) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, authPath) {
		front.Serve(w, r)
		return
	}

	if _, err := front.Context.CurrentUser(w, r); err != nil {
		front.Context.Add(w, r, "after-login", r.URL.String())
		front.RequireLogin(w, r)
		return
	}

	front.Server.ServeHTTP(w, r)
}
