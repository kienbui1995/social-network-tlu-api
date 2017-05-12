package services

import (
	"errors"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// AccountServiceInterface include method list
type AccountServiceInterface interface {
	Authentication(username string, password string, device string) models.User
}

type accountService struct{}

// NewAccountService contructor
func NewAccountService() *accountService {
	return new(accountService)
}

// Authentication func
func (service accountService) Authentication(username string, password string, device string) models.User {
	return models.User{}
}

// Login func to user login system
// models.Account
// int error
func (service accountService) Login(account models.Account) (int, error) {
	stmt := `
	MATCH (u:User) WHERE u.username	 = {username} return ID(u) as id, u.password as password
	`
	params := neoism.Props{"username": account.Username, "password": account.Password}

	res := []struct {
		ID       int    `json:"id"`
		Password string `json:"password"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return -1, err
	}

	if len(res) == 0 {
		return -1, errors.New("No exist user")
	} else if res[0].Password == account.Password {
		return res[0].ID, nil
	}
	return res[0].ID, errors.New("Wrong password")
}

// SaveToken func to insert token to db
// int string string
// bool error
func (service accountService) SaveToken(account models.Account, tokenstring string) (bool, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid}
		MERGE (u)-[:LOGGED_IN]->(d:Device {device:{device}}) SET d.token = {token}
	` // chua test
	params := neoism.Props{"userid": account.ID, "token": tokenstring, "device": account.Device}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	return true, nil
}

// CheckExistToken func to check exist token in DB
// int string
// bool error
func (service accountService) CheckExistToken(accountid int64, token string) (bool, error) {
	//check exist token
	stmt := `
	 MATCH (u:User) WHERE ID(u) = {accountid} return exists( ((u)-[:LOGGED_IN]->(:Device{ token:{token}})) ) as exist_token
	`
	params := neoism.Props{"accountid": accountid, "token": token}

	res := []struct {
		ExistToken bool `json:"exist_token"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	if len(res) == 0 {
		return false, errors.New("Token don't exist")
	}
	if res[0].ExistToken != true {
		return false, errors.New("Wrong token")
	}
	return true, nil
}

// DeleteToken func to delete token of user
// int string
// bool error
func (service accountService) DeleteToken(accountid int, token string) (bool, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {accountid}
	MATCH ((u)-[:LOGGED_IN]->(d))
	WHERE d.token = {token}
	DETACH DELETE d
	`
	params := neoism.Props{
		"accountid": accountid,
		"token":     token,
	}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetDeviceByUserID func to get deive list
// int
// []string error
func (service accountService) GetDeviceByUserID(accountid int) ([]string, error) {

	stmt := `
		MATCH (u:User)-[:LOGGED_IN]->(d:Device)
		WHERE ID(u) = {accountid}
		RETURN d.device AS device
			`
	params := map[string]interface{}{"accountid": accountid}
	res := []struct {
		Device string `json:"device"`
	}{}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		var devices []string
		for index := 0; index < len(res); index++ {
			devices = append(devices, res[index].Device)
		}
		return devices, nil
	}
	return nil, nil
}
