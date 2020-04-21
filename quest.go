package main

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jogit82/google-oauth-gin/handlers"
	"github.com/jogit82/google-oauth-gin/middleware"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	token, err := handlers.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}
	store := sessions.NewCookieStore([]byte(token))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("goquestsession", store))
	router.Static("/css", "./static/css")
	router.Static("/img", "./static/img")
	router.LoadHTMLGlob("templates/*")

	// these pages have to have an authorized login
	authorized := router.Group("/auth")
	// per group middleware! in this case we use the custom created
	// AuthorizeRequest() middleware just in the "authorized" group.
	authorized.Use(middleware.AuthorizeRequest())
	{
		authorized.GET("/field", handlers.FieldHandler)
		authorized.GET("/user", handlers.ProfileHandler)
	}

	// these pages are public, no login needed
	router.GET("/", handlers.IndexHandler)
	router.GET("/login", handlers.LoginHandler)
	router.GET("/auth", handlers.AuthHandler)
	return router
}

func main() {
	router := setupRouter()
	if err := router.Run("127.0.0.1:9090"); err != nil {
		log.Fatal(err)
	}
}
