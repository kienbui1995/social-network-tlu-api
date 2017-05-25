package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
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
	CheckUserRole(groupID int64, userID int64) (int, error)
	GetJoinedGroup(params helpers.ParamsGetAll, userID int64, myUserID int64) ([]models.GroupJoin, error)

	GetMembers(params helpers.ParamsGetAll, groupID int64) ([]models.GroupMember, error)
	GetPendingUsers(params helpers.ParamsGetAll, groupID int64) ([]models.PendingUser, error)
	GetBlockedUsers(params helpers.ParamsGetAll, groupID int64) ([]models.UserFollowObject, error)

	CreateMember(groupID int64, userID int64) (bool, error)

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
					exists((me)-[:JOIN{role:1}]->(g)) AS is_member,
					exists((me)-[:REQUEST{status:1}]->(g)) AS is_pending,
					g.privacy = 2  and exists((me)-[:JOIN]->(g))=false AS can_request,
					g.privacy = 1  and exists((me)-[:JOIN]->(g))=false AS can_join,
					exists((me)-[:JOIN{role:2}]->(g)) OR exists((me)-[:JOIN{role:3}]->(g)) AS is_admin
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
					exists((me)-[:JOIN{role:1}]->(g)) AS is_member,
					exists((me)-[:REQUEST{status:1}]->(g)) AS is_pending,
					g.privacy = 2  AND exists((me)-[:JOIN]->(g))=false AND exists((me)-[:REQUEST]->(g))=false AS can_request,
					g.privacy = 1  AND exists((me)-[:JOIN]->(g))=false AS can_join,
					exists((me)-[:JOIN{role:2}]->(g)) OR exists((me)-[:JOIN{role:3}]->(g)) AS is_admin
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
		"name":             group.Name,
		"description":      group.Description,
		"avatar":           group.Avatar,
		"cover":            group.Cover,
		"privacy":          group.Privacy,
		"pending_requests": 0,
		"members":          0,
		"posts":            0,
	}
	stmt := `
			MATCH(u:User) WHERE ID(u) = {myUserID}
			CREATE (g:Group{props})<-[r:JOIN{role:3, status:1, created_at: TIMESTAMP()}]-(u)
			SET g.created_at = TIMESTAMP(), g.members= 1
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
			g{id:ID(g), .*} AS group
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

// CheckUserRole func
// int64 int64
// int error
func (service groupService) CheckUserRole(groupID int64, userID int64) (int, error) {
	stmt := `
	MATCH (g:Group)	WHERE ID(g)= {groupID}
	MATCH (u:User) WHERE ID(u) = {userID}
	RETURN
		exists((u)-[:JOIN{role:1}]->(g)) AS is_member,
		exists((u)-[:JOIN{role:2}]->(g)) OR exists((u)-[:JOIN{status:3}]->(g)) AS is_admin,
		exists((u)-[:REQUEST{status:1}]->(g)) AS pending,
		exists((u)-[:REQUEST{status:2}]->(g)) AS declined,
		exists((u)-[:JOIN{role:4}]->(g)) AS blocked
	`

	paramsQuery := neoism.Props{
		"groupID": groupID,
		"userID":  userID,
	}
	res := []struct {
		IsMember bool `json:"is_member"`
		IsAdmin  bool `json:"is_admin"`
		Pending  bool `json:"pending"`
		Declined bool `json:"declined"`
		Blocked  bool `json:"blocked"`
	}{}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return -1, err
	}
	if len(res) > 0 {
		if res[0].Blocked {
			return configs.IBlocked, nil
		}
		if res[0].Declined {
			return configs.IDeclined, nil
		}
		if res[0].Pending {
			return configs.IPending, nil
		}
		if res[0].IsMember {
			return configs.IMember, nil
		}
		if res[0].IsAdmin {
			return configs.IAdmin, nil
		}
	}
	return -1, nil
}

// GetJoinedGroup func
// helpers.ParamsGetAll int64 int64
// []models.group error
func (service groupService) GetJoinedGroup(params helpers.ParamsGetAll, userID int64, myUserID int64) ([]models.GroupJoin, error) {
	var stmt string
	stmt = fmt.Sprintf(`
				MATCH(me:User) WHERE ID(me) = {myuserid}
				MATCH (g:Group)<-[j:JOIN]-(u:User)
				WHERE ID(u) = {userID}
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
					exists((me)-[:JOIN{role:1}]->(g)) AS is_member,
					exists((me)-[:JOIN{role:2}]->(g)) AS is_admin,
					exists((me)-[:REQUEST{status:1}]->(g)) AS is_pending,
					g.privacy = 2  and exists((me)-[:JOIN]->(g))=false AS can_request,
					g.privacy = 1  and exists((me)-[:JOIN]->(g))=false AS can_join,
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
				`, params.Sort)

	paramsQuery := map[string]interface{}{
		"myuserid": myUserID,
		"userID":   userID,
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

// GetMembers func
// helpers.ParamsGetAll int64
// []models.UserJoinedObject error
func (service groupService) GetMembers(params helpers.ParamsGetAll, groupID int64) ([]models.GroupMember, error) {
	stmt := fmt.Sprintf(`
		MATCH (g:Group)<-[j:JOIN]-(u:User)
		WHERE ID(g)= {groupID} AND j.role <> 4
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
		Users []models.GroupMember `json:"users"`
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

// GetPendingUsers func
// helpers.ParamsGetAll int64
// []models.UserJoinedObject error
func (service groupService) GetPendingUsers(params helpers.ParamsGetAll, groupID int64) ([]models.PendingUser, error) {
	return nil, nil
}

// GetBlockedUsers func
// helpers.ParamsGetAll int64
// []models.UserObject error
func (service groupService) GetBlockedUsers(params helpers.ParamsGetAll, groupID int64) ([]models.UserFollowObject, error) {
	return nil, nil
}

// CreateMember func
// int64 int64
// bool error
func (service groupService) CreateMember(groupID int64, userID int64) (bool, error) {

	// p := neoism.Props{
	// 	"name":        group.Name,
	// 	"description": group.Description,
	// 	"avatar":      group.Avatar,
	// }
	stmt := `
			MATCH(u:User) WHERE ID(u) = {userID}
			MATCH(g:Group) WHERE ID(g) = {groupID}
			CREATE (g)<-[r:JOIN]-(u)
			SET r.created_at = TIMESTAMP()
			RETURN ID(r) as id
			`
	params := map[string]interface{}{
		"userID":  userID,
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
		if res[0].ID >= 0 {
			return true, nil
		}
	}
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
