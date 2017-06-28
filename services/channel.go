package services

import (
	"fmt"

	"github.com/jmcvetta/neoism"
	"github.com/kienbui1995/social-network-tlu-api/configs"
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// ChannelServiceInterface include method list
type ChannelServiceInterface interface {
	//  Normal User
	GetAll(params helpers.ParamsGetAll, myUserID int64) ([]models.Channel, error)
	Get(channelID int64, myUserID int64) (models.Channel, error)

	CheckExistChannel(channelID int64) (bool, error)
	CheckUserRole(channelID int64, userID int64) (int, error)
	GetFollowers(params helpers.ParamsGetAll, channelID int64, myUserID int64) ([]models.UserFollowObject, error)

	// Only Admin
	Create(channel models.Channel, adminID int64) (int64, error)
	Update(channelID int64, newChannel models.InfoChannel) (models.Channel, error)
	Delete(channelID int64) (bool, error)
}

// channelService struct
type channelService struct{}

// NewChannelService to constructor
func NewChannelService() ChannelServiceInterface {
	return channelService{}
}

// GetAll func
// helpers.ParamsGetAll int64
// []models.Channel error
func (service channelService) GetAll(params helpers.ParamsGetAll, myUserID int64) ([]models.Channel, error) {
	var stmt string
	stmt = fmt.Sprintf(`
				MATCH(me:User) WHERE ID(me) = {myUserID}
				MATCH (c:Channel)
				WITH
					c{
						id: ID(c),
						.*,
						is_followed: exists((me)-[:FOLLOW]->(c)),
						is_admin: exists((me)-[:MANAGE]->(c))
						} AS channel
				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
				RETURN collect(channel) AS channels
				`, "channel."+params.Sort)

	paramsQuery := map[string]interface{}{
		"myUserID": myUserID,
		"skip":     params.Skip,
		"limit":    params.Limit,
	}
	res := []struct {
		Channels []models.Channel `json:"channels"`
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
		return res[0].Channels, nil
	}
	return nil, nil
}

// Get func
// int64 int64
// models.Channel error
func (service channelService) Get(channelID int64, myUserID int64) (models.Channel, error) {
	stmt := `
				MATCH(me:User) WHERE ID(me) = {myUserID}
				MATCH (c:Channel) WHERE ID(c) = {channelID}
				RETURN
					c{
						id: ID(c),
						.*,
						is_followed: exists((me)-[:FOLLOW]->(c)),
						is_admin: exists((me)-[:MANAGE]->(c))
						} AS channel
				`

	paramsQuery := map[string]interface{}{
		"channelID": channelID,
		"myUserID":  myUserID,
	}
	res := []struct {
		Channel models.Channel `json:"channel"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Channel{}, err
	}
	if len(res) > 0 {
		return res[0].Channel, nil
	}
	return models.Channel{}, nil
}

// Delete func
// int64
// bool error
func (service channelService) Delete(channelID int64) (bool, error) {
	return true, nil
}

// func Create
// models.Channel int64
// int64 error
func (service channelService) Create(channel models.Channel, myUserID int64) (int64, error) {
	p := neoism.Props{
		"name":        channel.Name,
		"description": channel.Description,
		"avatar":      channel.Avatar,
		"cover":       channel.Cover,
		"followers":   0,
		"posts":       0,
		"status":      1,
	}
	stmt := `
			MATCH(u:User) WHERE ID(u) = {myUserID}
			CREATE (c:Channel{props})<-[r:MANAGE]-(u)
			SET c.created_at = TIMESTAMP()
			RETURN ID(c) as id
			`
	params := map[string]interface{}{
		"props":    p,
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

// Update func
// int64 models.Channel
// models.Channel error
func (service channelService) Update(channelID int64, newChannel models.InfoChannel) (models.Channel, error) {
	stmt := `
		MATCH (c:Channel)
		WHERE ID(c) = {channelID}
		SET c+= {p}, c.updated_at = TIMESTAMP()
		RETURN
			c{id:ID(c), .*} AS channel
		`
	params := neoism.Props{
		"channelID": channelID,
		"p":         newChannel,
	}
	res := []struct {
		Channel models.Channel `json:"channel"`
	}{}
	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: params,
		Result:     &res,
	}
	err := conn.Cypher(&cq)
	if err != nil {
		return models.Channel{}, err
	}
	if len(res) > 0 {
		if res[0].Channel.ID >= 0 {
			return res[0].Channel, nil
		}
	}
	return models.Channel{}, nil
}

// CheckExistChannel func
// int64
// bool error
func (service channelService) CheckExistChannel(channelID int64) (bool, error) {
	stmt := `
		MATCH (c:Channel)
		WHERE ID(c)={channelID}
		RETURN ID(c) AS id
		`
	params := neoism.Props{
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
		return false, err
	}

	if len(res) > 0 {
		if res[0].ID == channelID {
			return true, nil
		}
	}
	return false, nil
}

// CheckUserRole func
// int64 int64
// int error
func (service channelService) CheckUserRole(channelID int64, userID int64) (int, error) {
	stmt := `
	MATCH (c:Channel)	WHERE ID(c)= {channelID}
	MATCH (u:User) WHERE ID(u) = {userID}
	RETURN
		exists((u)-[:MANAGE]->(c)) AS is_admin,
		exists((u)-[:FOLLOW]->(c)) AS is_followed
	`
	paramsQuery := neoism.Props{
		"channelID": channelID,
		"userID":    userID,
	}
	res := []struct {
		IsAdmin    bool `json:"is_admin"`
		IsFollowed bool `json:"is_followed"`
	}{}

	cq := neoism.CypherQuery{
		Statement:  stmt,
		Parameters: paramsQuery,
		Result:     &res,
	}

	err := conn.Cypher(&cq)
	if err != nil {
		return -1, err
	}
	if len(res) > 0 {

		if res[0].IsAdmin {
			return configs.IIsAdminChannel, nil
		}
		if res[0].IsFollowed {
			return configs.IFollowedChannel, nil
		}

	}
	return -1, nil
}

// GetFollowers func
// helpers.ParamsGetAll int64 int64
// []models.UserFollowObject error
func (service channelService) GetFollowers(params helpers.ParamsGetAll, channelID int64, myUserID int64) ([]models.UserFollowObject, error) {
	var stmt string
	stmt = fmt.Sprintf(`
				MATCH(me:User) WHERE ID(me) = {myUserID}
				MATCH (c:Channel)<-[f:FOLLOW]-(u:User)
				WHERE ID(c) = {channelID}
				WITH
					u{
						id: ID(u),
						username: .username,
						full_name: .full_name,
						avatar: .avatar,
						is_followed: exists((me)-[:FOLLOW]->(u))
					} AS user

				ORDER BY %s
				SKIP {skip}
				LIMIT {limit}
				RETURN collect(user) AS users
				`, "user."+params.Sort)

	paramsQuery := map[string]interface{}{
		"myuserid":  myUserID,
		"channelID": channelID,
		"skip":      params.Skip,
		"limit":     params.Limit,
	}
	res := []struct {
		Users []models.UserFollowObject `json:"users"`
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
		return res[0].Users, nil
	}
	return nil, nil
}
