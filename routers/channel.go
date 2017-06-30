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
		routes.DELETE("/:id", controller.Delete)
		routes.PUT("/:id", controller.Update)

		// work for follow
		routes.GET("/:id/followers", controller.GetFollowers)
		routes.POST("/:id/followers", controller.CreateFollower)
		routes.DELETE("/:id/followers", controller.DeleteFollower)
	}
	routes2 := r.Group("/users")
	{
		routes2.GET("/:id/channels", controller.GetFollowedChannels)
	}

}
