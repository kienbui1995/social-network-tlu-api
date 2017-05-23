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
		MATCH(u:User) WHERE ID(u)= {userid}
		MATCH(u)-[:FOLLOW]->(u1:User)-[:POST]->(s:Post)
		WHERE s.privacy = 1 OR (s.privacy = 2 AND exists((u)-[:FOLLOW]->(u1))) OR u1 = u
		RETURN
			ID(s) AS id, s.message AS message, s.created_at AS created_at,
		  case s.uploaded_at when null then "" else s.uploaded_at end AS uploaded_at,
			case s.photo when null then "" else s.photo end AS photo,
			case s.privacy when null then 1 else s.privacy end AS privacy,
			case s.status when null then 1 else s.status end AS status,
			u1{id:ID(u1), .username, .full_name, . avatar} as owner,
			s.likes AS likes, s.comments AS comments, s.shares AS shares,
			exists((u)-[:LIKE]->(s)) AS is_liked,
			CASE WHEN ID(u1) = {userid} THEN true ELSE false END AS can_edit,
			CASE WHEN ID(u1) = {userid} THEN true ELSE false END AS can_delete
	ORDER BY %s
	SKIP {skip}
	LIMIT {limit}
	`, params.Sort)
	res := []models.Post{}
	paramsQuery := map[string]interface{}{
		"userid": myUserID,
		"skip":   params.Skip,
		"limit":  params.Limit,
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
		if res[0].ID >= 0 {
			return res, nil
		}
	}
	return nil, nil
}

//GetNewsFeedWithPageRank func
// helpers.ParamsGetAll int64
// []models.Post error
func (service homeService) GetNewsFeedWithPageRank(params helpers.ParamsGetAll, myUserID int64) ([]models.Post, error) {
	stmt := `
	MATCH(u:User) WHERE ID(u)= {myUserID}
	MATCH(u)-[:FOLLOW]->(u1:User)-[:POST]->(p:Post)
	WHERE p.privacy = 1 OR (p.privacy = 2 AND exists((u)-[:FOLLOW]->(u1))) OR u1 = u
		WITH u, u1, COLLECT(p) as posts
  CALL apoc.algo.pageRank(posts) YIELD node AS s, score
  RETURN
		ID(s) AS id, s.message AS message, s.created_at AS created_at,
		case s.updated_at when null then "" else s.updated_at end AS updated_at,
		case s.photo when null then "" else s.photo end AS photo,
		case s.privacy when null then 1 else s.privacy end AS privacy,
		case s.status when null then 1 else s.status end AS status,
		u1{id:ID(u1), .username, .full_name, .avatar} AS owner,
		s.likes AS likes, s.comments AS comments, s.shares AS shares,
		exists((u)-[:LIKE]->(s)) AS is_liked,
					exists ((u)-[:FOLLOW]->(s)) AS is_followed,
		CASE WHEN ID(u1) = {myUserID} THEN true ELSE false END AS can_edit,
		CASE WHEN ID(u1) = {myUserID} THEN true ELSE false END AS can_delete
	ORDER BY score*TIMESTAMP()/((TIMESTAMP()- created_at)/10+1) DESC
	SKIP {skip}
	LIMIT {limit}
	`
	res := []models.Post{}
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
		if res[0].ID >= 0 {
			return res, nil
		}
	}
	return nil, nil
}
