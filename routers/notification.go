package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesNotifications router
func AddRoutesNotifications(r *gin.RouterGroup) {
	controller := controllers.NotificationController{Service: services.NewNotificationService()}
	routes := r.Group("/notifications")
	{
		routes.GET("", controller.GetAll) // get all noti of me
		// routes.GET("/:id", controller)    // get a noti of me
		// routes.PUT("/:id", controlle)     // update Seen Notification
	}

}
