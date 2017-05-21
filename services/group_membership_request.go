package services

import (
	"github.com/kienbui1995/social-network-tlu-api/helpers"
	"github.com/kienbui1995/social-network-tlu-api/models"
)

// GroupMembershipRequestServiceInterface include method list
type GroupMembershipRequestServiceInterface interface {
	GetAll(params helpers.ParamsGetAll, groupID int64) ([]models.GroupMembershipRequest, error)
	Get(requestID int64) (models.GroupMembershipRequest, error)
	Delete(requestID int64) (bool, error)
	Create(groupID int64, userID int64) (int64, error)
	Update(requestID int64, status int) (models.GroupMembershipRequest, error)
	CheckExistRequest(requestID int64) (bool, error)
}

// groupService struct
type groupMembershipRequestService struct{}

// NewGroupMemberShipRequestService to constructor
func NewGroupMemberShipRequestService() groupMembershipRequestService {
	return groupMembershipRequestService{}
}
