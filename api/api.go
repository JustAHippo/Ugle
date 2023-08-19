package api

import (
	"Ugle/db"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"os"
	"strings"
)

type SearchResponse struct {
	ErrorMsg string   `json:error`
	Results  []string `json:results`
}

func MongoApiSearch(ctx *gin.Context) {
	searchQuery := strings.ToLower(ctx.Query("q"))

	if searchQuery == "" {
		ctx.JSON(400, SearchResponse{ErrorMsg: "No query given!"})
		return
	} else if len(searchQuery) > 50 {
		ctx.JSON(400, SearchResponse{ErrorMsg: "Query length longer than 50!"})
		return
	}
	_ = mongo.IndexModel{Keys: bson.D{{"description", "text"}, {"domain", "text"}, {"title", "text"}}}
	filter := bson.D{{"$text", bson.D{{"$search", fmt.Sprintf("\"%s\"", searchQuery)}}}}
	println(searchQuery)
	cursor, err := db.SiteDirectory.Find(context.TODO(), filter)
	if err != nil {
		panic(err)
	}
	var results []db.Site
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	ctx.JSON(200, results)
}

func ApiCache(ctx *gin.Context) {
	resp, err := http.Get("https://raw.githubusercontent.com/ucanet/ucanet-registry/main/ucanet-registry.txt")
	if err != nil {
		panic(err)
		return
	}
	defer resp.Body.Close()
	out, err := os.Create("./registry/registry.txt")
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}
