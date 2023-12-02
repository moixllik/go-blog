package main

import (
	"net/http"

	"app/routes/home"
	"app/routes/tools"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	mux := http.NewServeMux()

	mux.HandleFunc("/home", home.Handle)
	mux.HandleFunc("/sitemap.txt", tools.Sitemap)

	// Static
	mux.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":8080", mux)
}
