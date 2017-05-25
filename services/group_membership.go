package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// GroupMembershipServiceInterface include method list
type GroupMembershipServiceInterface interface {
	GetAll(params helpers.ParamsGetAll, groupID int64) ([]models.GroupMembership, error)
	Get(membershipID int64) (models.GroupMembership, error)
	Delete(membershipID int64) (bool, error)
	DeleteByUser(groupID int64, userID int64) (bool, error)
	Create(groupID int64, userID int64) (int64, error)
	Update(membership models.GroupMembership) (models.GroupMembership, error)
	CheckExistGroupMembership(groupID int64, userID int64) (bool, error)
}

// groupMembershipService struct
type groupMembershipService struct{}

// NewGroupMemberShipService to constructor
func NewGroupMembershipService() groupMembershipService {
	return groupMembershipService{}
}

// GetAll func
// models.ParamsGetAll int64
// []models.GroupMembership error
func (service groupMembershipService) GetAll(params helpers.ParamsGetAll, groupID int64) ([]models.GroupMembership, error) {
	stmt := fmt.Sprintf(`
			MATCH (g:Group)<-[r:JOIN]-(u:User)
			WHERE ID(g) = {groupID} AND r.role <>4
			WITH
				r{id:ID(r), .status, .created_at, .updated_at,
					group: g{id:ID(g), .name, .avatar},
					user: u{id:ID(u), .username, .full_name, .avatar}
					} AS membership
			ORDER BY %s
			SKIP {skip}
			LIMIT {limit}
			RETURN collect(membership) AS memberships
			`, params.Sort)
	paramsQuery := map[string]interface{}{
		"groupID": groupID,
		"skip":    params.Skip,
		"limit":   params.Limit,
	}
	res := []struct {
		Memberships []models.GroupMembership `json:"memberships"`
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
		return res[0].Memberships, nil
	}
	return nil, nil
}

// Create func
// int64 int64
// int64 error
func (service groupMembershipService) Create(groupID int64, userID int64) (int64, error) {
	// p := neoism.Props{
	// 	"name":        group.Name,
	// 	"description": group.Description,
	// 	"avatar":      group.Avatar,
	// }
	stmt := `
			MATCH(u:User) WHERE ID(u) = {userID}
			MATCH(g:Group) WHERE ID(g) = {groupID}
			CREATE (g)<-[r:JOIN]-(u)
			SET r.role=1, r.created_at = TIMESTAMP(), r.status = 1, g.members=g.members+1
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
		return -1, err
	}
	if len(res) > 0 {
		if res[0].ID >= 0 {
			return res[0].ID, nil
		}
	}
	return -1, nil
}

// Get func
// int64
// models.GroupMembership error
func (service groupMembershipService) Get(membershipID int64) (models.GroupMembership, error) {
	stmt := `
			MATCH (g:Group)<-[r:JOIN]-(u:User)
			WHERE ID(r) = {membershipID}
			RETURN
				ID(r) AS id, r.created_at AS created_at, r.updated_at AS updated_at, r.role AS role, r.status AS status,
				u{id:ID(u), .username, .full_name, .avatar} AS user,
				g{id:ID(g), .name, .avatar} AS group
			`
	params := map[string]interface{}{
		"membershipID": membershipID,
	}
	res := []models.GroupMembership{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.GroupMembership{}, err
	}
	if len(res) > 0 {
		if res[0].ID >= 0 {
			return res[0], nil
		}
	}
	return models.GroupMembership{}, nil
}

// Update func
// models.GroupMembership
// models.GroupMembership error
func (service groupMembershipService) Update(membership models.GroupMembership) (models.GroupMembership, error) {
	stmt := `
	MATCH (u:User)-[r:JOIN]->(g:Group)
	WHERE ID(r) = {membershipID}
	SET u += {p}, r.updated_at = TIMESTAMP()
	RETURN
		ID(r) AS id, r.created_at AS created_at, r.updated_at AS updated_at,
		r.request_message AS request_message, r.response_message AS response_message, r.status AS status,
		u{id:ID(u), .username, .full_name, .avatar} AS user,
		g{id:ID(g), .name, .avatar} AS group
	`
	p := map[string]interface{}{
		"status": membership.Status,
	}
	params := neoism.Props{
		"membershipID": membership.ID,
		"p":            p,
	}

	res := []models.GroupMembership{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.GroupMembership{}, err
	}

	if len(res) > 0 {
		if res[0].User.ID >= 0 {
			return res[0], nil
		}
	}
	return models.GroupMembership{}, nil
}

// Delete func
// int64
// bool error
func (service groupMembershipService) Delete(membershipID int64) (bool, error) {
	stmt := `
			MATCH (g:Group)<-[r:JOIN]-(u:User)
			WHERE ID(r) = {membershipID}
			SET g.members = g.members - 1
			DELETE r
			`
	params := map[string]interface{}{
		"membershipID": membershipID,
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

func (service groupMembershipService) DeleteByUser(groupID int64, userID int64) (bool, error) {
	stmt := `
			MATCH (g:Group)<-[r:JOIN]-(u:User)
			WHERE ID(g) = {groupID} AND ID(u) = {userID}
			SET g.members = g.members - 1
			DELETE r
			`
	params := map[string]interface{}{
		"groupID": groupID,
		"userID":  userID,
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

// CheckExistGroupMembership func
// int64 int64
// bool error
func (service groupMembershipService) CheckExistGroupMembership(groupID int64, userID int64) (bool, error) {
	stmt := `
		MATCH (u)-[r:JOIN]->(g:Group)
		WHERE ID(u)={userID} AND ID(g)= {groupID}
		RETURN ID(r) AS id
		`
	params := neoism.Props{
		"groupID": groupID,
		"userID":  userID,
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
