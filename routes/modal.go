package routes

import (
	"html/template"
	"net/http"
)

type ModalParams struct {
	IsWinner bool
}

func Modal(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		tmpl := template.Must(template.ParseFiles("views/modal.html"))

		tmplParams := ModalParams{
			IsWinner: true,
		}

		err := tmpl.Execute(w, tmplParams)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
