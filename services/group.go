package services

// import (
// 	"github.com/kienbui1995/social-network-tlu-api/helpers"
// 	"github.com/kienbui1995/social-network-tlu-api/models"
// )
//
// // GroupServiceInterface include method list
// type GroupServiceInterface interface {
// 	GetAll(params helpers.ParamsGetAll, myUserID int64) ([]models.Group, error)
// 	Get(groupID int64, myUserID int64) (models.Group, error)
// 	Delete(groupID int64) (bool, error)
// 	Create(group models.Group, myUserID int64) (int64, error)
// 	Update(group models.Group) (models.Group, error)
// 	CheckExistGroup(groupID int64) (bool, error)
// 	GetMembers(params helpers.ParamsGetAll, groupID int64) ([]models.UserFollowObject, error)
// 	IncreasePosts(userID int64) (bool, error)
// 	DecreasePosts(userID int64) (bool, error)
// }
//
// // groupService struct
// type groupService struct{}
//
// // NewGroupService to constructor
// func NewGroupService() groupService {
// 	return groupService{}
// }
