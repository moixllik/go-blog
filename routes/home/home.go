package home

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type Repo struct {
	Name string
	Url  string
	Desc string
}

type PageData struct {
	PageTitle string
	PageDesc  string
	Repos     []Repo
}

func Handle(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/home.html")
	pageData := PageData{
		PageTitle: "Moixllik",
		PageDesc:  "Computación, Contabilidad y Arte",
		Repos:     getRepos(),
	}
	tmpl.Execute(w, pageData)
}

func getRepos() []Repo {
	var repos []Repo
	tmp := filepath.Join(os.TempDir(), "repos"+fmt.Sprint(time.Now().Day()))
	_, err := os.Stat(tmp)
	if os.IsNotExist(err) {
		token := os.Getenv("GITHUB")
		req, err := http.NewRequest("GET", "https://api.github.com/users/moixllik/repos", nil)
		if err != nil {
			return repos
		}
		req.Header.Add("User-Agent", "Mozilla/5.0")
		req.Header.Add("Authorization", "Bearer "+token)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return repos
		}
		defer resp.Body.Close()

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return repos
		}
		err = os.WriteFile(tmp, b, 0666)
		if err != nil {
			return repos
		}
	}

	b, err := os.ReadFile(tmp)
	if err != nil {
		return repos
	}
	var data []map[string]interface{}
	if err := json.Unmarshal(b, &data); err != nil {
		return repos
	}
	for _, it := range data {
		if it["fork"] == false && it["archived"] == false {
			var repo = new(Repo)
			repo.Name = it["name"].(string)
			repo.Url = it["html_url"].(string)
			repo.Desc = it["description"].(string)

			repos = append(repos, *repo)
		}
	}
	return repos
}
