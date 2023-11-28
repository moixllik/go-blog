package home

import (
	"net/http"
	"html/template"
)

type PageData struct {
	PageTitle string
	PageDesc  string
}

func Handle(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/home.html")
	pageData := PageData{
		PageTitle: "Moixllik",
		PageDesc:  "Computación, Contabilidad y Arte",
	}
	tmpl.Execute(w, pageData)
}
