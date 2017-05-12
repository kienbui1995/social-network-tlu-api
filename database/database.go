package database

import (
	"fmt"

	"github.com/jmcvetta/neoism"
)

// Neoism struct to connect db
type Neoism struct {
	IP       string
	User     string
	Password string
	Port     int
	Type     string
}

// CreateConnect func to  init connect to db
func (neo Neoism) CreateConnect() (*neoism.Database, bool) {
	url := fmt.Sprintf("%s://%s:%s@%s:%d/db/data/", neo.Type, neo.User, neo.Password, neo.IP, neo.Port)
	db, err := neoism.Connect(url)
	if err != nil {
		return nil, true
	}
	return db, false
}
