package routes

import (
	"html/template"
	"net/http"
)

type TemplateParams struct {
	Unauthorized bool
}

func Index(w http.ResponseWriter, r *http.Request) {
	tmplParams := TemplateParams{
		Unauthorized: true,
	}

	switch r.Method {

	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("views/index.html", "views/auth.html", "views/quiz.html"))

		err := tmpl.ExecuteTemplate(w, "index.html", tmplParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
