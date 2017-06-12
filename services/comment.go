package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// CommentServiceInterface include method list
type CommentServiceInterface interface {
	GetAll(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.Comment, error)
	Get(commentID int64) (models.Comment, error)
	Create(comment models.Comment, postID int64) (int64, error)
	CreateWithMention(comment models.Comment, postID int64) (int64, error)
	Delete(commentID int64) (bool, error)
	Update(comment models.Comment) (bool, error)
	CheckExistComment(commentID int64) (bool, error)
	GetUserIDByComment(commentID int64) (int64, error)
	GetPostIDbyComment(commentID int64) (int64, error)
	CheckPostInteractivePermission(postID int64, userID int64) (bool, error)
}

// commentService struct
type commentService struct {
}

// NewCommentService to constructor
func NewCommentService() commentService {
	return commentService{}
}

// GetAll func
// int64 helpers.ParamsGetAll int64
// []models.Comment error
func (service commentService) GetAll(postID int64, params helpers.ParamsGetAll, myUserID int64) ([]models.Comment, error) {
	stmt := fmt.Sprintf(`
	MATCH (me:User) WHERE ID(me) = {userid}
	MATCH (u:User)-[w:WRITE]->(c:Comment)-[a:AT]->(s:Post)
	WHERE ID(s) = {postid}
	RETURN
		ID(c) AS id, c.message AS message, c.created_at AS created_at, c.updated_at AS updated_at ,c.status AS status,
		u{id:ID(u),.username, .full_name, .avatar} AS owner, c.mentions,
		ID(u) <> ID(me) AS can_report,
		ID(u) = ID(me) AS can_edit,
		ID(u) = ID(me) AS can_delete
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	`, params.Sort)
	paramsQuery := map[string]interface{}{
		"postid": postID,
		"skip":   params.Skip,
		"limit":  params.Limit,
		"userid": myUserID,
	}

	res := []models.Comment{}
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

// Get func
// int64
// models.Comment error
func (service commentService) Get(commentID int64) (models.Comment, error) {
	stmt := `
	MATCH (c:Comment)<-[:WRITE]-(u:User)
	WHERE ID(c) = {commentID}
	RETURN
		ID(c) AS id,
		c.message AS message,
		c.created_at AS created_at,
		c.updated_at AS updated_at ,
		c.status AS status,
		u{id:ID(u),.username, .full_name, .avatar} AS owner,
		apoc.convert.fromJsonList(c.mentions) AS mentions
	`
	params := map[string]interface{}{
		"commentID": commentID,
	}

	res := []models.Comment{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Comment{}, err
	}
	if len(res) > 0 {
		if res[0].ID >= 0 {
			return res[0], nil
		}
	}
	return models.Comment{}, nil
}

// Create func
// models.Comment int64
// int64 error
func (service commentService) Create(comment models.Comment, postID int64) (int64, error) {

	p := neoism.Props{
		"message": comment.Message,
		"status":  comment.Status,
	}
	params := map[string]interface{}{
		"props":    p,
		"userid":   comment.Owner.ID,
		"postid":   postID,
		"mentions": comment.Mentions,
	}
	// if comment.Mentions != nil {
	// 	var ids []int64
	// 	for index := 0; index < len(comment.Mentions); index++ {
	// 		ids = append(ids, comment.Mentions[index].ID)
	// 	}
	// }

	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid}
	MATCH (s:Post) WHERE ID(s) = {postid}
	CREATE (c:Comment { props } ) SET c.created_at = TIMESTAMP(),c.mentions=apoc.convert.toJson({mentions})
	CREATE (u)-[w:WRITE]->(c)-[a:AT]->(s)
	SET s.comments = s.comments+1
	RETURN ID(c) AS id
	`

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

// Create func
// models.Comment int64
// int64 error
func (service commentService) CreateWithMention(comment models.Comment, postID int64) (int64, error) {

	p := neoism.Props{
		"message": comment.Message,
		"status":  comment.Status,
	}
	params := map[string]interface{}{
		"props":    p,
		"userid":   comment.Owner.ID,
		"postid":   postID,
		"mentions": comment.Mentions,
	}
	stmt := `
	MATCH (u:User) WHERE ID(u) = {userid}
	MATCH (s:Post) WHERE ID(s) = {postid}
	CREATE (c:Comment { props } ) SET c.created_at = TIMESTAMP(), c.mentions=apoc.convert.toJson({mentions})
	CREATE (u)-[w:WRITE]->(c)-[a:AT]->(s)
	SET s.comments = s.comments+1
	RETURN ID(c) AS id
	`

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

// Delete func
// int64
// bool error
func (service commentService) Delete(commentID int64) (bool, error) {
	stmt := `
	MATCH (c:Comment)
	WHERE ID(c) = {commentID}
	OPTIONAL MATCH(c)-[]->(p:Post)
	SET p.comments = p.comments-1
	DETACH DELETE c
	`
	params := map[string]interface{}{
		"commentID": commentID,
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

// Update func
// models.comment
// bool error
func (service commentService) Update(comment models.Comment) (bool, error) {
	stmt := `
	MATCH (c:Comment)
	WHERE ID(c) = {commentid}
	SET c.message = {message}, c.updated_at = TIMESTAMP()
  RETURN ID(c) AS id
	`
	params := map[string]interface{}{
		"commentid": comment.ID,
		"message":   comment.Message,
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
		if res[0].ID == comment.ID {
			return true, nil
		}
	}
	return false, nil
}

// CheckExistComment func
// int64
// bool error
func (service commentService) CheckExistComment(commentID int64) (bool, error) {
	stmt := `
		MATCH (c:Comment)
		WHERE ID(c)={commentID}
		RETURN ID(c) AS id
		`
	params := neoism.Props{
		"commentID": commentID,
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
		if res[0].ID == commentID {
			return true, nil
		}
	}
	return false, nil
}

// GetUserIDByComment func
// int64
// int64 error
func (service commentService) GetUserIDByComment(commentID int64) (int64, error) {
	stmt := `
	    MATCH (u:User)-[w:WRITE]->(c:Comment)
			WHERE ID(c) = {commentID}
			RETURN ID(u) AS id
	  	`
	params := map[string]interface{}{
		"commentID": commentID,
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

// GetPostIDbyComment func
// int64
// int64 error
func (service commentService) GetPostIDbyComment(commentID int64) (int64, error) {
	stmt := `
    MATCH (c:Comment)-[:AT]->(s)
		WHERE ID(c) = {commentID}
		RETURN ID(s) AS id
  	`
	params := map[string]interface{}{
		"commentID": commentID,
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

// CheckPostInteractivePermission func to check interactive permisson for user with a post
// int64 int64
// bool error
func (service commentService) CheckPostInteractivePermission(postID int64, userID int64) (bool, error) {
	stmt := `
		MATCH (who:User)
		WHERE ID(who) = {userID}
		MATCH (u:User)-[r:POST]->(p:Post)
		WHERE ID(p) = {postID}
		RETURN
			exists((who)-[:FOLLOW]->(u)) AS followed,
			p.privacy AS privacy,
			CASE WHEN exists((who)-[:JOIN]->(:Group)-[:HAS]->(p)) THEN true END AS see_in_group,
			who = u AS owner
		`
	params := map[string]interface{}{"userID": userID, "postID": postID}
	res := []struct {
		Followed   bool `json:"followed"`
		Privacy    int  `json:"privacy"`
		Owner      bool `json:"owner"`
		SeeInGroup bool `json:"see_in_group,omitempty"`
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
		if res[0].Privacy == configs.Public || (res[0].Followed && res[0].Privacy == configs.ShareToFollowers || res[0].Owner) || res[0].SeeInGroup {
			return true, nil
		}
	}
	return false, nil
}
