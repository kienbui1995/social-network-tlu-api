package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// SemesterServiceInterface include method list
type SemesterServiceInterface interface {
	// User
	GetAllOfStudent(params helpers.ParamsGetAll, studentCode string) ([]models.Semester, error)
	GetAllOfTeacher(params helpers.ParamsGetAll, teacherCode string) ([]models.Semester, error)
	// Get(semesterID int64) (models.Semester, error)
	// Delete(semesterID int64) (bool, error)
	// Create(semester models.Semester) (int64, error)
	// Update(semester models.Semester) (models.Semester, error)
	// CheckExistSemester(semesterID int64) (bool, error)

	// Only Admin
	//update from TLU
	UpdateFromTLU(year string) (bool, error)
	GetAll(params helpers.ParamsGetAll) ([]models.Semester, error)
}

// semesterService struct
type semesterService struct{}

// NewSemesterService to constructor
func NewSemesterService() SemesterServiceInterface {
	return semesterService{}
}

// GetAll func
// helpers.ParamsGetAll
// []models.Semester error
func (service semesterService) GetAll(params helpers.ParamsGetAll) ([]models.Semester, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		MATCH (se:Semester)
		WITH se{id:ID(se),.* } AS semester
		ORDER BY %s
		SKIP {skip}
		LIMIT {limit}
		RETURN collect(semester) AS semesters
		  	`, "semester."+params.Sort)

	paramsQuery := map[string]interface{}{
		"skip":  params.Skip,
		"limit": params.Limit,
	}
	res := []struct {
		Semesters []models.Semester `json:"semesters"`
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
		return res[0].Semesters, nil
	}
	return nil, nil
}

// Get func to get a post
// int64 int64
// models.Post error
// func (service semesterService) Get(semesterID int64) (models.Semester, error) {
// 	stmt := `
// 			MATCH(me:User) WHERE ID(me) = {myuserid}
// 			MATCH (s:Post)<-[:POST]-(u:User)
// 			WHERE ID(s) = {postid}
// 			RETURN
// 				ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
// 				case s.photo when null then "" else s.photo end AS photo,
// 				case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
// 				ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
// 				u{id:ID(u), .username, .full_name, .avatar} AS owner,
// 				s.likes AS likes, s.comments AS comments, s.shares AS shares,
// 				exists((me)-[:LIKE]->(s)) AS is_liked,
// 				exists((me)-[:FOLLOW]->(s)) AS is_following,
// 				CASE WHEN ID(u) = {myuserid} THEN true ELSE false END AS can_edit,
// 				CASE WHEN ID(u) = {myuserid} THEN true ELSE false END AS can_delete
// 			`
// 	params := map[string]interface{}{
//
// 		"semesterID": semesterID,
// 	}
// 	res := []models.Semester{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return models.Semester{}, err
// 	}
// 	if len(res) > 0 {
// 		return res[0], nil
// 	}
// 	return models.Semester{}, nil
// }

// // Delete func
// // int64
// // bool error
// func (service semesterService) Delete(postID int64) (bool, error) {
// 	stmt := `
// 	  	MATCH (s:Post)
// 			WHERE ID(s) = {postid}
// 			OPTIONAL MATCH (g:Group)-[h:HAS]->(s)
// 			SET g.posts= g.posts - 1
// 			WITH s
// 			OPTIONAL MATCH (u:User)-[p:POST]->(s)
// 			WHERE exists((:Group)-[:HAS]->(s))=false
// 			SET u.posts = u.posts - 1
// 			WITH s
// 			OPTIONAL MATCH (c:Comment)-->(s)
// 			DETACH DELETE c
// 			WITH s
// 			OPTIONAL MATCH (s)-[]->(n:Notification)
// 			DETACH DELETE n
// 			DETACH DELETE s
// 	  	`
// 	params := map[string]interface{}{
// 		"postid": postID,
// 	}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 	}
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

// // Create func to create a new post
// // models.Post int64
// // int64 error
// func (service semesterService) Create(post models.Post, myUserID int64) (int64, error) {
//
// 	var p interface{}
// 	var action int
//
// 	var stmt string
// 	if len(post.Photo) == 0 {
// 		action = configs.IActionPostStatus
// 		p = neoism.Props{
// 			"message":  post.Message,
// 			"privacy":  post.Privacy,
// 			"status":   post.Status,
// 			"likes":    0,
// 			"comments": 0,
// 			"shares":   0,
// 		}
// 		stmt = `
// 		    MATCH(u:User) WHERE ID(u) = {fromid}
// 		  	CREATE (s:Status:Post { props } )<-[r:POST]-(u)
// 				SET u.posts = u.posts + 1
// 				CREATE (u)-[f:FOLLOW]->(s)
// 				SET s.created_at = TIMESTAMP(), f.created_at = TIMESTAMP()
// 				WITH s,u
// 				MATCH(u1:User)-[:FOLLOW]->(u)
// 				CREATE (s)-[g:GENERATE]->(n:Notification)<-[:REGISTERED]-(u1)
// 				SET n.action = {action}, g.created_at = TIMESTAMP(), n.updated_at = TIMESTAMP()
// 				RETURN ID(s) as id
// 		  	`
// 	} else {
// 		action = configs.IActionPostPhoto
// 		p = neoism.Props{
// 			"message":  post.Message,
// 			"photo":    post.Photo,
// 			"privacy":  post.Privacy,
// 			"status":   post.Status,
// 			"likes":    0,
// 			"comments": 0,
// 			"shares":   0,
// 		}
// 		stmt = `
// 		    MATCH(u:User) WHERE ID(u) = {fromid}
// 		  	CREATE (s:Photo:Post { props } )<-[r:POST]-(u)
// 				SET u.posts = u.posts + 1
// 				CREATE (u)-[f:FOLLOW]->(s)
// 				SET s.created_at = TIMESTAMP(), f.created_at = TIMESTAMP()
// 				WITH s,u
// 				MATCH(u1:User)-[:FOLLOW]->(u)
// 				CREATE (s)-[g:GENERATE]->(n:Notification)<-[:REGISTERED]-(u1)
// 				SET n.action = {action}, g.created_at = TIMESTAMP(), n.updated_at = TIMESTAMP()
// 				RETURN ID(s) as id
// 		  	`
// 	}
// 	params := map[string]interface{}{
// 		"props":  p,
// 		"fromid": myUserID,
// 		"action": action,
// 	}
// 	res := []struct {
// 		ID int64 `json:"id"`
// 	}{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return -1, err
// 	}
// 	if len(res) > 0 {
// 		return res[0].ID, nil
// 	}
// 	return -1, nil
// }

// // Update func
// // models.Post
// // models.Post error
// func (service semesterService) Update(post models.Post) (models.Post, error) {
// 	var stmt string
// 	var params map[string]interface{}
// 	if len(post.Photo) > 0 {
// 		stmt = `
// 			MATCH (s:Post)<-[r:POST]-(u:User)
// 			WHERE ID(s) = {postid}
// 			SET s.message = {message}, s.photo = {photo}, s.privacy = {privacy}, s.updated_at = TIMESTAMP(), s.status = {status}, s:Photo
// 			RETURN
// 				ID(s) AS id, s.message AS message, s.photo AS photo, s.created_at AS created_at, s.updated_at AS updated_at,
// 				case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
// 				ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
// 				exists((u)-[:LIKE]->(s)) AS is_liked,
// 				exists((u)-[:FOLLOW]->(s)) AS is_following,
// 				s.likes AS likes, s.comments AS comments, s.shares AS shares,
// 				true AS can_edit,
// 				true AS can_delete
// 			`
// 		params = map[string]interface{}{
// 			"postid":  post.ID,
// 			"message": post.Message,
// 			"photo":   post.Photo,
// 			"privacy": post.Privacy,
// 			"status":  post.Status,
// 		}
// 	} else {
// 		stmt = `
//   	MATCH (s:Post)<-[r:POST]-(u:User)
//     WHERE ID(s) = {postid}
// 		SET s.message = {message}, s.privacy = {privacy}, s.updated_at = TIMESTAMP(), s.status = {status}
//     RETURN
// 			ID(s) AS id, s.message AS message, s.created_at AS created_at, s.updated_at AS updated_at,
// 			case s.photo when null then "" else s.photo end AS photo,
// 			case s.privacy when null then 1 else s.privacy end AS privacy, case s.status when null then 1 else s.status end AS status,
// 			ID(u) AS userid, u.username AS username, u.full_name AS full_name, u.avatar AS avatar,
// 			exists((u)-[:LIKE]->(s)) AS is_liked,
// 			exists((u)-[:FOLLOW]->(s)) AS is_following,
// 			s.likes AS likes, s.comments AS comments, s.shares AS shares,
// 			true AS can_edit,
// 			true AS can_delete
//   	`
// 		params = map[string]interface{}{
// 			"postid":  post.ID,
// 			"message": post.Message,
// 			"privacy": post.Privacy,
// 			"status":  post.Status,
// 		}
// 	}
//
// 	res := []models.Post{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return models.Post{}, err
// 	}
// 	if len(res) > 0 {
// 		return res[0], nil
// 	}
// 	return models.Post{}, errors.New("Dont' update user status")
// }

// // CheckExistSemester func
// // int64
// // bool error
// func (service semesterService) CheckExistSemester(semesterID int64) (bool, error) {
// 	stmt := `
// 		MATCH (s:Semester)
// 		WHERE ID(s)={semesterID}
// 		RETURN ID(s) AS id
// 		`
// 	params := neoism.Props{
// 		"semesterID": semesterID,
// 	}
//
// 	res := []struct {
// 		ID int64 `json:"id"`
// 	}{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
//
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return false, err
// 	}
//
// 	if len(res) > 0 {
// 		if res[0].ID == semesterID {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// UpdateFromTLU func
// models.Post
// models.Post error
func (service semesterService) UpdateFromTLU(year string) (bool, error) {

	stmt := fmt.Sprintf(`
    CALL apoc.load.json("%s") YIELD value AS d
    UNWIND d.data AS hocki
    MERGE (s:Semester{code:toString(hocki.Ma)})
    ON CREATE SET
    	s.code =toString(hocki.Ma),
      s.year = toString(hocki.Nam),
      s.group = hocki.Nhom,
      s.symbol = toString(hocki.Kyhieu),
      s.start_at = toString(hocki.Thoigianbd),
      s.finish_at=toString(hocki.Thoigiankt),
      s.name=toString(hocki.Tenky),
				s.status =1,
      s.created_at= timestamp()
			`, configs.SURLGetSemesterListByYear+year)
	// params := map[string]interface{}{
	// 	"url": " + configs.SURLGetSemesterListByYear + year + "\"",
	// }
	//

	cq := neoism.CypherQuery{
		Statement: stmt,
		// // Parameters: params,
		// Result: &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	return true, nil

	// return false, errors.New("Dont' update semester")
}

// GetAllOfStudent func
// helpers.ParamsGetAll
// models.Post error
func (service semesterService) GetAllOfStudent(params helpers.ParamsGetAll, studentCode string) ([]models.Semester, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		MATCH(s:Student) WHERE toLower(s.code) = toLower({studentCode})
		MATCH (se:Semester) WHERE exists((se)--(:Class)--(s))
		WITH se{id:ID(se),.* } AS semester
		ORDER BY %s
		SKIP {skip}
		LIMIT {limit}
		RETURN collect(semester) AS semesters
		  	`, "semester."+params.Sort)

	paramsQuery := map[string]interface{}{
		"studentCode": studentCode,
		"skip":        params.Skip,
		"limit":       params.Limit,
	}
	res := []struct {
		Semesters []models.Semester `json:"semesters"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}
	fmt.Printf("cq: %v\n", cq)
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		return res[0].Semesters, nil
	}
	return nil, nil
}

// GetAllOfStudent func
// helpers.ParamsGetAll
// models.Post error
func (service semesterService) GetAllOfTeacher(params helpers.ParamsGetAll, teacherCode string) ([]models.Semester, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		MATCH(t:Teacher) WHERE toLower(t.code) = toLower({teacherCode})
		MATCH (se:Semester) WHERE exists((se)--(:Class)--(t))
		WITH se{id:ID(se),.* } AS semester
		ORDER BY %s
		SKIP {skip}
		LIMIT {limit}
		RETURN collect(semester) AS semesters
		  	`, "semester."+params.Sort)

	paramsQuery := map[string]interface{}{
		"teacherCode": teacherCode,
		"skip":        params.Skip,
		"limit":       params.Limit,
	}
	res := []struct {
		Semesters []models.Semester `json:"semesters"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}
	fmt.Printf("cq: %v\n", cq)
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		return res[0].Semesters, nil
	}
	return nil, nil
}
