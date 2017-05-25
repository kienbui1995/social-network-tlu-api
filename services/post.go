package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// PostServiceInterface include method list
type PostServiceInterface interface {
	GetAll(params helpers.ParamsGetAll, userID int64, myUserID int64) ([]models.Post, error)
	Get(postID int64, myUserID int64) (models.Post, error)
	Delete(postID int64) (bool, error)
	Create(post models.Post, myUserID int64) (int64, error)
	Update(post models.Post) (models.Post, error)
	CheckExistPost(postID int64) (bool, error)
	GetUserIDByPostID(postID int64) (int64, error)
	IncreasePosts(userID int64) (bool, error)
	DecreasePosts(userID int64) (bool, error)

	// work with likes
	CreateLike(postID int64, userID int64) (int, error)
	GetLikes(postID int64, myUserID int64, params helpers.ParamsGetAll) ([]models.UserLikedObject, error)
	DeleteLike(postID int64, userID int64) (int, error)
	CheckExistLike(postID int64, userID int64) (bool, error)
	IncreaseLikes(postID int64) (bool, error)
	DecreaseLikes(postID int64) (bool, error)
	CheckPostInteractivePermission(postID int64, userID int64) (bool, error)

	// work with FOLLOW
	CreateFollow(postID int64, userID int64) (int64, error)
	DeleteFollow(postID int64, userID int64) (bool, error)
	CheckExistFollow(postID int64, userID int64) (bool, error)

	// work with users (can_mention, mentioned, liked, commented, followed)
	GetUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error)
	GetCanMentionedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error)
	GetMentionedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error)
	GetLikedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error)
	GetCommentedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error)
	GetFollowedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error)
}

// postService struct
type postService struct{}

// NewPostService to constructor
func NewPostService() postService {
	return postService{}
}

// GetAll func
// helpers.ParamsGetAll
// models.Post error
func (service postService) GetAll(params helpers.ParamsGetAll, userID int64, myUserID int64) ([]models.Post, error) {
	var stmt string
	if params.Type == configs.SPostPhoto {
		stmt = fmt.Sprintf(`
		    MATCH(u:User) WHERE ID(u) = {userid}
				MATCH(me:User) WHERE ID(me) = {myuserid}
		  	MATCH (s:Photo:Post)<-[r:POST]-(u)
				WHERE s.privacy = 1 OR (s.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR {userid} = {myuserid}
				RETURN
					ID(s) AS id,
					substring(s.message,0,250) AS message, length(s.message)>250 AS summary,
					CASE s.photo when null then "" else s.photo end AS photo,
					s.created_at AS created_at, s.updated_at AS updated_at,
					CASE s.privacy when null then 1 else s.privacy end AS privacy, CASE s.status when null then 1 else s.status end AS status,
					s.likes AS likes, s.comments AS comments, s.shares AS shares,
					u{id:ID(u), .username, .full_name, .avatar} AS owner,
					exists((me)-[:LIKE]->(s)) AS is_liked,
					CASE WHEN {userid} = {myuserid} THEN true ELSE false END AS can_edit,
					CASE WHEN {userid} = {myuserid} THEN true ELSE false END AS can_delete
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, params.Sort)
	} else if params.Type == configs.SPostStatus {
		stmt = fmt.Sprintf(`
		    MATCH(u:User) WHERE ID(u) = {userid}
				MATCH(me:User) WHERE ID(me) = {myuserid}
		  	MATCH (s:Status:Post)<-[r:POST]-(u)
				WHERE s.privacy = 1 OR (s.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR {userid} = {myuserid}
				RETURN
					ID(s) AS id,
					substring(s.message,0,250) AS message, length(s.message)>250 AS summary,
					s.created_at AS created_at, s.updated_at AS updated_at,
					case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
					s.likes AS likes, s.comments AS comments, s.shares AS shares,
					u{id:ID(u), .username, .full_name, .avatar} AS owner,
					exists((me)-[:LIKE]->(s)) AS is_liked,
					CASE WHEN {userid} = {myuserid} THEN true ELSE false END AS can_edit,
					CASE WHEN {userid} = {myuserid} THEN true ELSE false END AS can_delete
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, params.Sort)
	} else if params.Type == configs.SPost {
		stmt = fmt.Sprintf(`
		    MATCH(u:User) WHERE ID(u) = {userid}
				MATCH(me:User) WHERE ID(me) = {myuserid}
		  	MATCH (s:Post)<-[r:POST]-(u)
				WHERE s.privacy = 1 OR (s.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR {userid} = {myuserid}
				RETURN
					ID(s) AS id,
					substring(s.message,0,250) AS message, length(s.message)>250 AS summary,
					case s.photo when null then "" else s.photo end AS photo,
					s.created_at AS created_at, s.updated_at AS updated_at,
					case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
					s.likes AS likes, s.comments AS comments, s.shares AS shares,
					u{id:ID(u), .username, .full_name, .avatar} AS owner,
					exists((me)-[:LIKE]->(s)) AS is_liked,
					CASE WHEN {userid} = {myuserid} THEN true ELSE false END AS can_edit,
					CASE WHEN {userid} = {myuserid} THEN true ELSE false END AS can_delete
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, params.Sort)
	}
	paramsQuery := map[string]interface{}{
		"userid":   userID,
		"myuserid": myUserID,
		"skip":     params.Skip,
		"limit":    params.Limit,
	}
	res := []models.Post{}
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

// Get func to get a post
// int64 int64
// models.Post error
func (service postService) Get(postID int64, myUserID int64) (models.Post, error) {
	stmt := `
			MATCH(me:User) WHERE ID(me) = {myuserid}
			MATCH (s:Post)<-[:POST]-(u:User)
			WHERE ID(s) = {postid}
			RETURN
				ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
				case s.photo when null then "" else s.photo end AS photo,
				case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
				ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
				u{id:ID(u), .username, .full_name, .avatar} AS owner,
				s.likes AS likes, s.comments AS comments, s.shares AS shares,
				exists((me)-[:LIKE]->(s)) AS is_liked,
				CASE WHEN ID(u) = {myuserid} THEN true ELSE false END AS can_edit,
				CASE WHEN ID(u) = {myuserid} THEN true ELSE false END AS can_delete
			`
	params := map[string]interface{}{
		"postid":   postID,
		"myuserid": myUserID,
	}
	res := []models.Post{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Post{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Post{}, nil
}

// Delete func
// int64
// bool error
func (service postService) Delete(postID int64) (bool, error) {
	stmt := `
	  	MATCH (s:Post)
			WHERE ID(s) = {postid}
			DETACH DELETE s
	  	`
	params := map[string]interface{}{
		"postid": postID,
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

// Create func to create a new post
// models.Post int64
// int64 error
func (service postService) Create(post models.Post, myUserID int64) (int64, error) {

	var p interface{}
	var stmt string
	if len(post.Photo) == 0 {
		p = neoism.Props{
			"message":  post.Message,
			"privacy":  post.Privacy,
			"status":   post.Status,
			"likes":    0,
			"comments": 0,
			"shares":   0,
		}
		stmt = `
		    MATCH(u:User) WHERE ID(u) = {fromid}
		  	CREATE (s:Status:Post { props } )<-[r:POST]-(u)
				SET s.created_at = TIMESTAMP()
				RETURN ID(s) as id
		  	`
	} else {
		p = neoism.Props{
			"message":  post.Message,
			"photo":    post.Photo,
			"privacy":  post.Privacy,
			"status":   post.Status,
			"likes":    0,
			"comments": 0,
			"shares":   0,
		}
		stmt = `
		    MATCH(u:User) WHERE ID(u) = {fromid}
		  	CREATE (s:Photo:Post { props } )<-[r:POST]-(u)
				SET s.created_at = TIMESTAMP()
				RETURN ID(s) as id
		  	`
	}
	params := map[string]interface{}{
		"props":  p,
		"fromid": myUserID,
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
// models.Post
// models.Post error
func (service postService) Update(post models.Post) (models.Post, error) {
	var stmt string
	var params map[string]interface{}
	if len(post.Photo) > 0 {
		stmt = `
			MATCH (s:Post)<-[r:POST]-(u:User)
			WHERE ID(s) = {postid}
			SET s.message = {message}, s.photo = {photo}, s.privacy = {privacy}, s.updated_at = TIMESTAMP(), s.status = {status}, s:Photo
			RETURN
				ID(s) AS id, s.message AS message, s.photo AS photo, s.created_at AS created_at, s.updated_at AS updated_at,
				case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
				ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
				exists((u)-[:LIKE]->(s)) AS is_liked,
				s.likes AS likes, s.comments AS comments, s.shares AS shares,
				true AS can_edit,
				true AS can_delete
			`
		params = map[string]interface{}{
			"postid":  post.ID,
			"message": post.Message,
			"photo":   post.Photo,
			"privacy": post.Privacy,
			"status":  post.Status,
		}
	} else {
		stmt = `
  	MATCH (s:Post)<-[r:POST]-(u:User)
    WHERE ID(s) = {postid}
		SET s.message = {message}, s.privacy = {privacy}, s.updated_at = TIMESTAMP(), s.status = {status}
    RETURN
			ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
			case s.photo when null then "" else s.photo end AS photo,
			case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
			exists((u)-[:LIKE]->(s)) AS is_liked,
			s.likes AS likes, s.comments AS comments, s.shares AS shares,
			true AS can_edit,
			true AS can_delete
  	`
		params = map[string]interface{}{
			"postid":  post.ID,
			"message": post.Message,
			"privacy": post.Privacy,
			"status":  post.Status,
		}
	}

	res := []models.Post{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Post{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Post{}, errors.New("Dont' update user status")
}

// CheckExistPost func
// int64
// bool error
func (service postService) CheckExistPost(postID int64) (bool, error) {
	stmt := `
		MATCH (u:Post)
		WHERE ID(u)={postid}
		RETURN ID(u) AS id
		`
	params := neoism.Props{
		"postid": postID,
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
		if res[0].ID == postID {
			return true, nil
		}
	}
	return false, nil
}

// GetUserIDByPostID func
// int64
// int64 error
func (service postService) GetUserIDByPostID(postID int64) (int64, error) {
	stmt := `
	    MATCH (u:User)-[r:POST]->(s:Post)
			WHERE ID(s) = {postid}
			RETURN ID(u) AS id
	  	`
	params := map[string]interface{}{
		"postid": postID,
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

// IncreasePosts func
// int64
// bool error
func (service postService) IncreasePosts(userID int64) (bool, error) {
	stmt := `
	MATCH (u:User)
	WHERE ID(u)= {userid}
	SET u.posts = u.posts+1
	RETURN ID(u) AS id
	`
	params := neoism.Props{"userid": userID}
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
		if res[0].ID == userID {
			return true, nil
		}
	}
	return false, nil
}

// DecreasePosts func
// int64
// bool error
func (service postService) DecreasePosts(userID int64) (bool, error) {
	stmt := `
	MATCH (u:User)
	WHERE ID(u)= {userid}
	SET u.posts = u.posts-1
	RETURN ID(u) AS id
	`
	params := neoism.Props{
		"userid": userID,
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
		if res[0].ID == userID {
			return true, nil
		}
	}
	return false, nil
}

// CreateLike func
// int64 int64
// int error
func (service postService) CreateLike(postID int64, userID int64) (int, error) {
	stmt := `
		MATCH(u:User) WHERE ID(u) = {userid}
		MATCH(s:Post) WHERE ID(s) = {postid}
		MERGE(u)-[l:LIKE]->(s)
		ON CREATE SET l.created_at = TIMESTAMP()
		RETURN exists((u)-[l]->(s)) AS liked, s.likes AS likes
		`
	params := map[string]interface{}{
		"postid": postID,
		"userid": userID,
	}
	res := []struct {
		Liked bool `json:"liked"`
		Likes int  `json:"likes"`
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
	if len(res) > 0 && res[0].Liked == true {
		return res[0].Likes + 1, nil
	}
	return -1, nil
}

// GetLikes func
// int64 helpers.ParamsGetAll
// []models.UserLikedObject
func (service postService) GetLikes(postID int64, myUserID int64, params helpers.ParamsGetAll) ([]models.UserLikedObject, error) {
	stmt := fmt.Sprintf(`
	MATCH (me:User) WHERE ID(me) = {myuserid}
	MATCH (u:User)-[l:LIKE]->(s:Post)
	WHERE ID(s) = {postid}
	RETURN
		ID(u) AS id, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
		l.created_at as liked_at,
		exists((me)-[:FOLLOW]->(u)) AS is_followed
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	`, params.Sort)
	paramsQuery := map[string]interface{}{
		"postid":   postID,
		"myuserid": myUserID,
		"skip":     params.Skip,
		"limit":    params.Limit,
	}
	res := []models.UserLikedObject{}
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

// DeleteLike func
// int64 int64
// int error
func (service postService) DeleteLike(postID int64, userID int64) (int, error) {
	stmt := `
	MATCH (u:User)-[l:LIKE]->(s:Post)
	WHERE ID(s) = {postid} AND ID(u) = {userid}
	DELETE l
	RETURN s.likes AS likes
	`
	params := map[string]interface{}{
		"postid": postID,
		"userid": userID,
	}
	res := []struct {
		Likes int `json:"likes"`
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
		return res[0].Likes - 1, nil
	}
	return -1, nil
}

// CheckExistLike func
// int64 int64
// bool error
func (service postService) CheckExistLike(postID int64, userID int64) (bool, error) {
	stmt := `
  	MATCH (u:User)-[l:LIKE]->(s:Post)
		WHERE ID(u) = {userid} AND ID(s) = {postid}
		RETURN ID(l) as id
  	`
	params := neoism.Props{
		"postid": postID,
		"userid": userID,
	}
	res := []struct {
		ID int `json:"id"`
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

// IncreaseLikes func
// int64
// bool error
func (service postService) IncreaseLikes(postID int64) (bool, error) {
	stmt := `
	MATCH (s:Post)
	WHERE ID(s)= {postid}
	SET s.likes = s.likes+1
	RETURN ID(s) AS id
	`
	params := neoism.Props{
		"postid": postID,
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
		if res[0].ID == postID {
			return true, nil
		}
	}
	return false, nil
}

// DecreaseLikes func
// int64
// bool error
func (service postService) DecreaseLikes(postID int64) (bool, error) {
	stmt := `
	MATCH (s:Post)
	WHERE ID(s)= {postid}
	SET s.likes = s.likes-1
	RETURN ID(s) AS id
	`
	params := neoism.Props{
		"postid": postID,
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
		if res[0].ID == postID {
			return true, nil
		}
	}
	return false, nil
}

// CheckPostInteractivePermission func to check interactive permisson for user with a post
// int64 int64
// bool error
func (service postService) CheckPostInteractivePermission(postID int64, userID int64) (bool, error) {
	stmt := `
		MATCH (who:User) WHERE ID(who) = {userid}
		MATCH (u:User)-[r:POST]->(s:Post)
		WHERE ID(s) = {postid}
		RETURN exists((who)-[:FOLLOW]->(u)) AS followed, s.privacy AS privacy, who = u AS owner
		`
	params := map[string]interface{}{
		"userid": userID,
		"postid": postID,
	}
	res := []struct {
		Followed bool `json:"followed"`
		Privacy  int  `json:"privacy"`
		Owner    bool `json:"owner"`
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
		if res[0].Privacy == configs.Public || (res[0].Followed && res[0].Privacy == configs.ShareToFollowers || res[0].Owner) {
			return true, nil
		}
	}
	return false, nil

}

// CreateFollow func
// int64 int64
// int64 error
func (service postService) CreateFollow(postID int64, userID int64) (int64, error) {
	stmt := `
		MATCH(u:User) WHERE ID(u) = {userID}
		MATCH(p:Post) WHERE ID(p) = {postID}
		MERGE(u)-[f:FOLLOW]->(p)
		ON CREATE SET f.created_at = TIMESTAMP()
		RETURN ID(f) AS id
		`
	params := map[string]interface{}{
		"postID": postID,
		"userID": userID,
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

// DeleteFollow func
// int64 int64
// bool error
func (service postService) DeleteFollow(postID int64, userID int64) (bool, error) {
	stmt := `
		MATCH (u:User)-[f:FOLLOW]->(p:Post)
		WHERE ID(u) = {userID} AND ID(p) = {postID}
		DELETE f
		`
	params := map[string]interface{}{
		"postID": postID,
		"userID": userID,
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

// CheckExistFollow func
// int64 int64
// bool error
func (service postService) CheckExistFollow(postID int64, userID int64) (bool, error) {
	stmt := `
		MATCH (u:User)-[f:FOLLOW]->(p:Post)
		WHERE ID(u) = {userID} AND ID(p) = {postID}
		RETURN ID(f) as id
		`
	params := map[string]interface{}{
		"postID": postID,
		"userID": userID,
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

// GetUsers func
// helpers.ParamsGetAll
// models.PublicUsers error
func (service postService) GetUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error) {
	stmt := `
	MATCH (u:User)
	with u {id:ID(u), .*}  AS user

		SKIP {skip}
		LIMIT {limit}
		RETURN collect(user) AS users
		`
	p := map[string]interface{}{
		"skip":  params.Skip,
		"limit": params.Limit,
	}
	var res []struct {
		Users []models.UserFollowObject `json:"users"`
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

	if len(res) > 0 {
		return res[0].Users, nil
	}
	return nil, nil
}

// GetCanMentionedUsers func to get users who could mentioned in Comment
// int64 helpers.ParamsGetAll int64
// []models.UserFollowObject error
func (service postService) GetCanMentionedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error) {
	stmt := fmt.Sprintf(`
	MATCH (me:User)-[f:FOLLOW]-(u:User)
 	WHERE id(me) = {myUserID}
 	WITH ID(u) AS id,me
 	MATCH (u1:User)-[l:LIKE|:FOLLOW|:POST]->(p:Post) ,(u2)-[cr:WRITE]->(c:Comment)-[:AT]->(p)
 	WHERE ID(p)= {posID} and u2<>u1
 	WITH  collect(id)+collect(id(u1))+ collect(id(u2)) AS users, me
 	UNWIND users AS x
 	WITH DISTINCT x, me
 	MATCH (mention:User)
 	WHERE ID(mention) = x
 	WITH
 		mention{id:ID(mention),.username, .avatar, .full_name, is_followed: exists((me)-[:FOLLOW]->(mention)) } AS user,
    mention.created_at AS created_at, mention.username AS username, mention.full_name AS full_name, ID(mention) AS id
 	ORDER BY %s
 	SKIP {skip}
  LIMIT {limit}
 	RETURN  collect(user) AS users

		`,
		params.Sort)
	p := map[string]interface{}{
		"postID":   postID,
		"myUserID": myUserID,
		"skip":     params.Skip,
		"limit":    params.Limit,
	}
	var res []struct {
		Users []models.UserFollowObject `json:"users"`
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

	if len(res) > 0 {
		return res[0].Users, nil
	}
	return nil, nil
}

// GetMentionedUsers func
// int64 helpers.ParamsGetAll int64
// []models.UserFollowObject error
func (service postService) GetMentionedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error) {
	return nil, nil
}

// GetLikedUsers func
// int64 helpers.ParamsGetAll int64
// []models.UserFollowObject error
func (service postService) GetLikedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error) {
	return nil, nil
}

// GetCommentedUsers func
// int64 helpers.ParamsGetAll int64
// []models.UserFollowObject error
func (service postService) GetCommentedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error) {
	return nil, nil
}

// GetFollowedUsers func
// int64 helpers.ParamsGetAll int64
// []models.UserFollowObject error
func (service postService) GetFollowedUsers(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.UserFollowObject, error) {
	return nil, nil
}
