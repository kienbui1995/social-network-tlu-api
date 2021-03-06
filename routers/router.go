package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/middlewares"
)

// InitRoutes to start router
func InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORS())

	v1 := router.Group("")
	{
		//test
		AddRouterTest(v1)

		AddRoutesAccounts(v1)

		v1.Use(middlewares.AuthRequired())
		{
			AddRouterHome(v1)
			AddRoutesUsers(v1)
			AddRouterSubscriptions(v1)
			AddRoutesPosts(v1)
			AddRoutesComments(v1)
			AddRoutesGroups(v1)
			AddRoutesChannels(v1)
			AddRoutesGroupMemberships(v1)
			AddRoutesNotifications(v1)
			AddRoutesSemesters(v1)
			AddRoutesSubjects(v1)
			AddRoutesTeachers(v1)
			AddRoutesClasses(v1)
			AddRoutesStudents(v1)
			AddRoutesExamSchedules(v1)
			AddRoutesChannelNotifications(v1)
			AddRoutesViolations(v1)
		}

	}
	return router
}
