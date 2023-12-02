package main

import (
	"net/http"

	"app/routes/home"
	"app/routes/docs"
	"app/routes/tools"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("/home", home.Handle)
	mux.HandleFunc("/d/", docs.HandleReader)
	mux.HandleFunc("/search", docs.HandleSearch)
	mux.HandleFunc("/sitemap.txt", tools.Sitemap)

	// Static
	mux.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":8080", mux)
}
