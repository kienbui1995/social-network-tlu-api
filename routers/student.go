package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesStudents router
func AddRoutesStudents(r *gin.RouterGroup) {
	controller := controllers.StudentController{Service: services.NewStudentService()}
	routes := r.Group("/students")
	{
		routes.PUT("", controller.UpdateFromTLU)
		routes.GET("", controller.GetAll)
	}
}
