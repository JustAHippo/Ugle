package api

import (
	"Ugle/db"
	"bufio"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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

func ApiSearch(ctx *gin.Context) {
	debug := false
	searchQuery := strings.ToLower(ctx.Query("q"))

	if searchQuery == "" {
		ctx.JSON(400, SearchResponse{ErrorMsg: "No query given!"})
		return
	} else if len(searchQuery) > 50 {
		ctx.JSON(400, SearchResponse{ErrorMsg: "Query length longer than 50!"})
		return
	}
	f, err := os.Open("./registry/registry.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	var queryResults []string
	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchQuery) && !strings.Contains(scanner.Text(), "0.0.0.0") {
			queryResults = append(queryResults, scanner.Text())
			if debug {
				requestURL := "http://" + strings.Split(scanner.Text(), " ")[0]
				response, _ := http.Get(requestURL)
				if response.StatusCode == 200 {
					document, _ := goquery.NewDocumentFromReader(response.Body)
					title := document.Find("title").Text()
					sitedesc, _ := document.Find("meta[name='description']").Attr("content")
					println("Title", title)
					println("Description", sitedesc)
					println("Inserting Cache to Database (currently unimplemented)")
				}
			}

		}

		line++

	}
	if len(queryResults) == 0 {
		ctx.JSON(404, SearchResponse{ErrorMsg: "No results for search!"})
		return
	}
	ctx.JSON(200, SearchResponse{Results: queryResults})
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
	for _, result := range results {
		//res, _ := json.Marshal(result)
		ctx.JSON(200, result)
	}
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
