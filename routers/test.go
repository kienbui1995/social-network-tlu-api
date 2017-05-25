package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// AddRouterTest router
func AddRouterTest(r *gin.RouterGroup) {

	routes := r.Group("/test")
	{
		routes.GET("/user", func(c *gin.Context) {
			c.JSON(200, models.User{})
		})
		routes.GET("/post", func(c *gin.Context) {
			c.JSON(200, models.Post{})
		})
		routes.GET("/comment", func(c *gin.Context) {
			c.JSON(200, models.Comment{})
		})
	}
}
