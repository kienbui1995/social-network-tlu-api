package services

import (
	"errors"

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
}

// UserService struct
type userService struct{}

// NewUserService to constructor
func NewUserService() userService {
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
