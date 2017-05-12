package main

import (
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/routers"
)

func main() {
	routers.InitRoutes().Run(":" + configs.APIPort)
}
