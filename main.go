package main

import (
	"Ugle/api"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.TrustedPlatform = gin.PlatformCloudflare
	router.GET("/api/v1/search", api.ApiSearch)
	router.Use(static.Serve("/", static.LocalFile("./static", false)))
	router.Run("localhost:8080")
}
