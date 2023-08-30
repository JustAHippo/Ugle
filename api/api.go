package api

import (
	"Ugle/db"
	"context"
	"fmt"
	"time"
	"html"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"unicode"
	"strings"
)

type SearchResponse struct {
	ErrorMsg string   `json:error`
	Results  []string `json:results`
}

func EllipticalTruncate(text string, maxLen int) string {
	lastSpaceIx := maxLen
	len := 0
	for i, r := range text {
		if unicode.IsSpace(r) {
			lastSpaceIx = i
		}
		len++
		if len > maxLen {
			return text[:lastSpaceIx] + "..."
		}
	}
	return text
}

func MongoApiSearch(ctx *gin.Context) {
	startTime := time.Now()
	page, ferr := os.ReadFile("./static/search/index.html")
	if ferr != nil {
		panic(ferr)
	}
	searchQuery := strings.ToLower(ctx.Query("q"))
	numAmount := strings.ToLower(ctx.Query("num"))
	feelingLucky := strings.ToLower(ctx.Query("lucky"))
	
	var realNum = 10
	var s10 = ""
	var s30 = ""
	var s100 = ""
	if numAmount == "30" {
		realNum = 30
		s30 = "selected=\"\""
	} else if numAmount == "100" {
		realNum = 100
		s100 = "selected=\"\""
	} else {
		s10 = "selected=\"\""
	}
	
	var searchResults = ""
	var displayed = 0
	var total = 0
	if searchQuery == "" {
		searchResults = "No query given!"
	} else if len(searchQuery) > 50 {
		searchResults = "Query length longer than 50!"
	} else {
		idxopts := options.CreateIndexes().SetMaxTime(10 * time.Second)
		db.SiteDirectory.Indexes().CreateOne(context.TODO(), mongo.IndexModel{Keys: bson.D{{"description", "text"}, {"domain", "text"}, {"title", "text"}, {"specifieddescription", "text"}, {"specifiedtags", "text"}}}, idxopts)
		filter := bson.D{{"$text", bson.D{{"$search", fmt.Sprintf("\"%s\"", searchQuery)}}}}

		cursor, err := db.SiteDirectory.Find(context.TODO(), filter)
		if err != nil {
			panic(err)
		}
		var results []db.Site
		if err = cursor.All(context.TODO(), &results); err != nil {
			panic(err)
		}
		
		var entry = "<p><a href=\"http://%s\">%s</a><font size=\"-1\"><br>%s<br><font color=\"#008000\">%s</font></font></p>"
		for index, csite := range results {
			total++
			if index < realNum {
				if feelingLucky != "" {
					ctx.Redirect(302, "http://" + csite.Domain)
					return
				}
				displayed++
				searchResults += fmt.Sprintf(entry, csite.Domain, html.EscapeString(EllipticalTruncate(csite.Title, 100)), html.EscapeString(EllipticalTruncate(csite.Description, 200)), csite.Domain)
			}
		}
	}
	
	var pagestr = fmt.Sprintf(string(page), html.EscapeString(searchQuery), s10, s30, s100, displayed, total, html.EscapeString(searchQuery), time.Since(startTime), searchResults)
	ctx.Data(200, "text/html; charset=utf-8", []byte(pagestr))
}
