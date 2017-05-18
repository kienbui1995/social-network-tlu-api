package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesComments router
func AddRoutesComments(r *gin.RouterGroup) {
	controller := controllers.CommentController{Service: services.NewCommentService()}
	routes := r.Group("/comments")
	{

		routes.GET("/:id", controller.Get)
		routes.PUT("/:id", controller.Update)
		routes.DELETE("/:id", controller.Delete)

	}
	routes2 := r.Group("/posts")
	{
		routes2.POST("/:id/comments", controller.Create)
		routes2.GET("/:id/comments", controller.GetAll)
	}
}
