package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesGroupMemberships router
func AddRoutesGroupMemberships(r *gin.RouterGroup) {
	controller := controllers.GroupMembershipController{Service: services.NewGroupMembershipService()}
	routes2 := r.Group("/groups")
	{
		routes2.GET("/:id/members", controller.GetAll)          // get memberships (members)
		routes2.POST("/:id/members", controller.Create)         // join public group/request when privacy is private
		routes2.DELETE("/:id/members", controller.DeleteByUser) // out group

	}

	routes := r.Group("/group_memberships")
	{
		routes.PUT("/:id", controller.Update)    // make/remove admin/block member by admin
		routes.DELETE("/:id", controller.Delete) // user removed by admin

	}

}
