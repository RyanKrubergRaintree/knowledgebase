package assets

import (
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Files struct {
	dir  string
	exts []string
}

func NewFiles(dir string, allowedExts []string) *Files {
	return &Files{
		dir:  dir,
		exts: allowedExts,
	}
}

func (a *Files) allowed(url string) bool {
	ext := path.Ext(url)
	for _, v := range a.exts {
		if v == ext {
			return true
		}
	}
	return false
}

func (a *Files) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !a.allowed(r.URL.Path) {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	upath = path.Clean(upath)
	http.ServeFile(w, r, path.Join(a.dir, upath[1:]))
}

type Presenter struct {
	Dir      string
	Glob     string
	SiteInfo interface{}
	Context  kb.Context
}

func NewPresenter(dir, glob string, siteinfo interface{}, context kb.Context) *Presenter {
	return &Presenter{
		Dir:      dir,
		Glob:     glob,
		SiteInfo: siteinfo,
		Context:  context,
	}
}

func (a *Presenter) Present(w http.ResponseWriter, r *http.Request, tname string, data interface{}) error {
	ts, err := template.New("").Funcs(
		template.FuncMap{
			"Site": func() interface{} { return a.SiteInfo },
			"User": func() kb.User {
				user, _ := a.Context.CurrentUser(w, r)
				return user
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
