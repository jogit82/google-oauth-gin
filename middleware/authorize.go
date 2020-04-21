package middleware

import (
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

// AuthorizeRequest is used to authorize a request for a certain end-point group.
func AuthorizeRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user-id")
		if v == nil {
			c.HTML(http.StatusUnauthorized, "unauthorized.tmpl", gin.H{"message": "You need to login in order to view this page"})
			c.Abort() // c.index = abortIndex
		}
		c.Next() // c.index++
	}
}

/*
Dev Comments:
* http.StatusUnauthorized? List of Status Code here: https://golang.org/pkg/net/http/
* What can you do with Go Context? https://github.com/gin-gonic/gin/blob/master/context.go
* func (c *Context) HTML(code int, name string, obj interface{})
* renders the HTTP template specified by its file name.
* gin.H() is a shortcut for map[string]interface{}, you can nest this!!
* gin.H{
	"status":  gin.H{
	"code": http.StatusOK,
	"status": "Login Successful",
	},
	"message": message
	}
*/
