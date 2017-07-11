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
	}
	routes2 := r.Group("/users")
	{
		routes2.GET("/:id/groups", controller.GetJoinedGroup)
	}
	routes3 := r.Group("/students")
	{
		routes3.GET("/:id/groups", controller.GetClassGroupOfStudent)
	}
}
