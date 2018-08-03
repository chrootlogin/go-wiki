package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/chrootlogin/go-wiki/src/page"
)

func main() {
	log.Println("Starting go-wiki.")
	initRouter()
	log.Println("go-wiki is running.")
}

func initRouter() {
	router := gin.Default()

	// Allow cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization", "*")
	corsConfig.AddAllowMethods("HEAD", "GET", "POST", "PUT", "DELETE")
	router.Use(cors.New(corsConfig))

	// load authentication
	//am := auth.GetAuthMiddleware()

	// public routes
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/webapp")
	})

	router.GET("/api/page/*page", page.GetPageHandler)


	// add webapp frontend
	//router.Static("/webapp", "./frontend")
	router.Run(":8000")
}
