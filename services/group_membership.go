package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// GroupMembershipServiceInterface include method list
type GroupMembershipServiceInterface interface {
	GetAll(params helpers.ParamsGetAll, groupID int64, myUserID int64) ([]models.GroupMembership, error)
	Get(membershipID int64) (models.GroupMembership, error)
	Delete(membershipID int64) (bool, error)
	DeleteByUser(groupID int64, userID int64) (bool, error)
	Create(groupID int64, userID int64) (int64, error)
	Update(membership models.GroupMembership) (models.GroupMembership, error)
	CheckExistGroupMembership(groupID int64, userID int64) (bool, error)
}

// groupMembershipService struct
type groupMembershipService struct{}

// NewGroupMembershipService to constructor
func NewGroupMembershipService() GroupMembershipServiceInterface {
	return groupMembershipService{}
}

// GetAll func
// models.ParamsGetAll int64
// []models.GroupMembership error
func (service groupMembershipService) GetAll(params helpers.ParamsGetAll, groupID int64, myUserID int64) ([]models.GroupMembership, error) {

	var stmt string
	var paramsQuery neoism.Props
	var role = -1
	if params.Properties["role"] != nil {
		if params.Properties["role"].(string) == configs.SAdmin {
			role = 2
		} else if params.Properties["role"].(string) == configs.SBlocked {
			role = 4
		} else if params.Properties["role"].(string) == configs.SCreator {
			role = 3
		} else if params.Properties["role"].(string) == configs.SMember {
			role = 1
		} else if params.Properties["role"].(string) == configs.SPending {
			role = 0
		}
	}
	if role >= 1 && role <= 3 { // member/admin/creator
		stmt = fmt.Sprintf(`
				MATCH (me:User) WHERE ID(me) = {myUserID}
				MATCH (g:Group)<-[r:JOIN{role: {role} }]-(u:User)
				WHERE ID(g) = {groupID}
				WITH
					r{id:ID(r), .*,
						user: u{id:ID(u), .username, .full_name, .avatar},
						can_edit:
						CASE r.role WHEN 1 THEN CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
																				 WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
																		END
												WHEN 2 THEN CASE WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
																		END
												WHEN 4 THEN CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
																				 WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
																		END
												ELSE false
						END,
						can_delete:
						CASE r.role WHEN 1 THEN CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
																		WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
		 									 							END
												WHEN 2 THEN CASE WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
																		END
												WHEN 4 THEN CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
																		WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
		 									 							END

												ELSE false
						END

					} AS membership
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
				RETURN collect(membership) AS memberships
				`, "membership."+params.Sort)

		paramsQuery = map[string]interface{}{
			"myUserID": myUserID,
			"groupID":  groupID,
			"skip":     params.Skip,
			"limit":    params.Limit,
			"role":     role,
		}
	} else if role == 4 { // blocked
		stmt = fmt.Sprintf(`
		MATCH (me:User) WHERE ID(me) = {myUserID}
		MATCH (g:Group)<-[r:JOIN{role:4}]-(u:User)
		WHERE ID(g) = {groupID} AND r.status = 1
		WITH
			r{id:ID(r), .*,
				user: u{id:ID(u), .username, .full_name, .avatar},
				can_edit:
				CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
						 WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
						 ELSE false
				END,
				can_delete:
				CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
						 WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
						 WHEN ID(me) = ID(u) THEN true
						 ELSE false
				END

			} AS membership
		ORDER BY %s
		SKIP {skip}
		LIMIT {limit}
		RETURN collect(membership) AS memberships
		`, "membership."+params.Sort)
		paramsQuery = map[string]interface{}{
			"myUserID": myUserID,
			"groupID":  groupID,
			"skip":     params.Skip,
			"limit":    params.Limit,
		}
	} else if role == 0 { // pending request
		stmt = fmt.Sprintf(`
		MATCH (me:User) WHERE ID(me) = {myUserID}
		MATCH (g:Group)<-[r:JOIN{status: 0}]-(u:User)
		WHERE ID(g) = {groupID}
		WITH
			r{id:ID(r), .*,
				user: u{id:ID(u), .username, .full_name, .avatar},
				can_edit:
				CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
						 WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
						 ELSE false
				END,
				can_delete:
				CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
						 WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
						 WHEN ID(me) = ID(u) THEN true
						 ELSE false
				END

			} AS membership
		ORDER BY %s
		SKIP {skip}
		LIMIT {limit}
		RETURN collect(membership) AS memberships
		`, "membership."+params.Sort)
		paramsQuery = map[string]interface{}{
			"myUserID": myUserID,
			"groupID":  groupID,
			"skip":     params.Skip,
			"limit":    params.Limit,
		}
	} else { // all member (include admin/creator)
		stmt = fmt.Sprintf(`
				MATCH (me:User) WHERE ID(me) = {myUserID}
				MATCH (g:Group)<-[r:JOIN]-(u:User)
				WHERE ID(g) = {groupID} AND r.role <>4 AND r.status = 1
				WITH
					r{id:ID(r), .*,
						user: u{id:ID(u), .username, .full_name, .avatar},
						can_edit:
						CASE r.role WHEN 1 THEN CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
																				 WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
																		END
												WHEN 2 THEN CASE WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
																		END
												WHEN 4 THEN CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
																				 WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
																		END
												ELSE false
						END,
						can_delete:
						CASE r.role WHEN 1 THEN CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
																		WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
		 									 							END
												WHEN 2 THEN CASE WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
																		END
												WHEN 4 THEN CASE WHEN exists((me)-[:JOIN{role:2}]->(g)) THEN true
																		WHEN exists((me)-[:JOIN{role:3}]->(g)) THEN true
		 									 							END
												ELSE false
						END

					} AS membership
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
				RETURN collect(membership) AS memberships
				`, "membership."+params.Sort)

		paramsQuery = map[string]interface{}{
			"myUserID": myUserID,
			"groupID":  groupID,
			"skip":     params.Skip,
			"limit":    params.Limit,
			"role":     role,
		}
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
			SET r.created_at = TIMESTAMP(),
					r += CASE g.privacy WHEN 1 THEN {status: 1, role: 1}
															WHEN 2 THEN {status: 0}
							 END,
					g.members = CASE g.privacy WHEN 1 THEN g.members+1
																		 ELSE g.members
											END,
					g.pending_requests = CASE g.privacy WHEN 2 THEN g.pending_requests+1
					 																		ELSE g.pending_requests
															 END
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
	SET g.pending_requests = CASE WHEN r.status = 0 AND {status}= 1 THEN  g.pending_requests -1  ELSE g.pending_requests END,
			g.members = CASE WHEN r.status = 0 AND {status}= 1 THEN  g.members +1
											 WHEN r.role = 4 AND ( {role}= 1 OR {role} =2) THEN  g.members +1
											 WHEN (r.role = 1 OR r.role=2) AND {role}= 4 THEN  g.members -1
											 ELSE g.members END,
			r.updated_at = TIMESTAMP(), r.role = {role}, r.status = {status}
	RETURN
		ID(r) AS id, r.created_at AS created_at, r.updated_at AS updated_at,
 		r.role AS role, r.status AS status,
		u{id:ID(u), .username, .full_name, .avatar} AS user,
		g{id:ID(g), .name, .avatar} AS group
	`

	params := neoism.Props{
		"membershipID": membership.ID,
		"status":       membership.Status,
		"role":         membership.Role,
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
			SET g.members = CASE r.status WHEN 1 THEN g.members - 1
																		ELSE g.members
											END,
					g.pending_requests = CASE r.status WHEN 0 THEN g.pending_requests - 1
																		         ELSE g.pending_requests
											END
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
			SET g.members = CASE r.status WHEN 1 THEN g.members - 1
																		ELSE g.members
											END,
					g.pending_requests = CASE r.status WHEN 0 THEN g.pending_requests - 1
																		         ELSE g.pending_requests
											END
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
