package docs

import (
	"html/template"
	"net/http"
	"os"

	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Doc struct {
	Uri      string   `bson:"uri"`
	Title    string   `bson:"title"`
	Desc     string   `bson:"desc"`
	Content  string   `bson:"content"`
	Extra    string   `bson:"extra"`
	Modified string   `bson:"modified"`
	Tags     []string `bson:"tags"`
	Authors  []string `bson:"authors"`
}

func HandleReader(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path[3:]
	doc := getDoc(uri)
	if doc.Title == "" {
		tmpl, _ := template.ParseFiles("templates/404.html")
		w.WriteHeader(404)
		tmpl.Execute(w, nil)
		return
	}
	tmpl, _ := template.ParseFiles("templates/reader.html")
	tmpl.Execute(w, doc)
}

func getDoc(uri string) Doc {
	var doc Doc

	url := os.Getenv("MONGODB")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return doc
	}
	coll := client.Database("moixllik").Collection("docs")

	err = coll.FindOne(context.TODO(), bson.D{{"public", true}, {"uri", uri}}).Decode(&doc)
	if err != nil {
		return doc
	}
	return doc
}
