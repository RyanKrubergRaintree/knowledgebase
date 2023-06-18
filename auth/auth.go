package auth

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/raintreeinc/knowledgebase/auth/session"
	"github.com/raintreeinc/knowledgebase/kb"
)

var (
	ErrUnauthorized = errors.New("Unauthorized")
	ErrTimeSkewed   = errors.New("TimeSkewed")
)

type Provider interface {
	Boot() template.HTML
	Verify(user, pass string) (kb.User, error)
}

type Rules interface {
	Login(user kb.User, db kb.Database) error
}

type Server struct {
	Rules    Rules
	DB       kb.Database
	Provider map[string]Provider

	Sessions *session.Store
}

func NewServer(rules Rules, db kb.Database) *Server {
	return &Server{
		Rules:    rules,
		DB:       db,
		Provider: make(map[string]Provider),
		Sessions: session.NewStore(time.Hour),
	}
}

func (server *Server) params(w http.ResponseWriter, r *http.Request) (kb.User, session.Token, error) {
	token, err := session.TokenFromString(r.Header.Get("X-Auth-Token"))
	if err != nil {
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	// creates session for LMS user;
	// todo: remove token before returning response
	lmsToken, lmsTokenExists := os.LookupEnv("LMSTOKEN")
	if lmsTokenExists && token.String() == lmsToken && lmsToken != "" {
		lmsUser, err1 := server.DB.Context("admin").Users().ByID(kb.Slugify("lmsuser"))
		if err1 != nil {
			return kb.User{}, session.ZeroToken, ErrUnauthorized
		}
		token, err = server.Sessions.New(lmsUser) // register the user in the store
		if err != nil {
			log.Println("Session failed:", err.Error())
			return kb.User{}, session.ZeroToken, ErrUnauthorized
		}
	}

	user, ok := server.Sessions.Load(token)
	if !ok {
		return kb.User{}, token, ErrUnauthorized
	}

	return user, token, nil
}

func (server *Server) Verify(w http.ResponseWriter, r *http.Request) (kb.User, error) {
	user, _, err := server.params(w, r)
	return user, err
}

type SessionInfo struct {
	Token string  `json:"token"`
	User  kb.User `json:"user"`

	Pages  []kb.Slug         `json:"pages"`
	Params map[string]string `json:"params,omitempty"`
}

func (server *Server) info(w http.ResponseWriter, r *http.Request) {
	tokstring := r.FormValue("token")
	if tokstring == "" {
		tokstring = r.Header.Get("X-Auth-Token")
	}
	token, err := session.TokenFromString(tokstring)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusBadRequest)
		return
	}

	user, ok := server.Sessions.Load(token)
	if !ok {
		http.Error(w, "Session expired.", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&SessionInfo{
		Token: token.String(),
		User:  user,
	})
}

func (server *Server) logout(w http.ResponseWriter, r *http.Request) {
	tokstring := r.FormValue("token")
	if tokstring == "" {
		tokstring = r.Header.Get("X-Auth-Token")
	}
	token, err := session.TokenFromString(tokstring)
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusBadRequest)
		return
	}

	server.Sessions.Delete(token)
}

func (server *Server) login(providername, username, pass string) (user kb.User, token session.Token, err error) {
	provider, ok := server.Provider[providername]
	if !ok {
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	user, err = provider.Verify(username, pass)
	if err != nil {
		log.Println("Verification failed:", err.Error())
		if strings.Contains(err.Error(), "time skewed") {
			return kb.User{}, session.ZeroToken, ErrTimeSkewed
		}
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	if err := server.Rules.Login(user, server.DB); err != nil {
		log.Println("Login failed:", err.Error())
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	token, err = server.Sessions.New(user)
	if err != nil {
		log.Println("Session failed:", err.Error())
		return kb.User{}, session.ZeroToken, ErrUnauthorized
	}

	return user, token, nil
}

func (server *Server) SessionFromHeader(r *http.Request) (*SessionInfo, error) {
	auth := r.Header.Get("Authorization")

	args := strings.SplitN(auth, " ", 3)
	if len(args) != 3 {
		return nil, ErrUnauthorized
	}

	user, token, err := server.login(args[0], args[1], args[2])
	if err != nil {
		return nil, err
	}

	slugs := []kb.Slug{}
	params := make(map[string]string)
	if err := r.ParseForm(); err == nil {
		for name, val := range r.PostForm {
			if len(val) > 0 {
				params[name] = val[0]
			}
		}

		if pagelist := params["pages"]; pagelist != "" {
			tags := []kb.Slug{}
			for _, tag := range strings.Split(pagelist, "|") {
				tags = append(tags, kb.Slugify(tag))
			}

			filter := kb.Slugify(params["branch"])
			index := server.DB.Context(user.ID).Index(user.ID)
			entries, err := index.ByTagFilter(tags, "help-", "help-"+string(filter))
			if err == nil {
				kb.SortPageEntriesByRank(entries, tags)
				for _, entry := range entries {
					slugs = append(slugs, entry.Slug)
				}
			}

			if len(entries) == 0 && filter != "" {
				slugs = append(slugs, "help-"+filter+"=welcome-to-raintree-help?className=no-results")
				slugs = append(slugs, "help-"+filter+"=index")
			}
		}
	}

	if r.RequestURI == "/apilogin" {
		server.Sessions.SetPageToShowAfterLogin(token, slugs)
	}

	return &SessionInfo{
		User:   user,
		Token:  token.String(),
		Pages:  slugs,
		Params: params,
	}, nil
}

func (server *Server) SessionFromToken(r *http.Request) (*SessionInfo, error) {
	token, err := session.TokenFromString(r.Header.Get("X-Auth-Token"))
	if err != nil {
		return nil, err
	}

	user, ok := server.Sessions.Load(token)
	if !ok {
		return nil, errors.New("Session expired.")
	}

	slugs := server.Sessions.GetPageToShowAfterLogin(token)
	// prevent re-using the login token by re-generating it
	server.Sessions.Delete(token)
	token, err = server.Sessions.New(user)
	if err != nil {
		return nil, errors.New("Could not re-generate the token.")
	}
	server.Sessions.SetPageToShowAfterLogin(token, slugs)
	params := make(map[string]string)

	slugs = server.Sessions.GetPageToShowAfterLogin(token)

	return &SessionInfo{
		User:   user,
		Token:  token.String(),
		Pages:  slugs,
		Params: params,
	}, nil
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/info") {
		server.info(w, r)
		return
	}
	if strings.HasPrefix(r.URL.Path, "/logout") {
		server.logout(w, r)
		return
	}

	providername := strings.TrimPrefix(r.URL.Path, "/")
	user, token, err := server.login(
		providername,
		r.FormValue("user"),
		r.FormValue("code"),
	)

	if err != nil {
		if err == ErrUnauthorized {
			http.Error(w, "Access denied.", http.StatusForbidden)
			return
		}

		if err == ErrTimeSkewed {
			http.Error(w, "Time skewed, authentication failed.", http.StatusUnauthorized)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(&SessionInfo{
		Token: token.String(),
		User:  user,
	})
}
