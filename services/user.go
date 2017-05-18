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
	Get(int64) (models.User, error)
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
func (service userService) Get(id int64) (models.User, error) {
	var user = models.User{}
	stmt := `
		MATCH (u:User)
		WHERE ID(u) = {id}
		RETURN
		u.avatar as avatar, u.about as about, u.birthday as birthday, u.gender as gender, u.cover as cover,
		ID(u) as id,
		u.username as username,
		u.full_name as full_name, u.first_name as first_name, u.last_name as last_name,
		u.email as email, u.status as Status,
		u.followers as followers, u.followings as followings, u.posts as posts,
		u.created_at as created_at, u.updated_at as updated_at
		LIMIT 25;
		`
	params := neoism.Props{"id": id}

	res := []models.User{}
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

		return res[0], nil
	} else if len(res) > 1 {
		return user, errors.New("Many User")
	} else {
		return user, errors.New("No User")
	}
}

// Delete func
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

// Update func
func (service userService) Update(userID int64, newUser models.InfoUser) (models.User, error) {
	stmt := `
	MATCH (u:User)
	WHERE ID(u) = {userid}
	SET u += {p}
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
		return models.User{}, nil
	}
	return models.User{}, nil
}

// CheckExistUsername
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

// 	CheckExistUser(int64)
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
		return false, nil
	}
	return false, errors.New("CheckExistUser fail")
}
