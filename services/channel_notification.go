package services

import (
	"errors"
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// ChannelNotificationServiceInterface include method list
type ChannelNotificationServiceInterface interface {
	GetAll(params helpers.ParamsGetAll, myUserID int64) ([]models.ChannelNotification, error)
	Get(channelNotificationID int64) (models.ChannelNotification, error)
	Delete(channelNotificationID int64) (bool, error)
	Create(channelNotification models.ChannelNotification, channelID int64) (int64, error)
	Update(channelNotification models.ChannelNotification) (models.ChannelNotification, error)
	CheckExistChannelNotification(channelNotificationID int64) (bool, error)

	GetAllOfChannel(params helpers.ParamsGetAll, channelID int64) ([]models.ChannelNotification, error)
	CheckPermission(channelID int64, myUserID int64) (bool, error)
	CheckPermissionByNotificationID(channelNotificationID int64, myUserID int64) (bool, error)
}

// channelNotificationService struct
type channelNotificationService struct{}

// NewChannelNotificationService to constructor
func NewChannelNotificationService() ChannelNotificationServiceInterface {
	return channelNotificationService{}
}

// GetAll func
// helpers.ParamsGetAll int64
// []models.ChannelNotification error
func (service channelNotificationService) GetAll(params helpers.ParamsGetAll, myUserID int64) ([]models.ChannelNotification, error) {
	var stmt string
	stmt = fmt.Sprintf(`
				MATCH(me:User) WHERE ID(me) = {myUserID}
		  	MATCH (s:ChannelNotification)<-[:CREATE]-(c:Channel)<-[f:FOLLOW]-(me)
				RETURN
					ID(s) AS id,
					s.title AS title,
					substring(s.message,0,250) AS message, length(s.message)>250 AS summary,
					s.place AS place,
					s.time AS time,
					CASE s.photo when null then "" else s.photo end AS photo,
					s.created_at AS created_at, s.updated_at AS updated_at,
					c{id:ID(c), .name, .short_name, .avatar} AS owner
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
		  	`, params.Sort)

	paramsQuery := map[string]interface{}{
		"myUserID": myUserID,
		"skip":     params.Skip,
		"limit":    params.Limit,
	}
	res := []models.ChannelNotification{}
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

// Get func to get a channelNotification
// int64 int64
// models.ChannelNotification error
func (service channelNotificationService) Get(channelNotificationID int64) (models.ChannelNotification, error) {
	stmt := `
		MATCH (s:ChannelNotification)<-[:CREATE]-(c:Channel)
		WHERE ID(s) = {channelNotificationID}
		RETURN
			ID(s) AS id,
			s.title AS title,
			substring(s.message,0,250) AS message, length(s.message)>250 AS summary,
			s.place AS place,
			s.time AS time,
			CASE s.photo when null then "" else s.photo end AS photo,
			s.created_at AS created_at, s.updated_at AS updated_at,
			c{id:ID(c), .name, .short_name, .avatar} AS owner
			`
	params := map[string]interface{}{
		"channelNotificationID": channelNotificationID,
	}
	res := []models.ChannelNotification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.ChannelNotification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.ChannelNotification{}, nil
}

// Delete func
// int64
// bool error
func (service channelNotificationService) Delete(channelNotificationID int64) (bool, error) {
	stmt := `
	  	MATCH (s:ChannelNotification)
			WHERE ID(s) = {channelNotificationID}
			MATCH (c:Channel)-[:CREATE]->(s)
			SET c.notifications = c.notifications - 1
			DETACH DELETE s
	  	`
	params := map[string]interface{}{
		"channelNotificationID": channelNotificationID,
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
// models.ChannelNotification int64
// int64 error
func (service channelNotificationService) Create(channelNotification models.ChannelNotification, channelID int64) (int64, error) {
	p := neoism.Props{
		"message": channelNotification.Message,
		"photo":   channelNotification.Photo,
		"title":   channelNotification.Title,
		"time":    channelNotification.Time,
		"place":   channelNotification.Place,
		"status":  1,
	}
	stmt := `
		    MATCH(c:Channel) WHERE ID(c) = {channelID}
		  	CREATE (s:ChannelNotification{ props } )<-[r:CREATE]-(c)
				SET c.notifications = c.notifications + 1, s.created_at = TIMESTAMP()
				RETURN ID(s) as id
		  	`

	params := map[string]interface{}{
		"props":     p,
		"channelID": channelID,
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
// models.ChannelNotification
// models.ChannelNotification error
func (service channelNotificationService) Update(channelNotification models.ChannelNotification) (models.ChannelNotification, error) {
	stmt := `
  	MATCH (s:ChannelNotification)<-[:CREATE]-(c:Channel)
    WHERE ID(s) = {channelNotificationID}
		SET s.message = {message}, s.title = {title}, s.time = {time}, s.place = {place}, s.updated_at = TIMESTAMP()
    RETURN
		ID(s) AS id,
		s.title AS title,
		substring(s.message,0,250) AS message, length(s.message)>250 AS summary,
		s.place AS place,
		s.time AS time,
		CASE s.photo when null then "" else s.photo end AS photo,
		s.created_at AS created_at, s.updated_at AS updated_at,
		c{id:ID(c), .name, .short_name, .avatar} AS owner
  	`
	params := map[string]interface{}{
		"channelNotificationID": channelNotification.ID,
		"message":               channelNotification.Message,
		"title":                 channelNotification.Title,
		"time":                  channelNotification.Time,
		"place":                 channelNotification.Place,
	}

	res := []models.ChannelNotification{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.ChannelNotification{}, err
	}
	if len(res) > 0 {
		return res[0], nil
	}
	return models.ChannelNotification{}, errors.New("Dont' update channel notification")
}

// CheckExistChannelNotification func
// int64
// bool error
func (service channelNotificationService) CheckExistChannelNotification(channelNotificationID int64) (bool, error) {
	stmt := `
		MATCH (u:ChannelNotification)
		WHERE ID(u)={channelNotificationID}
		RETURN ID(u) AS id
		`
	params := neoism.Props{
		"channelNotificationID": channelNotificationID,
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
		if res[0].ID == channelNotificationID {
			return true, nil
		}
	}
	return false, nil
}

// GetAllOfChannel func
// helpers.ParamsGetAll
// []models.ChannelNotification error
func (service channelNotificationService) GetAllOfChannel(params helpers.ParamsGetAll, channelID int64) ([]models.ChannelNotification, error) {
	stmt := fmt.Sprintf(`
		MATCH(c:Channel) WHERE ID(c) = {channelID}
		MATCH (s:ChannelNotification)<-[:CREATE]-(c)
		RETURN
			ID(s) AS id,
			s.title AS title,
			substring(s.message,0,250) AS message, length(s.message)>250 AS summary,
			s.place AS place,
			s.time AS time,
			CASE s.photo when null then "" else s.photo end AS photo,
			s.created_at AS created_at, s.updated_at AS updated_at,
			c{id:ID(c), .name, .short_name, .avatar} AS owner
		ORDER BY %s
		SKIP {skip}
		LIMIT {limit}
	  	`, params.Sort)

	paramsQuery := map[string]interface{}{
		"channelID": channelID,
		"skip":      params.Skip,
		"limit":     params.Limit,
	}
	res := []models.ChannelNotification{}
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

// CheckPermission func
// int64 int64
// bool error
func (service channelNotificationService) CheckPermission(channelID int64, myUserID int64) (bool, error) {
	stmt := `
		MATCH (u:User)
		WHERE ID(u)={myUserID}
		MATCH(c:Channel)
		WHERE ID(c)={channelID}
		RETURN
			exists((u)-[:MANAGE]->(c)) AS is_admin
		`
	params := neoism.Props{
		"channelID": channelID,
		"myUserID":  myUserID,
	}

	res := []struct {
		IsAdmin bool `json:"is_admin"`
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
		return res[0].IsAdmin, nil
	}
	return false, nil
}

// CheckPermissionByNotificationID func
// int64 int64
// bool error
func (service channelNotificationService) CheckPermissionByNotificationID(channelNotificationID int64, myUserID int64) (bool, error) {
	stmt := `
		MATCH (u:User)
		WHERE ID(u)={myUserID}
		MATCH(c:Channel)-[:CREATE]->(n:ChannelNotification)
		WHERE ID(n)={channelNotificationID}
		RETURN
			exists((u)-[:MANAGE]->(c)) AS is_admin
		`
	params := neoism.Props{
		"channelNotificationID": channelNotificationID,
		"myUserID":              myUserID,
	}

	res := []struct {
		IsAdmin bool `json:"is_admin"`
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
		return res[0].IsAdmin, nil
	}
	return false, nil
}
