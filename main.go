package main

import (
	"Ugle/api"
	"Ugle/db"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.TrustedPlatform = gin.PlatformCloudflare

	router.GET("/api/v1/updateCache", api.ApiCache)
	router.GET("/api/v1/search", api.MongoApiSearch)
	router.Use(static.Serve("/", static.LocalFile("./static", false)))
	db.Init()

	//utils.CreateDatabaseFromRegistry()
	//WILL INDEX ALL SITES, ATTENTION!!!!! ^^^^^
	router.Run("localhost:80")
}
