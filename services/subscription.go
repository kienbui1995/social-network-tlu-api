package services

import (
	"errors"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// SubscriptionServiceInterface include method list
type SubscriptionServiceInterface interface {
	CreateSubscription(fromID int64, toID int64) (int64, error)
	DeleteSubcription(fromID int64, toID int64) (bool, error)
	CheckExistSubscription(fromID int64, toID int64) (bool, error)
	GetSubscriptions(userID int64) ([]models.UserFollowObject, error)
	GetFollowers(userID int64) ([]models.UserFollowObject, error)
	GetFollowerIDs(userID int64) ([]int64, error)
	CheckExistObject(objectID int64, objectType string) (bool, error)
}

// subscriptionService struct
type subscriptionService struct{}

// NewSubscriberService to constructor
func NewSubscriberService() SubscriptionServiceInterface {
	return subscriptionService{}
}

// CreateUserSubscriber func
// int64 int64
// int64 error
func (service subscriptionService) CreateSubscription(fromID int64, toID int64) (int64, error) {
	stmt := `
	MATCH(from:User) WHERE ID(from) = {fromid}
  MATCH (to) WHERE ID(to) = {toid}
  CREATE UNIQUE (from)-[f:FOLLOW{created_at:TIMESTAMP()}]->(to)
	SET from.followings = from.followings+1, to.followers = to.followers+1
 	RETURN ID(f) as id
	`
	res := []struct {
		ID int64 `json:"id"`
	}{}
	params := neoism.Props{
		"fromid": fromID,
		"toid":   toID,
	}
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
	return -1, errors.New("Don't create follow relationship")
}

//DeleteUserSubscriber fun
func (service subscriptionService) DeleteSubcription(fromID int64, toID int64) (bool, error) {
	stmt := `
  	MATCH (from:User)-[f:FOLLOW]->(to)
    WHERE ID(from) = {fromID} AND ID(to) = {toID}
		SET from.followings = from.followings-1, to.followers = to.followers -1
    DELETE f
  	`
	params := neoism.Props{
		"fromID": fromID,
		"toID":   toID,
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

// GetSubscribers func
func (service subscriptionService) GetSubscriptions(userID int64) ([]models.UserFollowObject, error) {
	stmt := `
	MATCH (userid:User)-[f:FOLLOW]->(u:User)
	WHERE ID(userid)= {userid}
	RETURN ID(u) as id, u.username as username, u.avatar as avatar, u.full_name as full_name,
	TRUE	AS is_followed
  	`
	params := neoism.Props{
		"userid": userID,
	}
	res := []models.UserFollowObject{}

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

// GetFollowers func
func (service subscriptionService) GetFollowers(userID int64) ([]models.UserFollowObject, error) {
	stmt := `
	MATCH (u:User)-[f:FOLLOW]->(user:User)
	WHERE ID(user)= {userid}
	RETURN ID(u) as id, u.username as username, u.avatar as avatar, u.full_name as full_name,
	exists((user)-[:FOLLOW]->(u)) as is_followed
  	`
	params := neoism.Props{
		"userid": userID,
	}
	res := []models.UserFollowObject{}

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

// GetFollowerIDs func
func (service subscriptionService) GetFollowerIDs(userID int64) ([]int64, error) {
	stmt := `
	MATCH (u:User)-[f:FOLLOW]->(user:User)
	WHERE ID(user)= {userid}
	RETURN ID(u) as id
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
		return nil, err
	}
	if len(res) > 0 {
		var ids []int64
		for index := 0; index < len(res); index++ {
			ids = append(ids, res[index].ID)
		}
		return ids, nil

	}
	return nil, nil
}

// CheckExistObject
// int64 string
// bool error
func (service subscriptionService) CheckExistObject(objectID int64, objectType string) (bool, error) {

	stmt := `
		MATCH (s)
		WHERE ID(s) = {objectID} AND {objectType} IN LABELS(s)
		RETURN ID(s) as id
		`
	params := neoism.Props{
		"objectID":   objectID,
		"objectType": objectType,
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
		if res[0].ID == objectID {
			return true, nil
		}
	}
	return false, nil
}

// CheckExistSubscription func
// int64 int64
// bool error
func (service subscriptionService) CheckExistSubscription(fromID int64, toID int64) (bool, error) {
	stmt := `
  	MATCH (from:User)-[f:FOLLOW]->(to)
		WHERE ID(from) = {fromID} AND ID(to) = {toID}
		RETURN ID(f) AS id
  	`
	params := neoism.Props{
		"fromID": fromID,
		"toID":   toID,
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
