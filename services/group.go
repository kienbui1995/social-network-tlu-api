package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// GroupServiceInterface include method list
type GroupServiceInterface interface {
	GetAll(params helpers.ParamsGetAll, myUserID int64) ([]models.GroupJoin, error)
	Get(groupID int64, myUserID int64) (models.GroupJoin, error)
	Delete(groupID int64) (bool, error)
	Create(group models.Group, myUserID int64) (int64, error)
	Update(groupID int64, newGroup models.InfoGroup) (models.Group, error)
	CheckExistGroup(groupID int64) (bool, error)
	GetMembers(params helpers.ParamsGetAll, groupID int64) ([]models.UserJoinedObject, error)
	CreateMember(groupID int64, userID int64) (bool, error)
	CreateRequest(groupID int64, userID int64) (bool, error)
	IncreasePosts(groupID int64) (bool, error)
	DecreasePosts(groupID int64) (bool, error)
}

// groupService struct
type groupService struct{}

// NewGroupService to constructor
func NewGroupService() groupService {
	return groupService{}
}

// GetAll func
// helpers.ParamsGetAll int64
// []models.Group error
func (service groupService) GetAll(params helpers.ParamsGetAll, myUserID int64) ([]models.GroupJoin, error) {
	var stmt string
	stmt = fmt.Sprintf(`
				MATCH(me:User) WHERE ID(me) = {myuserid}
				MATCH (g:Group)
				RETURN
					ID(g) AS id,
					g.name AS name,
					g.description AS description,
					g.members AS members,
					g.posts AS posts,
					g.avatar AS avatar,
					g.cover AS cover,
					g.privacy AS privacy,
					g.created_at AS created_at,
					g.updated_at AS updated_at,
					g.status AS status,
					exists((me)-[:JOIN]->(g)) AS is_joined
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
				`, params.Sort)

	paramsQuery := map[string]interface{}{
		"myuserid": myUserID,
		"skip":     params.Skip,
		"limit":    params.Limit,
	}
	res := []models.GroupJoin{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		return res, nil
	}
	return nil, nil
}

// Get func
// int64 int64
// models.Group error
func (service groupService) Get(groupID int64, myUserID int64) (models.GroupJoin, error) {
	stmt := `
				MATCH(me:User) WHERE ID(me) = {myUserID}
				MATCH (g:Group) WHERE ID(g) = {groupID}
				RETURN
					ID(g) AS id,
					g.name AS name,
					g.description AS description,
					g.members AS members,
					g.posts AS posts,
					g.avatar AS avatar,
					g.cover AS cover,
					g.privacy AS privacy,
					g.created_at AS created_at,
					g.updated_at AS updated_at,
					g.status AS status,
					exists((me)-[:JOIN]->(g)) AS is_joined
				`

	paramsQuery := map[string]interface{}{
		"groupID":  groupID,
		"myUserID": myUserID,
	}
	res := []models.GroupJoin{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.GroupJoin{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.GroupJoin{}, nil
}

// Delete func
// int64
// bool error
func (service groupService) Delete(groupID int64) (bool, error) {
	return true, nil
}

// func Create
// models.Group int64
// int64 error
func (service groupService) Create(group models.Group, myUserID int64) (int64, error) {
	p := neoism.Props{
		"name":        group.Name,
		"description": group.Description,
		"avatar":      group.Avatar,
		"cover":       group.Cover,
		"privacy":     group.Privacy,
		"members":     0,
		"posts":       0,
		"status":      0,
	}
	stmt := `
			MATCH(u:User) WHERE ID(u) = {myUserID}
			CREATE (g:Group{ props } )<-[r:CREATE]-(u)
			SET g.created_at = TIMESTAMP()
			RETURN ID(g) as id
			`
	params := map[string]interface{}{
		"props":    p,
		"myUserID": myUserID,
	}
	res := []struct {
		ID int64 `json:"id"`
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
	if len(res) > 0 {
		return res[0].ID, nil
	}
	return -1, nil
}

// Update func
// int64 models.Group
// models.Group error
func (service groupService) Update(groupID int64, newGroup models.InfoGroup) (models.Group, error) {
	stmt := `
		MATCH (g:Group)
		WHERE ID(g) = {groupID}
		SET g+= {p}, g.updated_at = TIMESTAMP()
		RETURN
			properties(g) AS group
		`
	params := neoism.Props{
		"groupID": groupID,
		"p":       newGroup,
	}
	res := []struct {
		Group models.Group `json:"group"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Group{}, err
	}
	if len(res) > 0 {
		if res[0].Group.ID >= 0 {
			return res[0].Group, nil
		}
	}
	return models.Group{}, nil
}

// CheckExistGroup func
// int64
// bool error
func (service groupService) CheckExistGroup(groupID int64) (bool, error) {
	stmt := `
		MATCH (g:Group)
		WHERE ID(g)={groupID}
		RETURN ID(g) AS id
		`
	params := neoism.Props{
		"groupID": groupID,
	}

	res := []struct {
		ID int64 `json:"id"`
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

	if len(res) > 0 {
		if res[0].ID == groupID {
			return true, nil
		}
	}
	return false, nil
}

// GetMembers func
// helpers.ParamsGetAll int64
// []models.UserJoinedObject error
func (service groupService) GetMembers(params helpers.ParamsGetAll, groupID int64) ([]models.UserJoinedObject, error) {
	stmt := fmt.Sprintf(`
	MATCH (g:Group)<-[j:JOIN]-(u:User)
	WHERE ID(g)= {groupID}
	WITH
		u{id:ID(u),joined_at:j.created_at,joined_by:"", .*} AS user
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	RETURN  collect(user) AS users
	`, params.Sort)

	paramsQuery := neoism.Props{
		"groupID": groupID,
		"skip":    params.Skip,
		"limit":   params.Limit,
	}
	res := []struct {
		Users []models.UserJoinedObject `json:"users"`
	}{}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		if res[0].Users[0].ID >= 0 {
			return res[0].Users, nil
		}
	}
	return nil, nil
}

// CreateMember func
// int64 int64
// bool error
func (service groupService) CreateMember(groupID int64, userID int64) (bool, error) {
	return false, nil
}

// CreateRequest func
// int64 int64
// bool error
func (service groupService) CreateRequest(groupID int64, userID int64) (bool, error) {
	return false, nil
}

// IncreasePosts func
// int64
// bool error
func (service groupService) IncreasePosts(groupID int64) (bool, error) {
	return false, nil
}

// DecreasePosts func
// int64
// bool error
func (service groupService) DecreasePosts(groupID int64) (bool, error) {
	return false, nil
}
