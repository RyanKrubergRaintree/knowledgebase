package kbserver

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Files struct {
	dir string
}

func NewFiles(dir string) *Files {
	return &Files{
		dir: dir,
	}
}

func (a *Files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, SafeFile(a.dir, r.URL.Path))
}

type presenter struct {
	Dir      string
	Glob     string
	SiteInfo interface{}
	Context  Context
	Source   *Source
}

func NewPresenter(dir, glob string, siteinfo interface{}, source *Source, context Context) Presenter {
	return &presenter{
		Dir:      dir,
		Glob:     glob,
		SiteInfo: siteinfo,
		Context:  context,
		Source:   source,
	}
}

func (a *presenter) Present(w http.ResponseWriter, r *http.Request, tname string, data interface{}) error {
	ts, err := template.New("").Funcs(
		template.FuncMap{
			"Site": func() interface{} { return a.SiteInfo },
			"User": func() kb.User {
				user, _ := a.Context.CurrentUser(w, r)
				return user
			},
			"SourceFiles": func() []string { return a.Source.Files() },
		},
	).ParseGlob(filepath.Join(a.Dir, a.Glob))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if err := ts.ExecuteTemplate(w, tname, data); err != nil {
		return err
	}
	return nil
}
