package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesUsers router
func AddRoutesUsers(r *gin.RouterGroup) {
	controller := controllers.UserController{Service: services.NewUserService()}
	routes := r.Group("/users")
	{
		routes.GET("", controller.GetAll)
		routes.POST("", controller.Create)
		routes.GET("/:id", controller.Get)
		routes.DELETE("/:id", controller.Delete)
		routes.PUT("/:id", controller.Update)

		// work with studentCode
		routes.POST("/:id/requests_link_code", controller.RequestLinkCode)
	}
	routes2 := r.Group("/requests_link_code")
	{
		routes2.GET("", controller.GetAllRequestLinkCode)
		routes2.PUT("/:id", controller.AcceptLinkCode)
		routes2.DELETE("/:id", controller.DeleteRequestLinkCode)
	}
}
