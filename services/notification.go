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

	// Update if action is excute
	UpdateLikeNotification(postID int64) (models.Notification, error)
	UpdateFollowNotification(userID int64, objectID int64) (models.Notification, error)
	UpdateStatusNotification(userID int64) (models.Notification, error)
	UpdatePhotoNotification(userID int64) (models.Notification, error)
	UpdateCommentNotification(postID int64) (models.Notification, error)
	UpdateMentionNotification(postID int64, userID int64, commentID int64) (models.Notification, error)
	UpdateLikedPostNotification(userID int64) (models.Notification, error)
	UpdateCommentedPostNotification(userID int64) (models.Notification, error)
	UpdateMentionedPostNotification(userID int64) (models.Notification, error)
}

// notificationService struct
type notificationService struct{}

// NewNotificationService to constructor
func NewNotificationService() notificationService {
	return notificationService{}
}

// GetAll func
// []models.Notification
func (service notificationService) GetAll(params helpers.ParamsGetAll, userID int64) ([]models.Notification, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		MATCH(u:User)-[h:HAS]->(n:Notification) WHERE ID(u) = {userID}
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

// GetSubcriberNotification
// int64
// []int64 error
func (service notificationService) GetNotificationSubcriber(notificationID int64) ([]int64, error) {
	stmt := `
				MATCH(u:User)-[h:HAS]->(n:Notification)
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
				MERGE (n)<-[h:HAS]-(u)
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
				MERGE (n)<-[h:HAS]-(u)
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
				MATCH(u)-[h:HAS]->(n:Notification)
				WHERE ID(n) = {notificationID}
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

//
func (service notificationService) UpdateFollowNotification(userID int64, objectID int64) (models.Notification, error) {
	stmt := `
			MATCH(u1:User)
			WHERE ID(u1) = {userID}
			MATCH(u1)-[f1:FOLLOW]->(u11:User)
			WITH f1, u11,u1
			ORDER BY f1.created_at DESC LIMIT 1
			MERGE (u1)-[g:GENERATE]->(n:Notification{action:{action}})
			SET
			n.actor = u1{id:ID(u1),username:u1.username,full_name:u1.full_name,avatar:u1.avatar},
			n.total_action = u1.followings,
			n.last_user= u11{id:ID(u11),username: u11.username, full_name: u11.full_name, avatar: u11.avatar},
			n.updated_at = TIMESTAMP()
			MATCH ()
			RETURN
				ID(n) AS id,
				apoc.convert.getJsonProperty(n,"actor") AS actor,
				CASE exists(n.last_comment) WHEN true THEN apoc.convert.getJsonProperty(n,"last_comment") END AS last_comment,
				CASE exists(n.last_post) WHEN true THEN apoc.convert.getJsonProperty(n,"last_post") END AS last_post,
				CASE exists(n.last_user) WHEN true THEN apoc.convert.getJsonProperty(n,"last_user") END AS last_user,
				CASE exists(n.last_mention) WHEN true THEN apoc.convert.getJsonProperty(n,"last_mention") END AS last_mention,
				n.updated_at AS updated_at,
				n.action AS action,
				n.total_action AS total_action
			`
	params := map[string]interface{}{
		"userID": userID,

		"action": configs.IActionFollow,
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
			MATCH (u:User)-[w]->(c:Comment)-[a]->(p:Post)
			WHERE ID(p) = {postID}
			WITH c,p,u
			ORDER BY c.created_at DESC LIMIT 1
			MATCH (u1:User)-[w]->(c1:Comment)-[a]->(p)
			WHERE TIMESTAMP() - c1.created_at < {limit_time}
			WITH c,p,u, count(c1) AS total_action
			MERGE (p)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET g.created_at = TIMESTAMP(),n.updated_at = c.created_at, n.total_action = total_action,n.last_comment = c{id:ID(c),message: c.message}, n.actor = u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar},last_post= p{id:ID(p), message:p.message}
			ON MATCH SET n.updated_at = c.created_at, n.total_action = total_action,n.last_comment = c{id:ID(c),message: c.message}, n.last_actor = u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar},last_post= p{id:ID(p), message:p.message}
			RETURN
				ID(n) AS id,
				n.actor AS actor,
				n.action AS action,
				n.total_action AS total_action,
				n.last_post AS last_post,
				n.last_comment AS last_comment,
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
			MATCH(u:User)-[l:LIKE]->(p:Post)
			WHERE ID(p) = {postID}
			WITH u,l,p
			ORDER BY l.created_at DESC LIMIT 1
			MATCH(:User)-[l1:LIKE]->(p)
			WHERE TIMESTAMP() - l1.created_at < {limit_time}
			WITH u,l,p,count(l1) AS like_count
			MERGE (p)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET
				g.created_at = TIMESTAMP(),
				n.actor =apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
				n.last_post=apoc.convert.toJson(p{id:ID(p),message:p.message}),
				n.total_action = like_count,
				n.updated_at = l.created_at
			ON MATCH SET
			n.actor =apoc.convert.toJson(u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}),
			n.last_post=apoc.convert.toJson(p{id:ID(p),message:p.message}),
			n.total_action = like_count,
			n.updated_at = l.created_at
			WITH u,l,p,n
			OPTIONAL MATCH (u1:User)-[:FOLLOW]->(p)
			MERGE (n)<-[h:HAS]-(u1)
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
			ON CREATE SET g.created_at = TIMESTAMP(),n.updated_at = p.created_at,n.last_post = p{id:ID(p), message:p.message}, n.total_action = total_action, n.actor= u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}
			ON MATCH SET n.updated_at = p.created_at,n.last_post = p{id:ID(p), message:p.message}, n.total_action = total_action, n.actor= u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}
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
func (service notificationService) UpdateMentionNotification(postID int64, userID int64, commentID int64) (models.Notification, error) {
	// ~doing ~needfix
	return models.Notification{}, nil
}

// UpdateLikedPostNotification
// int64
// models.Notification error
func (service notificationService) UpdateLikedPostNotification(userID int64) (models.Notification, error) {
	stmt := `
			MATCH(p:Post)<-[l:LIKE]-(u:User)
			WHERE ID(u) = {userID} AND p.privacy=1
			WITH p,l,u
			ORDER BY l.created_at DESC LIMIT 1
			MATCH(p1:Post)<-[l1:LIKE]-(u)
			WHERE TIMESTAMP() - l1.created_at < {limit_time} AND exists((p1)<-[:HAS]-(:Group))=false AND p1.privacy <>3
			WITH p,l,u,count(l1) AS total_action
			MERGE (u)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET g.created_at = TIMESTAMP(),n.updated_at = l.created_at,n.last_post = p{id:ID(p),message:p.message},n.total_action=total_action,n.actor=u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}
			ON MATCH SET n.updated_at = l.created_at,n.last_post = p{id:ID(p),message:p.message},n.total_action=total_action,n.actor=n.actor= u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar}
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
			MATCH(p:Post)<-[a]-(c:Comment)<-[w]-(u:User)
			WHERE ID(u) = {userID} AND p.privacy=1 AND exists((p)<-[:HAS]-(:Group))=false
			WITH p,c,u
			ORDER BY c.created_at DESC LIMIT 1
			MATCH(p1:Post)<-[a]-(c1:Comment)<-[w]-(u)
			WHERE TIMESTAMP() - c1.created_at < {limit_time} AND p1.privacy=1 AND exists((p1)<-[:HAS]-(:Group))=false
			WITH p,c,u,count(c1) AS total_action
			MERGE (u)-[g:GENERATE]->(n:Notification{action:{action}})
			ON CREATE SET g.created_at = TIMESTAMP(),n.updated_at = c.created_at,n.last_post = p{id:ID(p),message:p.message},n.total_action=total_action,n.actor=u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar},n.last_comment = c{id:ID(c),message:c.message}
			ON MATCH SET n.updated_at = c.created_at,n.last_post = p{id:ID(p),message:p.message},n.total_action=total_action,n.actor=n.actor= u{id:ID(u),username:u.username,full_name:u.full_name,avatar:u.avatar},n.last_comment = c{id:ID(c),message:c.message}
			RETURN
				ID(n) AS id,
				n.actor AS actor,
				n.action AS action,
				n.total_action AS total_action,
				n.last_post AS last_post,
				n.last_comment,
				"" AS title,
				"" AS message,
				n.updated_at AS updated_at,
				n.created_at AS created_at
			`
	params := map[string]interface{}{
		"userID":     userID,
		"time_limit": configs.ITwoDays,
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
