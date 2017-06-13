package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesTeachers router
func AddRoutesTeachers(r *gin.RouterGroup) {
	controller := controllers.TeacherController{Service: services.NewTeacherService()}
	routes := r.Group("/teachers")
	{
		routes.PUT("", controller.UpdateFromTLU)
		routes.GET("", controller.GetAll)
	}
}
