package handlers

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// ProfileHandler to handled user's profile page
func ProfileHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	userImg := session.Get("user-img")
	c.HTML(http.StatusOK, "user", gin.H{"user": userID, "userImg": userImg})
}
