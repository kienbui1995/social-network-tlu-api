package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// AccountServiceInterface include method list
type AccountServiceInterface interface {
	Login(models.Account) (int64, error)
	SaveToken(models.Account, string) (bool, error)
	CheckExistToken(int64, string) (bool, error)
	DeleteToken(int64, string) (bool, error)
	Create(models.User) (int64, error)

	GetDeviceByUserID(int64) ([]string, error)
	CheckExistUsername(string) (bool, error)
	CheckExistEmail(string) (bool, error)
	CreateEmailActive(string, string, int64) error

	CreateRecoverPassword(string, string) error
	VerifyRecoveryCode(string, string) (int64, error)
	AddUserRecoveryKey(int64, interface{}) error
	RenewPassword(int64, string, string) (bool, error)
	DeleteRecoveryProperty(int64) (bool, error)
	CheckExistFacebookID(string) (int64, error)
}

// accountService struct
type accountService struct{}

// NewAccountService contructor
func NewAccountService() accountService {
	return accountService{}
}

// Login func to user login system
// models.Account
// int error
func (service accountService) Login(account models.Account) (int64, error) {
	stmt := `
	MATCH (u:User) WHERE u.username	 = {username} return ID(u) as id, u.password as password
	`
	params := neoism.Props{"username": account.Username, "password": account.Password}

	res := []struct {
		ID       int64  `json:"id"`
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
func (service accountService) SaveToken(account models.Account, token string) (bool, error) {
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid}
		MERGE (u)-[:LOGGED_IN]->(d:Device {device:{device}}) SET d.token = {token}
	` // chua test
	params := neoism.Props{"userid": account.ID, "token": token, "device": account.Device}

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
	 MATCH (u:User) WHERE ID(u) = {accountid} return exists( (u)-[:LOGGED_IN]->(:Device{ token:{token} }) ) as exist_token
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
	fmt.Printf("res: %v\n", res)
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
func (service accountService) DeleteToken(accountid int64, token string) (bool, error) {
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
func (service accountService) GetDeviceByUserID(accountid int64) ([]string, error) {

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

// Create func
func (service accountService) Create(user models.User) (int64, error) {
	stmt := `
	Create (u:User{
		username: {username},
		password: {password},
		email: {email},
		first_name: {firstname},
		middle_name: {middlename},
		last_name: {lastname},
		full_name: {fullname},
		about: {about},
		gender: {gender},
		birthday: {birthday},
		avatar: {avatar},
		cover: {cover},
		status: {status},
		is_vertified: {isvertified},
		facebook_id: {facebookid},
		facebook_token: {facebooktoken},
		posts: {posts},
		followers: {followers},
		followings: {followings}
		}) SET u.created_at = TIMESTAMP()
	return ID(u) as id
	`
	params := neoism.Props{
		"username":      user.Username,
		"password":      user.Password,
		"email":         user.Email,
		"firstname":     user.FirstName,
		"middlename":    user.MiddleName,
		"lastname":      user.LastName,
		"fullname":      user.FullName,
		"about":         user.About,
		"gender":        user.Gender,
		"birthday":      user.Birthday,
		"avatar":        user.Avatar,
		"cover":         user.Cover,
		"status":        user.Status,
		"isvertified":   user.IsVertified,
		"facebookid":    user.FacebookID,
		"facebooktoken": user.FacebookToken,
		"posts":         user.Posts,
		"followers":     user.Followers,
		"followings":    user.Followings,
	}
	type resStruct struct {
		ID int64 `json:"id"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return -1, err

	}
	if len(res) > 0 {
		return res[0].ID, nil
	}
	return -1, nil
}

// CheckExistUsername
func (service accountService) CheckExistUsername(username string) (bool, error) {
	stmt := `
	MATCH (u:User)
	WHERE u.username = {username}
	RETURN ID(u) AS id
	`
	params := neoism.Props{
		"username": username,
	}
	type resStruct struct {
		ID int64 `json:"id"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	if len(res) > 0 && res[0].ID >= 0 {
		return true, nil
	}
	return false, nil
}

// CheckExistEmail func
func (service accountService) CheckExistEmail(email string) (bool, error) {
	stmt := `
	MATCH (u:User)
	WHERE u.email = {email}
	RETURN ID(u) AS id
	`
	params := neoism.Props{
		"email": email,
	}
	type resStruct struct {
		ID int64 `json:"id"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	if len(res) > 0 && res[0].ID >= 0 {
		return true, nil
	}
	return false, nil
}

//
func (service accountService) CreateEmailActive(email string, activecode string, userid int64) error {
	stmt := `
	MATCH (u:User)
	WHERE ID(u) = {userid} AND u.email = {email}
	SET u.active_code = {activecode}
	`
	params := neoism.Props{
		"userid":     userid,
		"email":      email,
		"activecode": activecode,
	}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return err
	}
	return nil

}

//CreateRecoverPassword func
func (service accountService) CreateRecoverPassword(email string, recoverycode string) error {
	stmt := `
	MATCH (u:User)
	WHERE u.email = {email}
	SET u.recovery_code = {recoverycode}, u.recovery_expired_at = TIMESTAMP()+1800000
	`
	params := neoism.Props{
		"email":        email,
		"recoverycode": recoverycode,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
	}
	if err := conn.Cypher(&cq); err != nil {
		return err
	}
	return nil

}

//VerifyRecoveryCode func
func (service accountService) VerifyRecoveryCode(email string, recoverycode string) (int64, error) {

	stmt := `
		MATCH (u:User)
		WHERE u.email ={email} and u.recovery_code = {recoverycode}  and u.recovery_expired_at > TIMESTAMP()
		RETURN ID(u) as id
		`
	params := neoism.Props{
		"email":        email,
		"recoverycode": recoverycode,
	}
	res := []struct {
		ID int64 `json:"id"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	if err := conn.Cypher(&cq); err != nil {
		return -1, err
	}
	if len(res) > 0 {
		return res[0].ID, nil
	}
	return -1, nil

}

//AddUserRecoveryKey func to add a property of user
func (service accountService) AddUserRecoveryKey(userid int64, value interface{}) error {
	stmt := `
	MATCH(u:User) WHERE ID(u)= {userid} SET u.recovery_key = {value}
	`
	params := neoism.Props{
		"userid": userid,
		"value":  value,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
	}
	err := conn.Cypher(&cq)
	return err
}

//RenewPassword func
func (service accountService) RenewPassword(userid int64, recoverykey string, newpassword string) (bool, error) {
	stmt := `
	MATCH(u:User) WHERE ID(u)= {userid} AND u.recovery_key = {recoverykey} SET u.password = {newpassword}
	RETURN u.password as password
	`
	res := []struct {
		Password string `json:"password"`
	}{}
	params := neoism.Props{
		"userid":      userid,
		"recoverykey": recoverykey,
		"newpassword": newpassword,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	if len(res) > 0 && res[0].Password == newpassword {
		return true, nil
	}
	return false, nil
}

//DeleteRecoveryProperty func
func (service accountService) DeleteRecoveryProperty(userid int64) (bool, error) {
	stmt := fmt.Sprintf("MATCH (u:User) WHERE ID(u) = %d REMOVE u.recovery_key, u.recovery_code, u.recovery_expired_at RETURN ID(u) as id, exists(u.recovery_key) as k, exists(u.recovery_code) as c,exists(u.recovery_expired_at) as e ", userid)
	type resStruct struct {
		ID int64 `json:"id"`
		K  bool  `json:"k"`
		C  bool  `json:"c"`
		E  bool  `json:"e"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement: stmt,

		Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	if len(res) != 0 {
		if res[0].ID == userid && res[0].K == true && res[0].C == true && res[0].E == true {
			return true, nil
		}
		return false, nil
	}
	return false, errors.New("No exist user")
}

//CheckExistFacebookID func
func (service accountService) CheckExistFacebookID(facebookid string) (int64, error) {
	stmt := `
	MATCH (u:User)
	WHERE u.facebook_id = {facebookid}
	RETURN ID(u) AS id
	`
	params := neoism.Props{
		"facebookid": facebookid,
	}
	type resStruct struct {
		ID int64 `json:"id"`
	}
	res := []resStruct{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return -1, err
	}

	if len(res) > 0 {
		return res[0].ID, nil
	}
	return -1, nil
}
