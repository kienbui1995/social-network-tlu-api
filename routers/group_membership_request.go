package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesGroupMembershipRequests router
func AddRoutesGroupMembershipRequests(r *gin.RouterGroup) {
	controller := controllers.GroupMembershipRequestController{Service: services.NewGroupMembershipRequestService()}
	routes := r.Group("/group_membership_requests")
	{

		routes.GET("/:id", controller.Get)
		routes.DELETE("/:id", controller.Delete)
		routes.PUT("/:id", controller.Update)

	}
	routes2 := r.Group("/groups")
	{
		routes2.GET("/:id/requests", controller.GetAll)
		routes2.POST("/:id/requests", controller.Create)
	}

}
