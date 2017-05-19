package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRouterHome router
func AddRouterHome(r *gin.RouterGroup) {
	controller := controllers.HomeController{Service: services.NewHomeService()}
	routes := r.Group("")
	{
		routes.GET("/news_feed", controller.NewsFeed)
		routes.GET("/find_user", controller.FindUser)
	}
}
