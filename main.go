package main

import (
	"Ugle/api"
	"Ugle/db"
	//"Ugle/utils"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"net/http"
)

func ApiCache(url string, file string) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
		return
	}
	defer resp.Body.Close()
	out, err := os.Create(file)
	if err != nil {
		// panic?
	}
	defer out.Close()
	io.Copy(out, resp.Body)
}

func main() {
	ApiCache("https://raw.githubusercontent.com/ucanet/ucanet-registry/main/ucanet-registry.txt", "./registry/registry.txt")
	ApiCache("http://ucanet.net/sitelist.txt", "./registry/sitelist.txt")
	db.Init()
	
	//utils.CreateDatabaseFromRegistry()
	//WILL INDEX ALL SITES, ATTENTION!!!!! ^^^^^
	
	router := gin.Default()
	router.TrustedPlatform = gin.PlatformCloudflare
	router.GET("/search", api.MongoApiSearch)
	router.Use(static.Serve("/", static.LocalFile("./static", false)))
	router.Run("127.0.0.1:80")
}
