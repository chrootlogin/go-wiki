package main

import (
	"os"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"github.com/chrootlogin/go-wiki/src/page"
	"github.com/chrootlogin/go-wiki/src/user"
	"github.com/chrootlogin/go-wiki/src/auth"
	"github.com/chrootlogin/go-wiki/src/lib/common"
)

var port = ""

func main() {
	if len(os.Getenv("PORT")) == 0 {
		port = "8000"
	} else {
		port = os.Getenv("PORT")
	}

	log.Println("Starting go-wiki.")
	common.LoadPlugins()
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

	//common.GetPluginRegistry().RunEngine(router)

	router.Run(":" + port)
}
