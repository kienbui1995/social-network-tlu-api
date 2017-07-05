package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesViolations router
func AddRoutesViolations(r *gin.RouterGroup) {
	controller := controllers.ViolationController{Service: services.NewViolationService()}
	routes := r.Group("/violations")
	{
		routes.GET("", controller.GetAll)        // get all channel notification of me
		routes.GET("/:id", controller.Get)       // get a noti of me
		routes.PUT("/:id", controller.Update)    // update a channel Notification
		routes.DELETE("/:id", controller.Delete) // Delete a channel notification
	}

	routes2 := r.Group("/students")
	{
		routes2.GET("/:id/violations", controller.GetAllOfStudent)

	}

	routes3 := r.Group("/supervisiors")
	{
		routes3.GET("/:id/violations", controller.GetAllOfSupervisior)

	}
}
