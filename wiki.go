package main

import (
	"os"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/chrootlogin/go-wiki/src/page"
	"github.com/chrootlogin/go-wiki/src/frontend"
	"github.com/chrootlogin/go-wiki/src/user"
	"github.com/chrootlogin/go-wiki/src/auth"
)

var port = ""

func main() {
	if len(os.Getenv("PORT")) == 0 {
		port = "8000"
	} else {
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

	// public routes
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/wiki")
	})
	router.GET("/wiki/*path", frontend.GetFrontendHandler)

	// authentication
	am := auth.GetAuthMiddleware()
	router.POST("/user/login", am.LoginHandler)
	router.POST("/user/register", user.RegisterHandler)

	// API
	api := router.Group("/api/")
	api.Use(am.MiddlewareFunc())
	{
		api.GET("/page/*path", page.GetPageHandler)
		api.POST("/page/*path", page.PostPageHandler)
		api.PUT("/page/*path", page.PutPageHandler)

		api.POST("/preview", page.PostPreviewHandler)
	}

	router.Run(":" + port)
}
