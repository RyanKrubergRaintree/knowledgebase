package farm

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

type renderer struct {
	template *template.Template
}

func NewRenderer(config Config) (Renderer, error) {
	glob := filepath.Join(config.AssetsDir, "templates", "*")
	t, err := template.New("").Funcs(template.FuncMap{
		"SiteTitle": func() string { return config.Site },
	}).ParseGlob(glob)
	return &renderer{t}, err
}

func (r *renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	if err := r.template.ExecuteTemplate(w, name, data); err != nil {
		log.Println(err)
	}
}
