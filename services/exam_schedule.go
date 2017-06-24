package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// ExamScheduleServiceInterface include method list
type ExamScheduleServiceInterface interface {
	// Admin
	GetAll(params helpers.ParamsGetAll) ([]models.ExamSchedule, error)
	// Get(semesterID int64) (models.ExamSchedule, error)
	// Delete(semesterID int64) (bool, error)
	// Create(semester models.ExamSchedule) (int64, error)
	// Update(semester models.ExamSchedule) (models.ExamSchedule, error)
	// CheckExistSubject(subjectID int64) (bool, error)

	//update from TLU
	UpdateFromTLU(semesterCode string) (bool, error)

	// A Student
	GetAllByStudent(params helpers.ParamsGetAll, semesterCode string, studentCode string) ([]models.ExamSchedule, error)
}

// examScheduleService struct
type examScheduleService struct{}

// NewExamScheduleService to constructor
func NewExamScheduleService() ExamScheduleServiceInterface {
	return examScheduleService{}
}

// GetAll func
// helpers.ParamsGetAll
// models.Post error
func (service examScheduleService) GetAll(params helpers.ParamsGetAll) ([]models.ExamSchedule, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		    MATCH(s:Subject)
				with s
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
				RETURN
					collect(s{id:ID(s),.*}) AS subject


		  	`, "s."+params.Sort)

	paramsQuery := map[string]interface{}{
		"skip":  params.Skip,
		"limit": params.Limit,
	}
	res := []struct {
		ExamSchedules []models.ExamSchedule `json:"subject"`
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
		return res[0].ExamSchedules, nil
	}
	return nil, nil
}

// GetAllByStudent func
// helpers.ParamsGetAll string
// []models.Class error
func (service examScheduleService) GetAllByStudent(params helpers.ParamsGetAll, semesterCode string, studentCode string) ([]models.ExamSchedule, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		MATCH (s:Student)-[:ATTEND]->(exam:ExamSchedule)<-[:ORGANIZE]-(semester:Semester),
					(exam)--(sub:Subject),
					(exam)--(r:Room)
		WHERE toLower(s.code) = toLower({studentCode}) AND semester.code = {semesterCode}
		WITH
		exam{
		id: ID(exam),.*,
		subject: sub{id:ID(sub),.code,.name},
		room: r{id:ID(r),.code}
		} AS exam_schedule
		ORDER BY %s
		SKIP {skip}
		LIMiT {limit}
		return collect(exam_schedule) AS exam_schedules
		  	`, "exam_schedule."+params.Sort)

	paramsQuery := map[string]interface{}{
		"skip":         params.Skip,
		"limit":        params.Limit,
		"studentCode":  studentCode,
		"semesterCode": semesterCode,
	}
	res := []struct {
		ExamSchedules []models.ExamSchedule `json:"exam_schedules"`
	}{}
	// res := []interface{}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}
	// fmt.Printf("cq: %v\n", cq)
	err := conn.Cypher(&cq)
	if err != nil {
		return nil, err
	}
	if len(res) > 0 {
		// fmt.Printf("res: %v\n", res)
		return res[0].ExamSchedules, nil
	}
	return nil, nil
}

// Get func to get a post
// int64 int64
// models.Post error
// func (service subjectService) Get(semesterCode string) (models.Subject, error) {
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
// 	res := []models.ExamSchedule{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: params,
// 		Result:     &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return models.ExamSchedule{}, err
// 	}
// 	if len(res) > 0 {
// 		return res[0], nil
// 	}
// 	return models.ExamSchedule{}, nil
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

// // CheckExistPost func
// // int64
// // bool error
// func (service subjectService) CheckExistSubject(subjectID int64) (bool, error) {
// 	stmt := `
// 		MATCH (s:Subject)
// 		WHERE ID(s)={subjectID}
// 		RETURN ID(s) AS id
// 		`
// 	params := neoism.Props{
// 		"subjectID": subjectID,
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
// 		if res[0].ID == subjectID {
// 			return true, nil
// 		}
// 	}
// 	return false, nil
// }

// UpdateFromTLU func
// models.Post
// models.Post error
func (service examScheduleService) UpdateFromTLU(semesterCode string) (bool, error) {
	stmt := fmt.Sprintf(`
		CALL apoc.load.json("%s") YIELD value AS d
      UNWIND d.data AS schedule
      MATCH(semester:Semester{code:"%s"})
      MATCH(s:Student{code:toString(schedule.masv)})
      MATCH(sub:Subject{code:toString(schedule.mahp)})
			WITH CASE schedule.maph
						WHEN NULL THEN ""
						ELSE toString(schedule.maph)
			END AS maph, schedule
      MERGE(r:Room{code:maph})
      MERGE(r)<-[:EXAM_AT]-(exam:ExamSchedule{day:toString(schedule.ngay),exam_time:toString(schedule.giothi)})-[:EXAM_FOR]->(sub)
      ON CREATE SET
        exam.status =1,
        exam.created_at = timestamp()
      MERGE (s)-[:ATTEND]->(exam)
      MERGE (semester)-[:ORGANIZE]->(exam)
			`, configs.SURLGetExamScheduleListBySemesterCode+semesterCode, semesterCode)
	// params := map[string]interface{}{
	// 	"url": " + configs.SURLGetSemesterListByYear + year + "\"",
	// }
	//

	cq := neoism.CypherQuery{
		Statement: stmt,
		// // Parameters: params,
		// Result: &res,
	}
	fmt.Printf("semester: %v\n", semesterCode)
	err := conn.Cypher(&cq)
	if err != nil {
		return false, err
	}

	return true, nil

	// return false, errors.New("Dont' update semester")
}
