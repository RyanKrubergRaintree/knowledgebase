package kb

import "net/http"

type Auth interface {
	LoggedIn(http.ResponseWriter, *http.Request) bool
	RequestCredentials(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request) (User, error)
	LogOut(http.ResponseWriter, *http.Request)

	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Admin interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
