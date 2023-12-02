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

type Search struct {
	Query   string
	Updates []Update
}

type Update struct {
	Uri   string `bson:"uri"`
	Title string `bson:"title"`
	Desc  string `bson:"desc"`
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

func HandleSearch(w http.ResponseWriter, r *http.Request) {
	var search Search

	query := r.URL.Query().Get("q")
	search.Query = query
	search.Updates = getUpdates(query)

	tmpl, _ := template.ParseFiles("templates/search.html")
	tmpl.Execute(w, search)
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

func getUpdates(query string) []Update {
	var updates []Update

	url := os.Getenv("MONGODB")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return updates
	}
	coll := client.Database("moixllik").Collection("docs")

	opts := options.Find().SetLimit(7)
	opts.SetSort(bson.D{{"modified", -1}})
	opts.SetProjection(bson.D{{"_id", 0}, {"uri", 1}, {"title", 1}, {"desc", 1}})

	filter := bson.D{{"public", true}}
	if len(query) > 0 {
		switch query[0:1] {
		case "@":
			filter = bson.D{{"public", true}, {"authors", bson.D{{"$in", bson.A{
				query[1:],
			}}}}}
		case "#":
			filter = bson.D{{"public", true}, {"tags", bson.D{{"$in", bson.A{
				query[1:],
			}}}}}
		default:
			q_string := bson.D{{"$regex", query}, {"$options", "i"}}
			filter = bson.D{{"public", true}, {"$or", bson.A{
				bson.D{{"uri", q_string}},
				bson.D{{"title", q_string}},
				bson.D{{"desc", q_string}},
			}}}
		}
	}

	cursor, err := coll.Find(context.TODO(), filter, opts)
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
