package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesSemester router
func AddRoutesSemester(r *gin.RouterGroup) {
	controller := controllers.SemesterController{Service: services.NewSemesterService()}
	routes := r.Group("/semesters")
	{
		routes.PUT("", controller.UpdateFromTLU)
		routes.GET("", controller.GetAll)
	}
}
