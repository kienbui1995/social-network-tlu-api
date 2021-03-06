package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// NotificationServiceInterface include method list
type NotificationServiceInterface interface {
	GetAll(params helpers.ParamsGetAll, userID int64) ([]models.Notification, error)
	Get(notificationID int64) (models.Notification, error)
	Delete(notificationID int64) (bool, error)
	Create(notification models.Notification, myUserID int64) (int64, error)
	Update(notification models.Notification) (models.Notification, error)

	GetNotificationSubcriber(notificationID int64) ([]int64, error)
	CreateNotificationSubcription(objectID int64, userID int64) (bool, error)
	CreateNotificationSubcriptionList(objectID int64, userIDs []int64) (bool, error)
	// user have seen noti and click
	SeenNotification(notificationID int64, userID int64) (bool, error)
	CheckSeenNotification(notificationID int64, userID int64) (bool, error)

	// Update if action is excute
	UpdateLikeNotification(postID int64) (models.Notification, error)
	UpdateFollowNotification(userID int64, objectID int64) (models.Notification, error)
	UpdateStatusNotification(userID int64) (models.Notification, error)
	UpdatePhotoNotification(userID int64) (models.Notification, error)
	UpdateCommentNotification(postID int64) (models.Notification, error)
	UpdateMentionNotification(postID int64, userID int64) (models.Notification, error)
	UpdateLikedPostNotification(userID int64) (models.Notification, error)
	UpdateCommentedPostNotification(userID int64) (models.Notification, error)
	UpdateMentionedPostNotification(userID int64) (models.Notification, error)
	UpdateRequestJoinNotification(groupID int64) (models.Notification, error)
}

// notificationService struct
type notificationService struct{}

// NewNotificationService to constructor
func NewNotificationService() NotificationServiceInterface {
	return notificationService{}
}

// GetAll func
// []models.Notification
func (service notificationService) GetAll(params helpers.ParamsGetAll, userID int64) ([]models.Notification, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		MATCH(u:User)-[h:REGISTERED]->(n:Notification) WHERE ID(u) = {userID}
		RETURN
			ID(n) AS id,
			CASE exists(n.actor) WHEN true THEN apoc.convert.getJsonProperty(n,"actor") END AS actor,
			CASE exists(n.last_comment) WHEN true THEN apoc.convert.getJsonProperty(n,"last_comment") END AS last_comment,
			CASE exists(n.last_post) WHEN true THEN apoc.convert.getJsonProperty(n,"last_post") END AS last_post,
			CASE exists(n.last_user) WHEN true THEN apoc.convert.getJsonProperty(n,"last_user") END AS last_user,
			CASE exists(n.last_mention) WHEN true THEN apoc.convert.getJsonProperty(n,"last_mention") END AS last_mention,
			n.updated_at AS updated_at,
			n.action AS action,
			n.total_action AS total_action
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
				`, params.Sort)

	paramsQuery := map[string]interface{}{
		"userID": userID,
		"skip":   params.Skip,
		"limit":  params.Limit,
	}
	res := []models.Notification{}

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
		fmt.Printf("res: %v\n", res)
		return res, nil
	}
	return nil, nil
}

// Get func
// int64
// models.Notification error
func (service notificationService) Get(notificationID int64) (models.Notification, error) {
	return models.Notification{}, nil
}

func (service notificationService) Delete(notificationID int64) (bool, error) {
	stmt := `
				MATCH(n:Notification)
				WHERE ID(n) = {notificationID}
				DETACH DELETE n
				`

	paramsQuery := map[string]interface{}{
		"notificationID": notificationID,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Create func
// models.Comment int64
// int64 error
func (service notificationService) Create(notification models.Notification, myUserID int64) (int64, error) {
	return -1, nil
}

// Update func
// models.Notification
// models.Notification error
func (service notificationService) Update(notification models.Notification) (models.Notification, error) {
	return models.Notification{}, nil
}

// Update func
// models.Notification
// models.Notification error
func (service notificationService) SeenAll(userID int64) (bool, error) {
	stmt := `
				MATCH(u:User)-[h:REGISTERED]->(n:Notification)
				WHERE ID(u) = {userID} AND h.seen_at= 0
				SET h.seen_at = TIMESTAMP()
				`
	paramsQuery := map[string]interface{}{
		"userID": userID,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	return true, nil
}

// GetSubcriberNotification
// int64
// []int64 error
func (service notificationService) GetNotificationSubcriber(notificationID int64) ([]int64, error) {
	stmt := `
				MATCH(u:User)-[h:REGISTERED]->(n:Notification)
				WHERE ID(n) = {notificationID}
				RETURN collect(ID(u)) AS ids
				`

	paramsQuery := map[string]interface{}{
		"notificationID": notificationID,
	}
	res := []struct {
		IDs []int64 `json:"ids"`
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
		return res[0].IDs, nil
	}
	return nil, nil
}

func (service notificationService) CreateNotificationSubcription(objectID int64, userID int64) (bool, error) {
	stmt := `
				MATCH (u:User) WHERE ID(u) = {userID}
				MATCH(obj)-[g:GENERATE]->(n:Notification)
				WHERE ID(obj) = {objectID}
				MERGE (n)<-[h:REGISTERED]-(u)
				RETURN ID(h) AS id
				`

	paramsQuery := map[string]interface{}{
		"userID":   userID,
		"objectID": objectID,
	}
	res := []struct {
		ID int64 `json:"id"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
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

func (service notificationService) CreateNotificationSubcriptionList(objectID int64, userIDs []int64) (bool, error) {
	stmt := `
				MATCH (u:User) WHERE ID(u) IN {userIDs}
				MATCH(obj)-[g:GENERATE]->(n:Notification)
				WHERE ID(obj) = {objectID}
				MERGE (n)<-[h:REGISTERED]-(u)
				ON CREATE h.created_at = TIMESTAMP(), h.seen=0
				RETURN ID(u) AS id
				`

	paramsQuery := map[string]interface{}{
		"userIDs":  userIDs,
		"objectID": objectID,
	}
	res := []struct {
		ID int64 `json:"id"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
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

// user have seen noti and click
func (service notificationService) SeenNotification(notificationID int64, userID int64) (bool, error) {
	stmt := `
				MATCH (u:User) WHERE ID(u) = {userID}
				MATCH(u)-[h:REGISTERED]->(n:Notification)
				WHERE ID(n) = {notificationID} AND h.seen_at = 0
				SET h.seen_at = TIMESTAMP()
				`

	paramsQuery := map[string]interface{}{
		"userID":         userID,
		"notificationID": notificationID,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (service notificationService) CheckSeenNotification(notificationID int64, userID int64) (bool, error) {
	stmt := `
				MATCH(u:User)-[h:REGISTERED]->(n:Notification)
				WHERE ID(n) = {notificationID} AND ID(u) = {userID}
				RETURN h.seen_at AS seen_at
				`

	paramsQuery := map[string]interface{}{
		"userID":         userID,
		"notificationID": notificationID,
	}
	res := []struct {
		SeenAt int64 `json:"seen_at"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}
	if len(res) > 0 {
		if res[0].SeenAt > 0 {
			return false, nil
		}
	}
	return true, nil
}

// UpdateFollowNotification func
// int64 int64
// int64 error
func (service notificationService) UpdateFollowNotification(userID int64, objectID int64) (models.Notification, error) {
	stmt := `
			MATCH(u:User)
			WHERE ID(u) = {userID}
			MATCH(u)-[f:FOLLOW]->(u1:User)
			WITH u,f,u1
			ORDER BY f.created_at DESC LIMIT 1
			MATCH(u)-[f_count:FOLLOW]->(u_count:User)
			WHERE TIMESTAMP() - f_count.created_at < {limit_time}
			WITH u,f,u1,count(f_count) AS total_action
			MERGE (u)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET
				g.created_at = TIMESTAMP(),
				n.updated_at = f.created_at,
				n.actor = apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.total_action = total_action,
				n.last_user= apoc.convert.toJson(u1{id:ID(u1),username: u1.username, full_name: u1.full_name, avatar: CASE u1.avatar THEN "" })
			ON MATCH SET
				n.updated_at = f.created_at,
				n.actor = apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.total_action = total_action,
				n.last_user= apoc.convert.toJson(u1{id:ID(u1),username: u1.username, full_name: u1.full_name, avatar: u1.avatar})
			WITH u,n
			OPTIONAL MATCH (u1:User)-[:FOLLOW]->(u)
			MERGE (n)<-[h:REGISTERED]-(u1)
			ON CREATE SET h.created_at = TIMESTAMP(), h.seen_at = 0
			ON MATCH SET h.seen_at = 0
			RETURN
				ID(n) AS id,
				apoc.convert.fromJsonMap(n.actor) AS actor,
				CASE exists(n.last_comment) WHEN true THEN apoc.convert.fromJsonMap(n.last_comment) END AS last_comment,
				CASE exists(n.last_post) WHEN true THEN apoc.convert.fromJsonMap(n.last_post) END AS last_post,
				CASE exists(n.last_user) WHEN true THEN apoc.convert.fromJsonMap(n.last_user) END AS last_user,
				CASE exists(n.last_mention) WHEN true THEN apoc.convert.fromJsonMap(n.last_mention) END AS last_mention,
				n.updated_at AS updated_at,
				n.action AS action,
				n.total_action AS total_action
			`
	params := map[string]interface{}{
		"userID":     userID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionFollow,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}

//
func (service notificationService) UpdateCommentNotification(postID int64) (models.Notification, error) {
	stmt := `
			MATCH (u:User)-[w]->(c:Comment)-[a]->(p:Post)<-[:POST]-(owner:User)
			WHERE ID(p) = {postID}
			WITH c,p,u,owner
			ORDER BY c.created_at DESC LIMIT 1
			MATCH (u1:User)-[w]->(c1:Comment)-[a]->(p)
			WHERE TIMESTAMP() - c1.created_at < {limit_time}
			WITH c,p,u, count(u1) AS total_action,owner
			MERGE (p)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET
				g.created_at = TIMESTAMP(),
				n.updated_at = c.created_at,
				n.total_action = total_action,
				n.last_comment = apoc.convert.toJson(c{id:ID(c),message: c.message, mentions: c.mentions}),
				n.actor = apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_post= apoc.convert.toJson(p{id:ID(p), message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}})
			ON MATCH SET
				n.updated_at = c.created_at,
				n.total_action = total_action,
				n.last_comment = apoc.convert.toJson(c{id:ID(c),message: c.message, mentions: c.mentions}),
				n.last_actor = apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_post= apoc.convert.toJson(p{id:ID(p),message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}})
			WITH u,p,n
			OPTIONAL MATCH (u1:User)-[:FOLLOW]->(p)
			MERGE (n)<-[h:REGISTERED]-(u1)
			ON CREATE SET h.created_at = TIMESTAMP(), h.seen_at = 0
			ON MATCH SET h.seen_at = 0
			RETURN
				ID(n) AS id,
				apoc.convert.fromJsonMap(n.actor) AS actor,
				n.action AS action,
				n.total_action AS total_action,
				apoc.convert.fromJsonMap(n.last_post) AS last_post,
				apoc.convert.fromJsonMap(n.last_comment) AS last_comment,
				"" AS title,
				"" AS message,
				n.updated_at AS updated_at,
				n.created_at AS created_at
			`
	params := map[string]interface{}{
		"postID":     postID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionComment,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}

// UpdateLikeNotification func
// int64 int64
// models.Notification error
func (service notificationService) UpdateLikeNotification(postID int64) (models.Notification, error) {
	stmt := `
			MATCH(u:User)-[l:LIKE]->(p:Post)<-[r:POST]-(owner:User)
			WHERE ID(p) = {postID}
			WITH u,l,p,owner
			ORDER BY l.created_at DESC LIMIT 1
			MATCH(:User)-[l1:LIKE]->(p)
			WHERE TIMESTAMP() - l1.created_at < {limit_time}
			WITH u,l,p,count(l1) AS like_count,owner
			MERGE (p)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET
				g.created_at = TIMESTAMP(),
				n.actor =apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_post=apoc.convert.toJson(p{id:ID(p),message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}}),
				n.total_action = like_count,
				n.updated_at = l.created_at
			ON MATCH SET
			n.actor =apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
			n.last_post=apoc.convert.toJson(p{id:ID(p),message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}}),
			n.total_action = like_count,
			n.updated_at = l.created_at
			WITH u,l,p,n
			OPTIONAL MATCH (u1:User)-[:FOLLOW]->(p)
			MERGE (n)<-[h:REGISTERED]-(u1)
			ON CREATE SET h.created_at = TIMESTAMP(), h.seen_at = 0
			ON MATCH SET h.seen_at = 0
			RETURN
				ID(n) AS id,
				apoc.convert.fromJsonMap(n.actor) AS actor,
				n.action AS action,
				n.total_action AS total_action,
				apoc.convert.fromJsonMap(n.last_post) AS last_post,
				"" AS title,
				"" AS message,
				n.updated_at AS updated_at,
				n.created_at AS created_at
			`
	params := map[string]interface{}{
		"postID":     postID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionLike,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}

// Notification with status
// int64 int64
// models.Notification error
func (service notificationService) UpdateStatusNotification(userID int64) (models.Notification, error) {
	stmt := `
			MATCH(p:Status:Post)<-[:POST]-(u:User)
			WHERE ID(u) = {userID} AND p.privacy <>3
			WITH p,u
			ORDER BY p.created_at DESC LIMIT 1
			MATCH(p1:Status:Post)<-[:POST]-(u)
			WHERE p1.privacy <>3 AND TIMESTAMP()- p1.created_at < {limit_time}
			WITH p,u, count(p1) AS total_action
			MERGE (u)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET
				g.created_at = TIMESTAMP(),
				n.updated_at = p.created_at,
				n.last_post = apoc.convert.toJson(p{id:ID(p), message:p.message}),
				n.total_action = total_action,
				n.actor= apoc.convert.toJson(u{id:ID(u), username:u.username, full_name:u.full_name, avatar:u.avatar})
			ON MATCH SET
			n.updated_at = p.created_at,
			n.last_post = apoc.convert.toJson(p{id:ID(p), message:p.message}),
			n.total_action = total_action,
			n.actor= apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar})
			RETURN
				ID(n) AS id,
				n.actor AS actor,
				n.action AS action,
			 	n.total_action AS total_action,
				n.last_post AS last_post,
				"" AS title,
				"" AS message,
				n.updated_at AS updated_at,
				n.created_at AS created_at
			`
	params := map[string]interface{}{
		"userID":     userID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionPostStatus,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}

// Notification with photo
// int64 int64
// models.Notification error
func (service notificationService) UpdatePhotoNotification(userID int64) (models.Notification, error) {
	stmt := `
			MATCH(p:Photo:Post)<-[:POST]-(u:User)
			WHERE ID(u) = {userID} AND p.privacy <>3
			WITH p,u
			ORDER BY p.created_at DESC LIMIT 1
			MATCH(p1:Status:Post)<-[:POST]-(u)
			WHERE p1.privacy <>3 AND TIMESTAMP()- p1.created_at < {limit_time}
			WITH p,u, count(p1) AS total_action
			MERGE (u)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET g.created_at = TIMESTAMP(),n.updated_at = p.created_at,n.last_post = p{id:ID(p), message:p.message,photo:p.photo}, n.total_action = total_action, n.actor= u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}
			ON MATCH SET n.updated_at = p.created_at,n.last_post = p{id:ID(p), message:p.message,photo:p.photo}, n.total_action = total_action, n.actor= u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}
			RETURN
				ID(n) AS id,
				n.actor AS actor,
				n.action AS action,
			 	n.total_action AS total_action,
				n.last_post AS last_post,
				"" AS title,
				"" AS message,
				n.updated_at AS updated_at,
				n.created_at AS created_at
			`
	params := map[string]interface{}{
		"userID":     userID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionPostStatus,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}

// Notification with mention
// int64 int64 int64
// models.Notification error
func (service notificationService) UpdateMentionNotification(postID int64, userID int64) (models.Notification, error) {
	stmt := `
			MATCH (actor:User)-[w]->(c:Comment)-[a]->(p:Post)<-[:POST]-(owner:User)
			WHERE ID(p) = {postID} AND exists(c.mentions)=true
			WITH actor, c,p, owner
			ORDER BY c.created_at DESC LIMIT 1
			MATCH (u1:User)-[w]->(c1:Comment)-[a]->(p)
			WHERE TIMESTAMP() - c1.created_at < {limit_time}
			WITH c,p,u, count(u1) AS total_action,owner
			MERGE (p)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET
				g.created_at = TIMESTAMP(),
				n.updated_at = c.created_at,
				n.total_action = total_action,
				n.last_comment = apoc.convert.toJson(c{id:ID(c),message: c.message}),
				n.actor = apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_post= apoc.convert.toJson(p{id:ID(p), message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}})
			ON MATCH SET
				n.updated_at = c.created_at,
				n.total_action = total_action,
				n.last_comment = apoc.convert.toJson(c{id:ID(c),message: c.message}),
				n.last_actor = apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_post= apoc.convert.toJson(p{id:ID(p),message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}})
			WITH u,p,n
			OPTIONAL MATCH (u1:User)-[:FOLLOW]->(p)
			MERGE (n)<-[h:REGISTERED]-(u1)
			ON CREATE SET h.created_at = TIMESTAMP(), h.seen_at = 0
			ON MATCH SET h.seen_at = 0
			RETURN
				ID(n) AS id,
				apoc.convert.fromJsonMap(n.actor) AS actor,
				n.action AS action,
				n.total_action AS total_action,
				apoc.convert.fromJsonMap(n.last_post) AS last_post,
				apoc.convert.fromJsonMap(n.last_comment) AS last_comment,
				"" AS title,
				"" AS message,
				n.updated_at AS updated_at,
				n.created_at AS created_at
			`
	params := map[string]interface{}{
		"postID":     postID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionComment,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}

// UpdateLikedPostNotification
// int64
// models.Notification error
func (service notificationService) UpdateLikedPostNotification(userID int64) (models.Notification, error) {
	stmt := `
				MATCH(u:User)-[l:LIKE]->(p:Post)<-[:POST]-(owner:User)
				WHERE ID(p) = {userID} AND p.privacy=1
				WITH u,l,p,owner
				ORDER BY l.created_at DESC LIMIT 1
				MATCH(u)-[l1:LIKE]->(p1:Post)
				WHERE TIMESTAMP() - l1.created_at < {limit_time} AND exists((p1)<-[:HAS]-(:Group))=false AND p1.privacy <>3
				WITH u,l,p,count(l1) AS like_count,owner
				MERGE (u)-[g:GENERATE]->(n:Notification{action:{action}})
				ON CREATE SET
					g.created_at = TIMESTAMP(),
					n.actor =apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
					n.last_post=apoc.convert.toJson(p{id:ID(p),message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}}),
					n.total_action = like_count,
					n.updated_at = l.created_at
				ON MATCH SET
				n.actor =apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_post=apoc.convert.toJson(p{id:ID(p),message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}}),
				n.total_action = like_count,
				n.updated_at = l.created_at
				WITH u,l,p,n
				OPTIONAL MATCH (u1:User)-[:FOLLOW]->(u)
				MERGE (n)<-[h:REGISTERED]-(u1)
				ON CREATE SET h.created_at = TIMESTAMP(), h.seen_at = 0
				ON MATCH SET h.seen_at = 0
				RETURN
					ID(n) AS id,
					apoc.convert.fromJsonMap(n.actor) AS actor,
					n.action AS action,
					n.total_action AS total_action,
					apoc.convert.fromJsonMap(n.last_post) AS last_post,
					"" AS title,
					"" AS message,
					n.updated_at AS updated_at,
					n.created_at AS created_at
			`
	params := map[string]interface{}{
		"userID":     userID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionLikedPost,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}

// UpdateCommentedPostNotification
// int64
// models.Notification error
func (service notificationService) UpdateCommentedPostNotification(userID int64) (models.Notification, error) {
	stmt := `
			MATCH (owner:User)-[:POST]->(p:Post)<-[a]-(c:Comment)<-[w]-(u:User)
			WHERE ID(u) = {userID} AND p.privacy=1 AND exists((p)<-[:HAS]-(:Group))=false
			WITH p,c,u,owner
			ORDER BY c.created_at DESC LIMIT 1
			MATCH(p1:Post)<-[a]-(c1:Comment)<-[w]-(u)
			WHERE TIMESTAMP() - c1.created_at < {limit_time} AND p1.privacy=1 AND exists((p1)<-[:HAS]-(:Group))=false
			WITH p,c,u,count(c1) AS total_action,owner
			MERGE (u)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET
				g.created_at = TIMESTAMP(),
				n.updated_at = c.created_at,
				n.last_post = apoc.convert.toJson(p{id:ID(p),message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}}),
				n.total_action=total_action,
				n.actor=apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_comment = apoc.convert.toJson(c{id:ID(c),message:c.message,mentions: c.mentions})
			ON MATCH SET
				n.updated_at = c.created_at,
				n.last_post = apoc.convert.toJson(p{id:ID(p),message:p.message,photo:p.photo,owner:owner{id:ID(owner),username:owner.username,full_name:owner.full_name,avatar:owner.avatar}}),
				n.total_action=total_action,
				n.actor=apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_comment = apoc.convert.toJson(c{id:ID(c),message:c.message, mentions: c.mentions})
			RETURN
				ID(n) AS id,
				apoc.convert.fromJsonMap(n.actor) AS actor,
				n.action AS action,
				n.total_action AS total_action,
				apoc.convert.fromJsonMap(n.last_post) AS last_post,
				apoc.convert.fromJsonMap(n.last_comment) AS last_comment,
				"" AS title,
				"" AS message,
				n.updated_at AS updated_at,
				n.created_at AS created_at
			`
	params := map[string]interface{}{
		"userID":     userID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionCommentedPost,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}

// Notification with mention
// int64 int64 int64
// models.Notification error
func (service notificationService) UpdateMentionedPostNotification(userID int64) (models.Notification, error) {
	// ~doing ~needfix
	return models.Notification{}, nil
}

// UpdateFollowNotification func
// int64 int64
// int64 error
func (service notificationService) UpdateRequestJoinNotification(groupID int64) (models.Notification, error) {
	stmt := `
		MATCH(g:Group)
		WHERE ID(g) = {groupID}
		MATCH(g)<-[f:JOIN{status:0}]-(u1:User)
		WITH g,f,u1
		ORDER BY f.created_at DESC LIMIT 1
		MATCH(g)<-[f_count:JOIN{status:0}]-(u_count:User)
		WHERE TIMESTAMP() - f_count.created_at < {limit_time}
		WITH g,f,u1,count(f_count) AS total_action
		MERGE (g)-[gen:GENERATE]->(n:Notification{action:{action}})
		ON CREATE SET
			gen.created_at = TIMESTAMP(),
			n.updated_at = f.created_at,
			n.actor = apoc.convert.toJson(u1{id:ID(u1),username:u1.username,full_name:u1.full_name,avatar:u1.avatar}),
			n.total_action = total_action,
			n.group= apoc.convert.toJson(g{id:ID(g),name: g.name, avatar: CASE exists(g.avatar) WHEN true THEN g.avatar ELSE "" END })
		ON MATCH SET
			n.updated_at = f.created_at,
			n.actor = apoc.convert.toJson(u1{id:ID(u1),username:u1.username,full_name:u1.full_name,avatar:u1.avatar}),
			n.total_action = total_action,
			n.group= apoc.convert.toJson(g{id:ID(g),name: g.name, avatar: CASE exists(g.avatar) WHEN true THEN g.avatar ELSE "" END })
		WITH g,n
		OPTIONAL MATCH (u:User)-[join:JOIN{status:1}]->(g)
		WHERE join.role =2 OR join.role = 3
		MERGE (n)<-[h:REGISTERED]-(u)
		ON CREATE SET h.created_at = TIMESTAMP(), h.seen_at = 0
		ON MATCH SET h.seen_at = 0
		RETURN
			ID(n) AS id,
			apoc.convert.fromJsonMap(n.actor) AS actor,
			CASE exists(n.last_comment) WHEN true THEN apoc.convert.fromJsonMap(n.last_comment) END AS last_comment,
			CASE exists(n.last_post) WHEN true THEN apoc.convert.fromJsonMap(n.last_post) END AS last_post,
			CASE exists(n.last_user) WHEN true THEN apoc.convert.fromJsonMap(n.last_user) END AS last_user,
			CASE exists(n.last_mention) WHEN true THEN apoc.convert.fromJsonMap(n.last_mention) END AS last_mention,
			CASE exists(n.group) WHEN true THEN apoc.convert.fromJsonMap(n.group) END AS group,
			n.updated_at AS updated_at,
			n.action AS action,
			n.total_action AS total_action
			`
	params := map[string]interface{}{
		"groupID":    groupID,
		"limit_time": configs.ITwoDays,
		"action":     configs.IActionRequestJoinGroup,
	}
	res := []models.Notification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Notification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Notification{}, nil
}
