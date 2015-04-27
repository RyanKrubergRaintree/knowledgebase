package farm

import (
	"html/template"
	"log"
	"net/http"
)

type renderer struct {
	template *template.Template
}

func NewRenderer(glob string) (Renderer, error) {
	t, err := template.New("").ParseGlob(glob)
	return &renderer{t}, err
}

func (r *renderer) Render(w http.ResponseWriter, name string, data interface{}) {
	if err := r.template.ExecuteTemplate(w, name, data); err != nil {
		log.Println(err)
	}
}
