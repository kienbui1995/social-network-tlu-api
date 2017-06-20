package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesExamSchedules router
func AddRoutesExamSchedules(r *gin.RouterGroup) {
	controller := controllers.ExamScheduleController{Service: services.NewExamScheduleService()}
	routes := r.Group("/exam_schedules")
	{
		routes.PUT("", controller.UpdateFromTLU)
		routes.GET("", controller.GetAllByStudent)
	}
}
