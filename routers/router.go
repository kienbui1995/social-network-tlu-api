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
		//
		AddRoutesAccounts(v1)

		v1.Use(middlewares.AuthRequired())
		{

			AddRoutesUsers(v1)
			AddRouterSubscriptions(v1)
			AddRoutesPosts(v1)
			// AddRoutesGroups(v1)
		}
	}

	return router
}
