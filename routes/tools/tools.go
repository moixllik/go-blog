package tools

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DocUri struct {
	Uri string `bson:"uri"`
}

func Sitemap(w http.ResponseWriter, r *http.Request) {
	uris := getUris()
	fmt.Fprintf(w, strings.Join(uris, "\n"))
}

func getUris() []string {
	var uris []string
	domain := "https://www.cix.ovh"

	url := os.Getenv("MONGODB")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		return uris
	}
	coll := client.Database("moixllik").Collection("docs")

	opts := options.Find()
	opts.SetProjection(bson.D{{"_id", 0}, {"uri", 1}})

	cursor, err := coll.Find(context.TODO(), bson.D{{"public", true}}, opts)
	if err != nil {
		return uris
	}

	for cursor.Next(context.TODO()) {
		var result DocUri

		if err := cursor.Decode(&result); err != nil {
			continue
		}
		uris = append(uris, domain+"/d/"+result.Uri)
	}

	return uris
}
