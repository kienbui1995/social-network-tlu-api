package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// GroupMembershipRequestServiceInterface include method list
type GroupMembershipRequestServiceInterface interface {
	GetAll(params helpers.ParamsGetAll, groupID int64) ([]models.GroupMembershipRequest, error)
	Get(requestID int64) (models.GroupMembershipRequest, error)
	Delete(requestID int64) (bool, error)
	DeleteByUser(groupID int64, userID int64) (bool, error)
	Create(request models.GroupMembershipSentRequest) (bool, error)
	Update(requestID int64, request models.GroupMembershipRequest) (models.GroupMembershipRequest, error)
	CheckExistRequest(groupID int64, userID int64) (bool, error)
}

// groupService struct
type groupMembershipRequestService struct{}

// NewGroupMembershipRequestService to constructor
func NewGroupMembershipRequestService() groupMembershipRequestService {
	return groupMembershipRequestService{}
}

// GetAll func
// helpers.ParamsGetAll int64
// []models.GroupMembershipRequest error
func (service groupMembershipRequestService) GetAll(params helpers.ParamsGetAll, groupID int64) ([]models.GroupMembershipRequest, error) {
	stmt := fmt.Sprintf(`
			MATCH (g:Group)<-[r:REQUEST]-(u:User)
			WHERE ID(g) = {groupID} AND r.status = {status}
			RETURN
				ID(r) AS id, r.created_at AS created_at, r.updated_at AS updated_at,
				r.request_message AS request_message, r.response_message AS response_message, r.status AS status,
				u{id:ID(u), .username, .full_name, .avatar} AS user,
				g{id:ID(g), .name, .avatar} AS group
			ORDER BY %s
			SKIP {skip}
			LIMIT {limit}
			`, params.Sort)
	paramsQuery := map[string]interface{}{
		"groupID": groupID,
		"skip":    params.Skip,
		"limit":   params.Limit,
		"status":  params.Type,
	}
	res := []models.GroupMembershipRequest{}
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
		if res[0].ID >= 0 {
			return res, nil
		}
	}
	return nil, nil
}

// Create func
// models.GroupMembershipSentRequest
// bool error
func (service groupMembershipRequestService) Create(request models.GroupMembershipSentRequest) (bool, error) {
	// p := neoism.Props{
	// 	"name":        group.Name,
	// 	"description": group.Description,
	// 	"avatar":      group.Avatar,
	// }
	stmt := `
			MATCH(u:User) WHERE ID(u) = {userID}
			MATCH(g:Group) WHERE ID(g) = {groupID}
			CREATE (g)<-[r:REQUEST]-(u)
			SET r.created_at = TIMESTAMP(), r.status = 1, r.request_message = {message}, g.pending_requests=g.pending_requests+1
			RETURN ID(r) as id
			`
	params := map[string]interface{}{
		"userID":  request.UserID,
		"groupID": request.GroupID,
		"message": request.Message,
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

// Create func
// int64
// models.GroupMembershipRequest error
func (service groupMembershipRequestService) Get(requestID int64) (models.GroupMembershipRequest, error) {
	// p := neoism.Props{
	// 	"name":        group.Name,
	// 	"description": group.Description,
	// 	"avatar":      group.Avatar,
	// }
	stmt := `
			MATCH (g:Group)<-[r:REQUEST]-(u:User)
			WHERE ID(r) = {requestID}
			RETURN
				ID(r) AS id, r.created_at AS created_at, r.updated_at AS updated_at,
				r.request_message AS request_message, r.response_message AS response_message, r.status AS status,
				u{id:ID(u), .username, .full_name, .avatar} AS user,
				g{id:ID(g), .name, .avatar} AS group
			`
	params := map[string]interface{}{
		"requestID": requestID,
	}
	res := []models.GroupMembershipRequest{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.GroupMembershipRequest{}, err
	}
	if len(res) > 0 {
		if res[0].ID >= 0 {
			return res[0], nil
		}
	}
	return models.GroupMembershipRequest{}, nil
}

// Update func
// int64 models.GroupMembershipRequest
// models.GroupMembershipRequest error
func (service groupMembershipRequestService) Update(requestID int64, request models.GroupMembershipRequest) (models.GroupMembershipRequest, error) {
	stmt := `
	MATCH (u:User)-[r:REQUEST]->(g:Group)
	WHERE ID(r) = {requestID}
	SET u += {p}, r.updated_at = TIMESTAMP()
	RETURN
		ID(r) AS id, r.created_at AS created_at, r.updated_at AS updated_at,
		r.request_message AS request_message, r.response_message AS response_message, r.status AS status,
		u{id:ID(u), .username, .full_name, .avatar} AS user,
		g{id:ID(g), .name, .avatar} AS group
	`

	params := neoism.Props{
		"requestID": requestID,
		"p":         request,
	}

	res := []models.GroupMembershipRequest{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.GroupMembershipRequest{}, err
	}

	if len(res) > 0 {
		if res[0].User.ID >= 0 {
			return res[0], nil
		}
	}
	return models.GroupMembershipRequest{}, nil
}

// Delete func
// int64
// bool error
func (service groupMembershipRequestService) Delete(requestID int64) (bool, error) {
	stmt := `
			MATCH (g:Group)<-[r:REQUEST]-(u:User)
			WHERE ID(r) = {requestID}
			DELETE r
			`
	params := map[string]interface{}{
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

// Delete func
// int64 int64
// bool error
func (service groupMembershipRequestService) DeleteByUser(groupID int64, userID int64) (bool, error) {
	stmt := `
			MATCH (g:Group)<-[r:REQUEST]-(u:User)
			WHERE ID(g) = {groupID} AND ID(u) = {userID}
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

// CheckExistRequest func
// int64 int64
// bool error
func (service groupMembershipRequestService) CheckExistRequest(groupID int64, userID int64) (bool, error) {
	stmt := `
		MATCH (u)-[r:REQUEST]->(g:Group)
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
