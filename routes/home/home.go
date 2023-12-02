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

	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo struct {
	Url  string
	Name string
	Desc string
}

type Update struct {
	Uri   string `bson:"uri"`
	Title string `bson:"title"`
	Desc  string `bson:"desc"`
}

type PageData struct {
	PageTitle string
	PageDesc  string
	Repos     []Repo
	Updates   []Update
}

func Handle(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("templates/home.html")
	pageData := PageData{
		PageTitle: "Moixllik",
		PageDesc:  "Computación, Contabilidad y Arte",
		Repos:     getRepos(),
		Updates:   getUpdates(),
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

func getUpdates() []Update {
	var updates []Update

	uri := os.Getenv("MONGODB")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return updates
	}
	coll := client.Database("moixllik").Collection("docs")

	opts := options.Find().SetLimit(5)
	opts.SetSort(bson.D{{"modified", -1}})
	opts.SetProjection(bson.D{{"_id", 0}, {"uri", 1}, {"title", 1}, {"desc", 1}})

	cursor, err := coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return updates
	}

	for cursor.Next(context.TODO()) {
		var result Update

		if err := cursor.Decode(&result); err != nil {
			continue
		}
		updates = append(updates, result)
	}
	return updates
}
