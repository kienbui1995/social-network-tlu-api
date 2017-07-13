package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesSemesters router
func AddRoutesSemesters(r *gin.RouterGroup) {
	controller := controllers.SemesterController{Service: services.NewSemesterService()}
	routes := r.Group("/semesters")
	{
		routes.PUT("", controller.UpdateFromTLU)
		routes.GET("", controller.GetAll)
	}
	routes2 := r.Group("/students")
	{

		routes2.GET("/:id/semesters", controller.GetSemesterOfStudent)
	}
	routes3 := r.Group("/teachers")
	{

		routes3.GET("/:id/semesters", controller.GetSemesterOfTeacher)
	}
}
