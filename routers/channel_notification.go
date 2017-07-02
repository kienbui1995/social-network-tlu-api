package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesChannelNotifications router
func AddRoutesChannelNotifications(r *gin.RouterGroup) {
	controller := controllers.ChannelNotificationController{Service: services.NewChannelNotificationService()}
	routes := r.Group("/channel_notifications")
	{
		routes.GET("", controller.GetAll)        // get all channel notification of me
		routes.GET("/:id", controller.Get)       // get a noti of me
		routes.PUT("/:id", controller.Update)    // update a channel Notification
		routes.DELETE("/:id", controller.Delete) // Delete a channel notification
	}

	routes2 := r.Group("/channels")
	{
		routes2.GET("/:id/channel_notifications", controller.GetAllOfChannel)
		routes2.POST("/:id/channel_notifications", controller.Create)

	}

}
