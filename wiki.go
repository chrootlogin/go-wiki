package main

import (
	"os"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/chrootlogin/go-wiki/src/page"
	"github.com/chrootlogin/go-wiki/src/frontend"
)

var port = "8000"

func main() {
	if len(os.Getenv("PORT")) > 0 {
		port = os.Getenv("PORT")
	}

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
		c.Redirect(http.StatusMovedPermanently, "/wiki")
	})

	router.GET("/api/page/*page", page.GetPageHandler)
	router.GET("/wiki/*path", frontend.GetFrontendHandler)

	router.Run(":" + port)
}
