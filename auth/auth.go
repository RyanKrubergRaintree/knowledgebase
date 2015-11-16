package auth

import (
	"errors"
	"net/http"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Rules interface {
	Login(user kb.User, db kb.Database) error
}

type Server struct {
	Rules Rules
	DB    kb.Database
}

func NewServer(rules Rules, db kb.Database) *Server {
	return &Server{
		Rules: rules,
		DB:    db,
	}
}

func (server *Server) Verify(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	return kb.User{}, errors.New("Invalid Password!")
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

/*
	// protect server with authentication
	/// url := "http://" + *domain
	/// auth.Register(os.Getenv("APPKEY"), url, "/system/auth", auth.ClientsFromEnv())

	sec := auth.NewServer()

	sec.Alternate["guest"] = auth.NewDB(db)
	if caskey := os.Getenv("CASKEY"); caskey != "" {
		data, err := base64.StdEncoding.DecodeString(caskey)
		if err != nil {
			log.Fatal(err)
		}
		sec.Alternate["community"] = auth.NewCAS("community", data)
	}

*/
