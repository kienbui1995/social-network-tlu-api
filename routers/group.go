package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesGroups router
func AddRoutesGroups(r *gin.RouterGroup) {
	controller := controllers.GroupController{Service: services.NewGroupService()}
	routes := r.Group("/groups")
	{
		routes.GET("", controller.GetAll)
		routes.POST("", controller.Create)
		routes.GET("/:id", controller.Get)
		routes.DELETE("/:id", controller.Delete)
		routes.PUT("/:id", controller.Update)

		routes.GET("/:id/members", controller.GetMembers)
		routes.POST("/:id/members", controller.CreateMember)

		routes.GET("/:id/requests", controller.GetRequests)
		routes.POST("/:id/requests", controller.CreateRequest)

		routes.GET("/:id/reports", controller.GetReports)
		routes.POST("/:id/reports", controller.CreateReport)

		routes.GET("/:id/posts", controller.GetPosts)
		routes.POST("/:id/posts", controller.CreatePost)

	}

}
