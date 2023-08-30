package utils

import (
	"Ugle/db"
	"bufio"
	"context"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"net/url"
	"net"
	"os"
	"time"
	"strings"
)

type list_entry struct {
	tag_list string
	specified_description string
}

func CreateDatabaseFromRegistry() {
	f, err := os.Open("./registry/sitelist.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Splits on newlines by default.
	var site_list = make(map[string]list_entry)
	var category_list = make(map[string]string)
	current_category := ""
	scanner := bufio.NewScanner(f)
	line := 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	for scanner.Scan() {
		first_item := strings.SplitN(scanner.Text(), " ", 3)[0]
		second_item := strings.SplitN(scanner.Text(), " ", 3)[1]
		third_item := strings.SplitN(scanner.Text(), " ", 3)[2]
		
		if first_item[0] == '!' {
			current_category = second_item + " " + third_item
		} else if first_item[0] == '-' {
			// skip
		} else if first_item[0] == '#' {
			category_list[third_item] = current_category + " " + second_item
		} else {
			current_tags := ""
			for _, letter := range second_item {
				current_tags += category_list[string(letter)] + " "
			}
			site_list[first_item] = list_entry{
				tag_list: strings.ReplaceAll(strings.ReplaceAll(current_tags, "_", " "), "-", " "),
				specified_description: third_item,
			}
		}

		line++
	}
	
	f, err = os.Open("./registry/registry.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Splits on newlines by default.
	scanner = bufio.NewScanner(f)
	line = 1
	// https://golang.org/pkg/bufio/#Scanner.Scan
	fullDatabase := []interface{}{}
	for scanner.Scan() {
		domainName := strings.Split(scanner.Text(), " ")[0]
		ipAddress := strings.Split(scanner.Text(), " ")[2]
		
		if ipAddress != "0.0.0.0" {
			var requestURL = "http://" + domainName + "/"
			var proxyURL = "http://"
			if net.ParseIP(ipAddress).To4() == nil {
				proxyURL += "135.148.41.26:80"
			} else {
				proxyURL += ipAddress + ":80"
			}
			proxyObj, perr := url.Parse(proxyURL)
			if perr != nil {
				panic(perr)
			}
			var siteTitle string
			var siteDescription string
			specifiedEntry, entryExists := site_list[domainName]
			if !entryExists {
				specifiedEntry, entryExists = site_list["www." + domainName]
			}
			req, qerr := http.NewRequest("GET", requestURL, nil)
			if qerr != nil {
				panic(qerr)
			}
			req.Host = domainName
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyObj),
			}
			client := &http.Client{
				Transport: transport,
				Timeout: 4 * time.Second,
			}
			response, err := client.Do(req)

			if err != nil {
				println("Request failed for", requestURL)
				siteTitle = domainName
			} else {
				if response.StatusCode == 200 {
					document, err := goquery.NewDocumentFromReader(response.Body)
					if err != nil {
						println("Failed to parse", requestURL)
						siteTitle = domainName
					} else {
						siteTitle = document.Find("title").Text()
						siteDescriptionExists := false
						siteDescription, siteDescriptionExists = document.Find("meta[name='description']").Attr("content")
						if !siteDescriptionExists {
							println("No site description")
						}
					}

					println("Title", siteTitle)
					println("Description", siteDescription)
					println(line)
				}
			}
			descriptionClean := strings.ToValidUTF8(siteDescription, "")
			titleClean := strings.ToValidUTF8(siteTitle, "")
			specified_description := ""
			tag_list := ""
			if entryExists {
				specified_description = specifiedEntry.specified_description
				tag_list = specifiedEntry.tag_list
			}
			if titleClean == "" {
				titleClean = strings.ToValidUTF8(domainName, "")
			}
			if descriptionClean == "" {
				if entryExists {
					descriptionClean = strings.ToValidUTF8(strings.ReplaceAll(specifiedEntry.specified_description, "[Protoweb]", ""), "")
				} else {
					descriptionClean = "No site description"
				}
			}
			if domainName == "ugle.com" {
				descriptionClean = "but why did you search for Ugle... on Ugle?"
			}
			dbEntry := db.Site{
				Domain:      domainName,
				IP:          ipAddress,
				DiscordID:   strings.Split(scanner.Text(), " ")[1],
				Title:       titleClean,
				Description: descriptionClean,
				SpecifiedDescription: specified_description,
				SpecifiedTags: tag_list,
			}
			fullDatabase = append(fullDatabase, dbEntry)
		}

		line++
	}
	db.SiteDirectory.Drop(context.TODO())
	db.SiteDirectory.InsertMany(context.TODO(), fullDatabase)
}
