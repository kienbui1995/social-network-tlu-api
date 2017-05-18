package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesPosts router
func AddRoutesPosts(r *gin.RouterGroup) {
	controller := controllers.PostController{Service: services.NewPostService()}
	routes := r.Group("/posts")
	{

		routes.GET("/:id", controller.Get)
		routes.DELETE("/:id", controller.Delete)
		// routes.POST("", controller.Create)
		routes.PUT("/:id", controller.Update)
		routes.POST("/:id/likes", controller.CreateLike)
		routes.DELETE("/:id/likes", controller.DeleteLike)
		routes.GET("/:id/likes", controller.GetLikes)

	}
	routes2 := r.Group("/users")
	{
		routes2.POST("/:id/posts", controller.Create)
		routes2.GET("/:id/posts", controller.GetAll)
	}
}
