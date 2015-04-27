package farm

import (
	"net/http"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data interface{})
}

type Auth interface {
	LoggedIn(http.ResponseWriter, *http.Request) bool
	RequestCredentials(http.ResponseWriter, *http.Request)
	GetUser(http.ResponseWriter, *http.Request) (kb.User, error)
	LogOut(http.ResponseWriter, *http.Request)

	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type Admin interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
