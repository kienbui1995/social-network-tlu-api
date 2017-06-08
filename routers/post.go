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

		routes.POST("/:id/subcriptions", controller.CreateFollow)
		routes.DELETE("/:id/subcriptions", controller.DeleteFollow)

		routes.POST("/:id/reports", controller.CreateReport)
		routes.DELETE("/:id/reports", controller.DeleteReport)

		routes.GET("/:id/users", controller.GetUsers) // get users who can_mentioned/mentioned or liked/commented post

	}
	// post with user
	routesUser := r.Group("/users")
	{
		routesUser.POST("/:id/posts", controller.Create)
		routesUser.GET("/:id/posts", controller.GetAll)
	}

	// post in group
	routesGroup := r.Group("/groups")
	{
		routesGroup.POST("/:id/posts", controller.CreateGroupPost)
		routesGroup.GET("/:id/posts", controller.GetAllGroupPost)

	}
}
