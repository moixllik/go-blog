package main

import (
	"bufio"
	"net/http"
	"os"
	"strings"

	"app/routes/home"
)

func init() {
	env, err := os.Open(".env")
	if err != nil {
		return
	}
	defer env.Close()

	scanner := bufio.NewScanner(env)
	for scanner.Scan() {
		line := scanner.Text()
		idx := strings.Index(line, "=")
		os.Setenv(line[:idx], line[idx+1:])
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/home", home.Handle)

	// Static
	mux.Handle("/", http.FileServer(http.Dir("public")))
	http.ListenAndServe(":8080", mux)
}
