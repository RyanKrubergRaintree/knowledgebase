package client

import (
	"encoding/base64"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/raintreeinc/livepkg"

	"github.com/raintreeinc/knowledgebase/kb"

	"github.com/raintreeinc/knowledgebase/auth"
)

var (
	errTimeSkewedMessage = "\n\n" +
		"Knowledge Base authentication failed due to incorrect computer time.\n" +
		"Accuracy within to 2 minutes is required.\n\n" +
		"Tip: Enable automatic time adjustment in your computer's settings."
)

type Info struct {
	Domain string

	ShortTitle string
	Title      string
	Company    string

	TrackingID string
	Version    string
}

type Server struct {
	Info
	Login *auth.Server

	development bool
	bootstrap   string
	dir         string
	assets      http.Handler
	client      http.Handler
}

func NewServer(info Info, login *auth.Server, dir string, development bool) *Server {
	return &Server{
		Info:  info,
		Login: login,

		development: development,
		bootstrap:   filepath.Join(dir, "index.html"),
		dir:         dir,
		assets: http.StripPrefix("/assets/",
			http.FileServer(http.Dir(filepath.Join(dir, "assets")))),
		client: livepkg.NewServer(
			http.Dir(dir),
			development,
			"/boot.js",
		),
	}
}

func (server *Server) index(w http.ResponseWriter, r *http.Request) {
	initialSession, initialSessionErr := server.Login.SessionFromHeader(r)
	if initialSessionErr == auth.ErrTimeSkewed {
		http.Error(w, errTimeSkewedMessage, http.StatusUnauthorized)
		return
	}

	CSPNonce := string(kb.AddCSPHeader(w))

	ts, err := template.New("").Funcs(
		template.FuncMap{
			"CSPNonce": func() template.HTMLAttr {
				return template.HTMLAttr(CSPNonce)
			},
			"Development": func() bool { return server.development },
			"Site":        func() Info { return server.Info },
			"InitialSession": func() template.JS {
				if initialSessionErr != nil {
					return "null"
				}

				data, _ := json.Marshal(initialSession)
				return template.JS(data)
			},
			"LoginProviders": func() interface{} {
				return server.Login.Provider
			},
		},
	).ParseGlob(server.bootstrap)

	if err != nil {
		log.Printf("Error parsing templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-UA-Compatible", "IE=edge")

	if err := ts.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Error executing template: %s", err)
		return
	}
}

func (server *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	kb.AddCommonHeaders(w)

	token := strings.TrimSpace(r.FormValue("token"))
	if token != "" {
		server.loginUsingTokenFromURL(w, r, token)
	} else {
		switch {
		case r.URL.Path == "/favicon.ico":
			http.ServeFile(w, r, filepath.Join(server.dir, "assets", "ico", "favicon.ico"))
		case r.URL.Path == "/":
			server.index(w, r)
		case r.URL.Path == "/apilogin":
			server.apiLogin(w, r)
		case strings.HasPrefix(r.URL.Path, "/assets/"):
			server.assets.ServeHTTP(w, r)
		default:
			server.client.ServeHTTP(w, r)
		}
	}
}

// rtAgent requests a token that can be used to log in a Web Client session
func (server *Server) apiLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sess, err := server.Login.SessionFromHeader(r)
	if err == auth.ErrTimeSkewed {
		http.Error(w, errTimeSkewedMessage, http.StatusUnauthorized)
		return
	}

	if sess == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// token will be sent back via query param
	token := base64.URLEncoding.EncodeToString([]byte(sess.Token))

	//nolint:errcheck
	json.NewEncoder(w).Encode(token)
}

// Web Client logs in with the token obtained from the agent earlier
func (server *Server) loginUsingTokenFromURL(w http.ResponseWriter, r *http.Request, token string) {
	base64token, err1 := base64.URLEncoding.DecodeString(token)
	if err1 != nil {
		http.Error(w, "Invalid token.", http.StatusUnauthorized)
		return
	}

	r.Header.Set("X-Auth-Token", string(base64token))

	session, err := server.Login.SessionFromToken(r)
	if err == auth.ErrTimeSkewed {
		http.Error(w, errTimeSkewedMessage, http.StatusUnauthorized)
		return
	}

	if session == nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	CSPNonce := string(kb.AddCSPHeader(w))

	ts, err := template.New("").Funcs(
		template.FuncMap{
			"CSPNonce": func() template.HTMLAttr {
				return template.HTMLAttr(CSPNonce)
			},
			"Development": func() bool { return server.development },
			"Site":        func() Info { return server.Info },
			"InitialSession": func() template.JS {
				if err != nil {
					return "null"
				}

				data, _ := json.Marshal(session)
				return template.JS(data)
			},
			"LoginProviders": func() interface{} {
				return server.Login.Provider
			},
		},
	).ParseGlob(server.bootstrap)

	if err != nil {
		log.Printf("Error parsing templates: %s", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("X-UA-Compatible", "IE=edge")

	if err := ts.ExecuteTemplate(w, "index.html", nil); err != nil {
		log.Printf("Error executing template: %s", err)
		return
	}

}
