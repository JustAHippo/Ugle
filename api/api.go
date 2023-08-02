package api

import (
	"bufio"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

type SearchResponse struct {
	ErrorMsg string   `json:error`
	Results  []string `json:results`
}

func ApiSearch(ctx *gin.Context) {
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
		}

		line++

	}
	if len(queryResults) == 0 {
		ctx.JSON(404, SearchResponse{ErrorMsg: "No results for search!"})
		return
	}
	ctx.JSON(200, SearchResponse{Results: queryResults})
}
