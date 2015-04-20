package kb

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/egonelbre/fedwiki"
	"github.com/egonelbre/fedwiki/pagestore"
	"github.com/egonelbre/fedwiki/pagestore/folderstore"
	"github.com/egonelbre/fedwiki/pagestore/mongostore"

	"github.com/raintreeinc/knowledgebase/kb/pageindex"
	"github.com/raintreeinc/knowledgebase/kb/pageindex/memindex"
	"github.com/raintreeinc/knowledgebase/kb/pageindex/mongoindex"
)

type Server struct {
	Domain   string
	Provider string
	Pages    http.Handler
	Index    http.Handler
}

// Currently supports following databases:
// folder://folder
// mongodb://user:pass@host1:port1,host2:port2/database?options
func NewStore(database, provider string) (fedwiki.PageStore, pageindex.Index, error) {
	tokens := strings.SplitN(database, "://", 2)
	if len(tokens) != 2 {
		return nil, nil, errors.New("invalid database specification")
	}

	scheme, params := tokens[0], tokens[1]
	switch scheme {
	case "mongodb":
		store, err := mongostore.New(params, provider)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to mongo: %s", err)
		}

		index, err := mongoindex.New(params, provider)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to connect to index: %s", err)
		}

		if err := index.Init(); err != nil {
			return nil, nil, fmt.Errorf("failed to initialize index: %s", err)
		}

		return store, index, nil
	case "folder":
		params = filepath.Join(params, provider)
		err := os.MkdirAll(params, 0755)
		if err != nil && !os.IsExist(err) {
			return nil, nil, fmt.Errorf("failed to create folder: %s", err)
		}
		store := folderstore.New(params)
		return store, memindex.New(store, 1*time.Minute), nil
	}

	panic("unknown database engine " + scheme)
}

func NewServer(provider, domain, database string) (*Server, error) {
	mainstore, mainindex, err := NewStore(database, provider)
	if err != nil {
		return nil, err
	}
	return &Server{
		Domain:   domain,
		Provider: provider,
		Pages:    &fedwiki.Server{Handler: pagestore.Handler{mainstore}},
		Index:    &fedwiki.Server{Handler: pageindex.Handler{mainindex}},
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/index/") {
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/index")
		s.Index.ServeHTTP(w, r)
		return
	}
	s.Pages.ServeHTTP(w, r)
}
