package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// UserServiceInterface include method list
type UserServiceInterface interface {
	GetAll(helpers.ParamsGetAll) (models.PublicUsers, error)
	Get(int64) (models.PublicUser, error)
	Delete(int64) (bool, error)
	Create(models.User) (int64, error)
	Update(userID int64, newUser models.InfoUser) (models.User, error)
	CheckExistUsername(string) (bool, error)
	CheckExistEmail(string) (bool, error)
	CreateEmailActive(string, string, int64) error
	CheckExistUser(int64) (bool, error)

	CreateRequestLinkCode(request models.RequestLinkCode, userID int64) (int64, error)
	AcceptLinkCode(requestID int64) (bool, error)
	AcceptLinkCodeByEmail(requestID int64, code string) (bool, error)
	DeleteRequestLinkCode(requestID int64) (bool, error)
	GetAllRequestsLinkCode(params helpers.ParamsGetAll) ([]models.RequestLinkCode, error)
	CheckExistRequestLinkCode(requestID int64) (bool, error)
}

// UserService struct
type userService struct{}

// NewUserService to constructor
func NewUserService() UserServiceInterface {
	return userService{}

}

// GetAll func
// helpers.ParamsGetAll
// models.PublicUsers error
func (service userService) GetAll(params helpers.ParamsGetAll) (models.PublicUsers, error) {
	stmt := `
	MATCH (u:User)
	return u {id:ID(u), .*}  AS user

		SKIP {skip}
		LIMIT {limit}
		`
	p := map[string]interface{}{
		"skip":  params.Skip,
		"limit": params.Limit,
	}
	var res []struct {
		User models.PublicUser `json:"user"`
	}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: p,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}
	//fmt.Printf("res: %v", res)
	var list models.PublicUsers
	for _, val := range res {
		list = append(list, val.User)
	}

	return list, nil
}

// Get func
// int64
// models.User error
func (service userService) Get(id int64) (models.PublicUser, error) {
	var user = models.PublicUser{}
	stmt := `
		MATCH (u:User)
		WHERE ID(u) = {id}
		RETURN
			u{id:ID(u), .*} as user
		`
	params := neoism.Props{"id": id}

	res := []struct {
		User models.PublicUser `json:"user"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)

	if err != nil {
		return user, err
	}
	if len(res) == 1 {

		return res[0].User, nil
	} else if len(res) > 1 {
		return user, errors.New("Many User")
	} else {
		return user, errors.New("No User")
	}
}

// Delete func
// int64
// bool error
func (service userService) Delete(id int64) (bool, error) {
	stmt := `
		MATCH (u:User) WHERE ID(u) = {id}
		DETACH DELETE u
		`
	params := neoism.Props{"id": id}

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

// Create func
// models.User
// int64 error
func (service userService) Create(user models.User) (int64, error) {
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
		large_avatar: {large_avatar},
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
		"large_avatar":  user.LargeAvatar,
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

// Update func
// int64 models.InfoUser
// models.User error
func (service userService) Update(userID int64, newUser models.InfoUser) (models.User, error) {
	stmt := `
	MATCH (u:User)
	WHERE ID(u) = {userid}
	SET u += {p}, u.updated_at = TIMESTAMP()
	RETURN properties(u) AS user
	`
	params := neoism.Props{
		"userid": userID,
		"p":      newUser,
	}

	res := []struct {
		User models.User `json:"user"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.User{}, err
	}

	if len(res) > 0 {
		if res[0].User.ID >= 0 {
			return res[0].User, nil
		}
	}
	return models.User{}, nil
}

// CheckExistUsername
// string
// bool error
func (service userService) CheckExistUsername(username string) (bool, error) {
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
// string
// bool error
func (service userService) CheckExistEmail(email string) (bool, error) {
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

// CreateEmailActive
// string string int64
// error
func (service userService) CreateEmailActive(email string, activecode string, userid int64) error {
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

// 	CheckExistUser
// int64
// bool error
func (service userService) CheckExistUser(id int64) (bool, error) {
	stmt := `
	MATCH (u:User)
	WHERE ID(u) = {id}
	RETURN ID(u) AS id
	`
	params := neoism.Props{
		"id": id,
	}

	var res []struct {
		ID int64 `json:"id"`
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
	if len(res) > 0 {
		if res[0].ID == id {
			return true, nil
		}
	}
	return false, nil
}

// 	CreateRequestLinkCode
// models.RequestLinkCode int64
// int64 error
func (service userService) CreateRequestLinkCode(request models.RequestLinkCode, userID int64) (int64, error) {
	var stmt string
	var params neoism.Props
	if len(request.Email) > 0 {
		stmt = `
			MATCH (u:User) WHERE ID(u) = {userID}
			MATCH (s:Student) WHERE toLower(s.email) = toLower({email})
			CREATE (u)-[f:IS_A]->(s)
			SET
				f.created_at = TIMESTAMP(),
				f.verifycation_code = {verifycationCode},
				f.verifycation_expired_at = TIMESTAMP()+1800000,
				f.status=0
			RETURN ID(f) AS id
	`
		params = neoism.Props{
			"userID":           userID,
			"email":            request.Email,
			"verifycationCode": request.VerificationCode,
		}
	} else {
		stmt = `
		MATCH (u:User) WHERE ID(u) = {userID}
		MATCH (s:Student) WHERE toLower(s.code) = toLower({code})
		CREATE (u)-[f:IS_A]->(s)
		SET
			f.created_at = TIMESTAMP(),
			f.full_name = {fullName},
			f.photo = {photo},
			f.status=0
		RETURN ID(f) AS id
	`
		params = neoism.Props{
			"code":     request.Code,
			"fullName": request.FullName,
			"photo":    request.Photo,
			"userID":   userID,
		}
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

// 	GetRequestsLinkCode
// helpers.ParamsGetAll
// []models.RequestLinkCode error
func (service userService) GetAllRequestsLinkCode(params helpers.ParamsGetAll) ([]models.RequestLinkCode, error) {
	stmt := fmt.Sprintf(`
		MATCH (u:User)-[f:IS_A{status:0}]->(s:Student)
		WHERE exists(f.photo)
		RETURN
			ID(f) AS id,
			u{id:ID(u),.username,.full_name,.avatar} AS user,
			s{id:ID(s),.code, name: s.first_name + " " + s.last_name} AS student,
			f.code AS code,
			f.full_name AS full_name,
			f.photo AS photo,
			f.created_at AS created_at
		ORDER BY %s
		SKIP {skip}
		LIMIT {limit}`, params.Sort)
	p := map[string]interface{}{
		"skip":  params.Skip,
		"limit": params.Limit,
	}
	res := []models.RequestLinkCode{}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: p,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res, nil
}

// 	AcceptLinkCode
// int64
// bool error
func (service userService) AcceptLinkCode(requestID int64) (bool, error) {
	stmt := `
	MATCH (u:User)-[f:IS_A{status:0}]->(s:Student)
	WHERE ID(f) = {requestID}
	REMOVE f.properities
	SET f.status = 1, f.updated_at = TIMESTAMP()
	`
	params := neoism.Props{
		"requestID": requestID,
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

// 	AcceptLinkCode
// int64
// bool error
func (service userService) AcceptLinkCodeByEmail(requestID int64, code string) (bool, error) {
	stmt := `
	MATCH (u:User)-[f:IS_A{status:0}]->(s:Student)
	WHERE ID(f) = {requestID} AND f.verifycation_code = {code} AND f.verifycation_expired_at > TIMESTAMP()
	REMOVE f.properities
	SET f.status = 1, f.updated_at = TIMESTAMP()
	RETURN CASE f.status WHEN 1 THEN true ELSE FALSE END AS accept
	`
	params := neoism.Props{
		"requestID": requestID,
		"code":      code,
	}
	var res []struct {
		Accept bool `json:"accept"`
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

	if len(res) > 0 {
		return res[0].Accept, nil
	}
	return false, nil
}

// 	DeleteRequestLinkCode
// int64
// bool error
func (service userService) DeleteRequestLinkCode(requestID int64) (bool, error) {
	stmt := `
		MATCH (u:User)-[f:IS_A{status:0}]->(s:Student) WHERE ID(f) = {requestID}
		DETACH DELETE f
		`
	params := neoism.Props{
		"requestID": requestID,
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

// 	CheckExistRequestLinkCode
// int64
// bool error
func (service userService) CheckExistRequestLinkCode(requestID int64) (bool, error) {
	stmt := `
	MATCH (u:User)-[f:IS_A{status:0}]->(s:Student)
	WHERE ID(f) = {requestID}
	RETURN ID(f) AS id
	`
	params := neoism.Props{
		"requestID": requestID,
	}

	var res []struct {
		ID int64 `json:"id"`
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
	if len(res) > 0 {
		if res[0].ID == requestID {
			return true, nil
		}
	}
	return false, nil
}

// //CreateRecoverPassword func
// func (service userService) CreateRecoverPassword(email string, recoveryCode string) (bool, error) {
// 	stmt := `
// 	MATCH (u:User)
// 	WHERE u.email = {email}
// 	SET u.recovery_code = {recoverycode}, u.recovery_expired_at = TIMESTAMP()+1800000
// 	RETURN ID(u) as id
// 	`
// 	params := neoism.Props{
// 		"email":        email,
// 		"recoverycode": recoveryCode,
// 	}
// 	var res []struct {
// 		ID int64 `json:"id"`
// 	}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	if err := conn.Cypher(&cq); err != nil {
// 		return false, err
// 	}
// 	if len(res) > 0 {
// 		if res[0].ID >= 0 {
// 			return true, nil
// 		}
// 		return false, errors.New("CreateRecoverPassword fail")
// 	}
// 	return false, nil
// }
//
// //VerifyRecoveryCode func
// func (service userService) VerifyRecoveryCode(email string, recoveryCode string) (int64, error) {
//
// 	stmt := `
// 		MATCH (u:User)
// 		WHERE u.email ={email} and u.recovery_code = {recoverycode}  and u.recovery_expired_at > TIMESTAMP()
// 		RETURN ID(u) as id
// 		`
// 	params := neoism.Props{
// 		"email":        email,
// 		"recoverycode": recoveryCode,
// 	}
// 	res := []struct {
// 		ID int64 `json:"id"`
// 	}{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	if err := conn.Cypher(&cq); err != nil {
// 		return -1, err
// 	}
// 	if len(res) > 0 {
// 		if res[0].ID >= 0 {
// 			return res[0].ID, nil
// 		}
// 		return -1, errors.New("VerifyRecoveryCode fail")
// 	}
// 	return -1, nil
// }
