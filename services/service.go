package services

import (
	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
)

var conn, _ = neoism.Connect(configs.URLDB)
