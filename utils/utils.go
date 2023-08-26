package utils

import (
	"Ugle/db"
	"bufio"
	"context"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net"
	"os"
	"time"
	"strings"
)

func CreateDatabaseFromRegistry() {
	f, err := os.Open("./registry/registry.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner := bufio.NewScanner(f)
	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	fullDatabase := []interface{}{}
	for scanner.Scan() {
		domainName := strings.Split(scanner.Text(), " ")[0]
		ipAddress := strings.Split(scanner.Text(), " ")[2]
		
		if ipAddress != "0.0.0.0" {
			var requestURL = "http://"
			if net.ParseIP(ipAddress).To4() == nil {
				requestURL += "135.148.41.26"
			} else {
				requestURL += ipAddress
			}
			var siteTitle string
			var siteDescription string
			
			req, qerr := http.NewRequest("GET", requestURL, nil)
			if qerr != nil {
				panic(qerr)
			}
			req.Host = domainName
			client := &http.Client{Timeout: 4 * time.Second}
			response, err := client.Do(req)

			if err != nil {
				println("Request failed for", requestURL)
				siteTitle = domainName
				siteDescription = ipAddress
			} else {
				if response.StatusCode == 200 {
					document, err := goquery.NewDocumentFromReader(response.Body)
					if err != nil {
						println("Failed to parse", requestURL)
						siteTitle = domainName
						siteDescription = ipAddress
					} else {
						siteTitle = document.Find("title").Text()
						siteDescriptionExists := false
						siteDescription, siteDescriptionExists = document.Find("meta[name='description']").Attr("content")
						if !siteDescriptionExists {
							println("No site description")
							siteDescription = ipAddress
						}
					}

					println("Title", siteTitle)
					println("Description", siteDescription)
					println(line)
				}
			}
			descriptionClean := strings.ToValidUTF8(siteDescription, "")
			titleClean := strings.ToValidUTF8(siteTitle, "")
			dbEntry := db.Site{
				Domain:      domainName,
				IP:          ipAddress,
				DiscordID:   strings.Split(scanner.Text(), " ")[1],
				Title:       titleClean,
				Description: descriptionClean,
			}
			fullDatabase = append(fullDatabase, dbEntry)
		}

		line++
	}
	db.SiteDirectory.InsertMany(context.TODO(), fullDatabase)
}
