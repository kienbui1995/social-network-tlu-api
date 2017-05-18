package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRouterSubscriptions router
func AddRouterSubscriptions(r *gin.RouterGroup) {
	controller := controllers.SubscriptionController{Service: services.NewSubscriberService()}
	routes := r.Group("/users")
	{
		routes.GET("/:id/followers", controller.GetFollowers)
		routes.GET("/:id/subscriptions", controller.GetSubcriptions)
		routes.POST("/:id/subscriptions", controller.CreateSubscription)   // follow user id
		routes.DELETE("/:id/subscriptions", controller.DeleteSubscription) // unfolow user id
	}

}
