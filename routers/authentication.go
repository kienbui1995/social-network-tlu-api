package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/kienbui1995/social-network-tlu-api/middlewares"
)

func AddRoutesAuthentication(r *gin.RouterGroup) {
	r.POST("/auth/auth_token", middlewares.SetToken)
}
