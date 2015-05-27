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

// TODO: merge this into server
type presenter struct {
	Dir      string
	Glob     string
	SiteInfo interface{}
	Context  Context
	Database Database
	Sources  Sources
}

func NewPresenter(dir, glob string, siteinfo interface{}, sources Sources, context Context, database Database) Presenter {
	return &presenter{
		Dir:      dir,
		Glob:     glob,
		SiteInfo: siteinfo,
		Context:  context,
		Sources:  sources,
		Database: database,
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
			"UserGroups": func() []string {
				user, _ := a.Context.CurrentUser(w, r)
				info, _ := a.Database.Users().ByID(user.ID)
				return info.Groups
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
