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
		routes.POST("/sign_up", controller.SignUp)
		routes.POST("/logout", controller.Logout)
		routes.POST("/login_facebook", controller.LoginViaFacebook)
		routes.POST("/forgot_password", controller.ForgotPassword)
		routes.POST("/verify_recovery_code", controller.VerifyRecoveryCode)
		routes.PUT("/renew_password", controller.RenewPassword)

	}
}
