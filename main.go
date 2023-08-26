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

func ApiCache() {
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

func main() {
	router := gin.Default()
	router.TrustedPlatform = gin.PlatformCloudflare

	router.GET("/search", api.MongoApiSearch)
	router.Use(static.Serve("/", static.LocalFile("./static", false)))
	
	ApiCache()
	db.Init()
	//utils.CreateDatabaseFromRegistry()
	//WILL INDEX ALL SITES, ATTENTION!!!!! ^^^^^
	
	router.Run("192.168.86.30:8080")
}
