package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesChannels router
func AddRoutesChannels(r *gin.RouterGroup) {
	controller := controllers.ChannelController{Service: services.NewChannelService()}
	routes := r.Group("/channels")
	{
		routes.GET("", controller.GetAll)
		routes.POST("", controller.Create)
		routes.GET("/:id", controller.Get)
		routes.GET("/:id/followers", controller.GetFollowers)
		routes.DELETE("/:id", controller.Delete)
		routes.PUT("/:id", controller.Update)
	}
	// routes2 := r.Group("/users")
	// {
	// 	routes2.GET("/:id/channels", controller.GetJoinedGroup)
	// }

}
