package handlers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// IndexHandler handles the location /.
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "Index title!",
		"add": func(a int, b int) int {
			return a + b
		},
	})
}

// PageHandler handles the location /.
func PageHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "page.html", gin.H{
		"title": "Page file title!!",
	})
}

// FieldHandler is a rudementary handler for logged in users.
func FieldHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.HTML(http.StatusOK, "field", gin.H{"user": userID})
}
