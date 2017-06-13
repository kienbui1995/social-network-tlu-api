package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesSubjects router
func AddRoutesSubjects(r *gin.RouterGroup) {
	controller := controllers.SubjectController{Service: services.NewSubjectService()}
	routes := r.Group("/subjects")
	{
		routes.PUT("", controller.UpdateFromTLU)
		routes.GET("", controller.GetAll)
	}
}
