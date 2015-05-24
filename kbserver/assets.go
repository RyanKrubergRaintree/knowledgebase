package kbserver

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Sources interface {
	Include() string
}

type presenter struct {
	Dir      string
	Glob     string
	SiteInfo interface{}
	Context  Context
	Sources  Sources
}

func NewPresenter(dir, glob string, siteinfo interface{}, sources Sources, context Context) Presenter {
	return &presenter{
		Dir:      dir,
		Glob:     glob,
		SiteInfo: siteinfo,
		Context:  context,
		Sources:  sources,
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
			"Include": func() template.HTML {
				return template.HTML(a.Sources.Include())
			},
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
