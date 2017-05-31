package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// HomeServiceInterface include method list
type HomeServiceInterface interface {
	FindUserByUsernameAndFullName(name string, myUserID int64) ([]models.UserFollowObject, error)
	GetNewsFeed(params helpers.ParamsGetAll, myUserID int64) ([]models.Post, error)
	GetNewsFeedWithPageRank(params helpers.ParamsGetAll, myUserID int64) ([]models.Post, error)
}

// homeService struct
type homeService struct{}

// NewHomeService to constructor
func NewHomeService() homeService {
	return homeService{}
}

// FindUserByUsernameAndFullName func
// string int64
// []models.UserFollowObject error
func (service homeService) FindUserByUsernameAndFullName(name string, myUserID int64) ([]models.UserFollowObject, error) {
	stmt := `
		 MATCH(a:User) where ID(a)={userid}
		 OPTIONAL MATCH(u:User)
		 WHERE toLower(u.username) CONTAINS toLower({s})  OR toLower(u.full_name)  CONTAINS toLower({s})
		 RETURN
		 	ID(u) as id, u.username as username,
			u.avatar as avatar,
			u.full_name as full_name,
		  exists((a)-[:FOLLOW]->(u)) as is_followed
	`
	res := []models.UserFollowObject{}
	params := neoism.Props{
		"userid": myUserID,
		"s":      name,
	}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
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

//GetNewsFeed func
// helpers.ParamsGetAll int64
// []models.Post error
func (service homeService) GetNewsFeed(params helpers.ParamsGetAll, myUserID int64) ([]models.Post, error) {
	stmt := fmt.Sprintf(`
		MATCH (u:User)-[:POST]->(p:Post)<-[:HAS]-(g:Group)<-[j:JOIN{status:1}]-(me:User)
		WHERE ID(me) = {myUserID} AND 3>=j.role>=1
		WITH
		collect(
		p{
		id:ID(p),
		.*,
		message: substring(p.message,0,250),
		summary: size(p.message)>250,
		owner: u{id:ID(u), .username, .full_name, .avatar},
		place: g{id:ID(g), .name, .avatar},
		is_liked: exists((me)-[:LIKE]->(p)),
		is_following: exists((me)-[:FOLLOW]->(p)),
		can_edit: CASE WHEN ID(u) = ID(me) THEN true ELSE false END,
		can_delete: CASE WHEN ID(u) = ID(me) OR exists((me)-[:JOIN{role:2}]->(g))OR exists((me)-[:JOIN{role:3}]->(g)) THEN true ELSE false END
		}) AS posts1, me

		MATCH(me)-[:FOLLOW]->(u:User)-[:POST]->(p:Post)
		WHERE (p.privacy = 1 OR (p.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR u = me) AND exists((:Group)-[:HAS]->(p))=false
		WITH
		collect(
		p{
		id:ID(p), .*,
		message: substring(p.message,0,250),
		summary: size(p.message)>250,
		owner: u{id:ID(u), .username, .full_name, .avatar},
		is_liked: exists((me)-[:LIKE]->(p)),
		is_following: exists((me)-[:FOLLOW]->(p)),
		can_edit: CASE WHEN ID(u) = ID(me) THEN true ELSE false END,
		can_delete:	CASE WHEN ID(u) = ID(me) THEN true ELSE false END
		}) AS posts2, posts1
		WITH  posts1+posts2 AS posts
		UNWIND posts AS p
		WITH p ORDER BY %s
		RETURN collect(p{
			.*
		} ) as post
	`, "p."+params.Sort)
	res := []struct {
		Post []models.Post `json:"post"`
	}{}
	paramsQuery := map[string]interface{}{
		"myUserID": myUserID,
		"skip":     params.Skip,
		"limit":    params.Limit,
	}
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
		return res[0].Post[params.Skip : params.Limit-params.Skip], nil
	}
	return nil, nil
}

//GetNewsFeedWithPageRank func
// helpers.ParamsGetAll int64
// []models.Post error
func (service homeService) GetNewsFeedWithPageRank(params helpers.ParamsGetAll, myUserID int64) ([]models.Post, error) {
	stmt := `
	MATCH (u:User)-[:POST]->(p:Post)<-[:HAS]-(g:Group)<-[j:JOIN{status:1}]-(me:User)
	WHERE ID(me) = {myUserID} AND 3>=j.role>=1
	WITH
	collect(
	p{
	id:ID(p),
	.*,
	message: substring(p.message,0,250),
	summary: size(p.message)>250,
	owner: u{id:ID(u), .username, .full_name, .avatar},
	place: g{id:ID(g), .name, .avatar},
	is_liked: exists((me)-[:LIKE]->(p)),
	is_following: exists((me)-[:FOLLOW]->(p)),
	can_edit: CASE WHEN ID(u) = ID(me) THEN true ELSE false END,
	can_delete: CASE WHEN ID(u) = ID(me) OR exists((me)-[:JOIN{role:2}]->(g))OR exists((me)-[:JOIN{role:3}]->(g)) THEN true ELSE false END
	}) AS posts1, me

	MATCH(me)-[:FOLLOW]->(u:User)-[:POST]->(p:Post)
	WHERE (p.privacy = 1 OR (p.privacy = 2 AND exists((me)-[:FOLLOW]->(u))) OR u = me) AND exists((:Group)-[:HAS]->(p))=false
	WITH
	collect(
	p{
	id:ID(p), .*,
	message: substring(p.message,0,250),
	summary: size(p.message)>250,
	owner: u{id:ID(u), .username, .full_name, .avatar},
	is_liked: exists((me)-[:LIKE]->(p)),
	is_following: exists((me)-[:FOLLOW]->(p)),
	can_edit: CASE WHEN ID(u) = ID(me) THEN true ELSE false END,
	can_delete:	CASE WHEN ID(u) = ID(me) THEN true ELSE false END
	}) AS posts2, posts1
	WITH  posts1+posts2 AS posts
	UNWIND posts AS p
	return collect(p{
		.*
	} ) AS post
	`
	res := []struct {
		Post []models.Post `json:"post"`
	}{}
	paramsQuery := map[string]interface{}{
		"myUserID": myUserID,
		"skip":     params.Skip,
		"limit":    params.Limit,
	}
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
		return res[0].Post[params.Skip : params.Limit-params.Skip], nil
	}
	return nil, nil
}
