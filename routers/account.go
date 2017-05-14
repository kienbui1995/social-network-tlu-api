package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/controllers"
	"github.com/kienbui1995/social-network-tlu-api/services"
)

// AddRoutesAccounts router
func AddRoutesAccounts(r *gin.RouterGroup) {
	controller := controllers.AccountController{Service: services.NewAccountService()}
	routes := r.Group("")
	{
		routes.POST("/login", controller.Login)
		routes.POST("/logout", controller.Logout)

	}
}
