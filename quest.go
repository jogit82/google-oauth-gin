/*
 * Copyright 2018 Foolin.  All rights reserved.
 *
 * Use of this source code is governed by a MIT style
 * license that can be found in the LICENSE file.
 *
 */

package main

import (
	"log"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jogit82/google-oauth-gin/handlers"
	"github.com/jogit82/google-oauth-gin/middleware"
)

func setupRouter() *gin.Engine {
	router := gin.Default()
	router.HTMLRender = ginview.Default()
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
	// router.LoadHTMLGlob("templates/*")

	// these pages have to have an authorized login
	authorized := router.Group("")
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
	router.POST("/search", handlers.SearchHandler)
	router.GET("/random", handlers.RandomHandler)
	// router.GET("/moviesearchform", handlers.MovieSearchFormHandler)
	router.POST("/searchmovies", handlers.SearchMoviesHandler)
	router.GET("/page", handlers.PageHandler)
	return router
}

func main() {
	router := setupRouter()
	if err := router.Run("127.0.0.1:9090"); err != nil {
		log.Fatal(err)
	}
}
