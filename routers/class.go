package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesClasses router
func AddRoutesClasses(r *gin.RouterGroup) {
	controller := controllers.ClassController{Service: services.NewClassService()}
	routes := r.Group("/classes")
	{
		routes.PUT("", controller.UpdateFromTLU)
		routes.GET("", controller.GetAllByStudent)
	}
	routes2 := r.Group("/rooms")
	{
		routes2.GET("/:id/classes", controller.GetAllByRoom)
	}
	routes3 := r.Group("/teachers/")
	{
		routes3.GET("/:id/classes", controller.GetAllByTeacher)
	}
}
