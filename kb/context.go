package kb

import "net/http"

type Context interface {
	LoggedIn(http.ResponseWriter, *http.Request) bool
	RequestCredentials(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request) (User, error)
	LogOut(http.ResponseWriter, *http.Request)

	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
