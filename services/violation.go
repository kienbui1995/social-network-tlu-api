package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// ViolationServiceInterface include method list
type ViolationServiceInterface interface {
	GetAll(params helpers.ParamsGetAll) ([]models.Violation, error)
	GetAllOfSupervisior(params helpers.ParamsGetAll, supervisiorCode string) ([]models.Violation, error)
	GetAllOfStudent(params helpers.ParamsGetAll, studentCode string) ([]models.Violation, error)
	//GetAllOfSemester(params helpers.ParamsGetAll, semesterCode string) ([]models.Violation, error)
	Delete(violationID int64) (bool, error)
	Create(violation models.Violation, studentCode string, supervisiorID int64) (int64, error)
	Update(violation models.Violation) (models.Violation, error)
	CheckExistViolation(violationID int64) (bool, error)
	Get(violationID int64) (models.Violation, error)
	CheckPermission(violationID int64, myUserID int64) (bool, error)
	GetSupervisiorID(myUserID int64) (int64, error)
}

// violationService struct
type violationService struct{}

// NewViolationService to constructor
func NewViolationService() ViolationServiceInterface {
	return violationService{}
}

// GetAll func
// helpers.ParamsGetAll
// []models.Violation error
func (service violationService) GetAll(params helpers.ParamsGetAll) ([]models.Violation, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		  	MATCH (s:Student)<-(v:Violation)<-[:CATCH]-(sv:Supervisior)
				RETURN
					ID(v) AS id,
					v.message AS message,
					v.place AS place,
					v.time_at AS time_at,
					CASE v.photo when null then "" else v.photo end AS photo,
					v.created_at AS created_at, v.updated_at AS updated_at,
					s{id:ID(s),name: s.last_name+" "+s.first_name, .code} AS owner,
					sv{id:ID(sv),.code,name: sv.last_name+" "+sv.first_name} AS catcher
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, params.Sort)

	paramsQuery := map[string]interface{}{
		"skip":  params.Skip,
		"limit": params.Limit,
	}
	res := []models.Violation{}
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

// GetAll func
// helpers.ParamsGetAll string
// []models.Violation error
func (service violationService) GetAllOfSupervisior(params helpers.ParamsGetAll, supervisiorCode string) ([]models.Violation, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		  	MATCH (s:Student)<-[:OF]-(v:Violation)<-[:CATCH]-(sv:Supervisior)
				WHERE toLower(sv.code) CONTAINS toLower({supervisiorCode})
				RETURN
					ID(v) AS id,
					v.message AS message,
					v.place AS place,
					v.time_at AS time_at,
					CASE v.photo when null then "" else v.photo end AS photo,
					v.created_at AS created_at, v.updated_at AS updated_at,
					s{id:ID(s),name: s.last_name+" "+s.first_name, .code} AS owner,
					sv{id:ID(sv),.code,name: sv.last_name+" "+sv.first_name} AS catcher
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, params.Sort)

	paramsQuery := map[string]interface{}{
		"skip":            params.Skip,
		"limit":           params.Limit,
		"supervisiorCode": supervisiorCode,
	}
	res := []models.Violation{}
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

// GetAll func
// helpers.ParamsGetAll string
// []models.Violation error
func (service violationService) GetAllOfStudent(params helpers.ParamsGetAll, studentCode string) ([]models.Violation, error) {
	var stmt string
	stmt = fmt.Sprintf(`
		  	MATCH (s:Student)<-[:OF]-(v:Violation)<-[:CATCH]-(sv:Supervisior)
				WHERE toLower(s.code) CONTAINS toLower({studentCode})
				RETURN
					ID(v) AS id,
					v.message AS message,
					v.place AS place,
					v.time_at AS time_at,
					CASE v.photo when null then "" else v.photo end AS photo,
					v.created_at AS created_at, v.updated_at AS updated_at,
					s{id:ID(s),name: s.last_name+" "+s.first_name, .code} AS owner,
					sv{id:ID(sv),.code,name: sv.last_name+" "+sv.first_name} AS catcher
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, params.Sort)

	paramsQuery := map[string]interface{}{
		"skip":        params.Skip,
		"limit":       params.Limit,
		"studentCode": studentCode,
	}
	res := []models.Violation{}
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

// // GetAll func
// // helpers.ParamsGetAll string
// // []models.Violation error
// func (service violationService) GetAllOfSemester(params helpers.ParamsGetAll, semesterCode string) ([]models.Violation, error) {
// 	var stmt string
// 	stmt = fmt.Sprintf(`
// 		  	MATCH (s:Student)<-(v:Violation)<-[:CATCH]-(sv:Supervisior)
// 				WHERE toLower(s.code) CONTAINS toLower({studentCode})
// 				RETURN
// 					ID(v) AS id,
// 					v.message AS message,
// 					v.place AS place,
// 					v.time_at AS time_at,
// 					CASE v.photo when null then "" else v.photo end AS photo,
// 					v.created_at AS created_at, v.updated_at AS updated_at,
// 					s{id:ID(s), .name, .code} AS owner,
// 					sv{id:ID(sv),.code,.name} AS catcher
// 				ORDER BY %s
// 				SKIP {skip}
// 				LIMIT {limit}
// 		  	`, params.Sort)
//
// 	paramsQuery := map[string]interface{}{
// 		"skip":        params.Skip,
// 		"limit":       params.Limit,
// 		"studentCode": studentCode,
// 	}
// 	res := []models.Violation{}
// 	cq := neoism.CypherQuery{
// 		Statement:  stmt,
// 		Parameters: paramsQuery,
// 		Result:     &res,
// 	}
// 	err := conn.Cypher(&cq)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(res) > 0 {
// 		return res, nil
// 	}
// 	return nil, nil
// }

// Get func to get a channelNotification
// int64
// models.Violation error
func (service violationService) Get(violationID int64) (models.Violation, error) {
	stmt := `
		MATCH (s:Student)<-(v:Violation)<-[:CATCH]-(sv:Supervisior)
		WHERE ID(v)= {violationID}
		RETURN
			ID(v) AS id,
			v.message AS message,
			v.place AS place,
			v.time_at AS time_at,
			CASE v.photo when null then "" else v.photo end AS photo,
			v.created_at AS created_at, v.updated_at AS updated_at,
			s{id:ID(s),name: s.last_name+" "+s.first_name, .code} AS owner,
			sv{id:ID(sv),.code,name: sv.last_name+" "+sv.first_name} AS catcher
			`
	params := map[string]interface{}{
		"violationID": violationID,
	}
	res := []models.Violation{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Violation{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Violation{}, nil
}

// Delete func
// int64
// bool error
func (service violationService) Delete(violationID int64) (bool, error) {
	stmt := `
	  	MATCH (v:Violation)
			WHERE ID(v) = {violationID}
			DETACH DELETE v
	  	`
	params := map[string]interface{}{
		"violationID": violationID,
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

// Create func to create a new ChannelNotification
// models.Violation
// int64 error
func (service violationService) Create(violation models.Violation, studentCode string, supervisiorID int64) (int64, error) {
	p := neoism.Props{
		"message": violation.Message,
		"photo":   violation.Photo,
		"time_at": violation.TimeAt,
		"place":   violation.Place,
		"status":  1,
	}
	stmt := `
		    MATCH(s:Student) WHERE toLower(s.code) = toLower({studentCode})
				MATCH(sv:Supervisior) WHERE ID(sv) = {supervisiorID}
		  	CREATE (s)<-[:OF]-(v:Violation{ props } )<-[r:CATCH]-(sv)
				SET v.created_at = TIMESTAMP()
				RETURN ID(v) as id
		  	`

	params := map[string]interface{}{
		"props":         p,
		"studentCode":   studentCode,
		"supervisiorID": supervisiorID,
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
// models.Violation
// models.Violation error
func (service violationService) Update(violation models.Violation) (models.Violation, error) {
	stmt := `
  	MATCH (sv:Supervisior)-[:CATCH]->(v:Violation)-[:OF]->(s:Student)
    WHERE ID(v) = {violationID}
		SET v.message = {message}, v.time_at = {time_at}, v.place = {place}, v.photo = {photo}, v.updated_at = TIMESTAMP()
    RETURN
		ID(s) AS id,
		v.message AS message,
		v.place AS place,
		v.time_at AS time_at,
		CASE v.photo when null then "" else v.photo end AS photo,
		v.created_at AS created_at, v.updated_at AS updated_at,
		s{id:ID(s), name: s.last_name+" "+s.first_name, .code} AS owner,
		v{id:ID(sv), name: sv.last_name+" "+sv.first_name, .code} AS catcher
  	`
	params := map[string]interface{}{
		"violationID": violation.ID,
		"message":     violation.Message,
		"time_at":     violation.TimeAt,
		"place":       violation.Place,
		"photo":       violation.Photo,
	}

	res := []models.Violation{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Violation{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.Violation{}, errors.New("Dont' update violation")
}

// CheckExistViolation func
// int64
// bool error
func (service violationService) CheckExistViolation(violationID int64) (bool, error) {
	stmt := `
		MATCH (v:Violation)
		WHERE ID(v)={violationID}
		RETURN ID(v) AS id
		`
	params := neoism.Props{
		"violationID": violationID,
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
		if res[0].ID == violationID {
			return true, nil
		}
	}
	return false, nil
}

// CheckPermission func
// int64 int64
// bool error
func (service violationService) CheckPermission(violationID int64, myUserID int64) (bool, error) {
	stmt := `
		MATCH (u:User)
		WHERE ID(u)={myUserID}
		MATCH(v:Violation)
		WHERE ID(v)={violationID}
		RETURN
			exists((u)-[:CATCHER]->(v)) OR exists((v)-[:OF]->(u)) AS allowed
		`
	params := neoism.Props{
		"myUserID":    myUserID,
		"violationID": violationID,
	}

	res := []struct {
		Allowed bool `json:"allowed"`
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
		return res[0].Allowed, nil
	}
	return false, nil
}

// CheckPermission func
// int64 int64
// bool error
func (service violationService) GetSupervisiorID(myUserID int64) (int64, error) {
	stmt := `
		MATCH (u:User)
		WHERE ID(u)={myUserID}
		MATCH(u)-[:IS_A{status:1}]->(sv:Supervisior)
		RETURN
			ID(sv) AS id
		`
	params := neoism.Props{
		"myUserID": myUserID,
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
