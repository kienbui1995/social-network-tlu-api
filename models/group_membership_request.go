package models

// GroupMembershipRequest struct
type GroupMembershipRequest struct {
	ID              int64        `json:"id"`
	CreatedAt       int64        `json:"created_at,omitempty"`
	UpdatedAt       int64        `json:"updated_at,omitempty"`
	User            *UserObject  `json:"user"`
	Group           *GroupObject `json:"group"`
	RequestMessage  string       `json:"request_message,omitempty"`
	ResponseMessage string       `json:"response_message,omitempty"`
	Status          int          `json:"status"` // 1: Pending; 2: Accepted; 3: Declined
}

// IsEmpty func to check empty entity
func (model GroupMembershipRequest) IsEmpty() bool {
	return model == GroupMembershipRequest{}
}

// InfoGroupMembershipRequest struct to update
type InfoGroupMembershipRequest struct {
	RequestMessage  string `json:"request_message,omitempty"`
	ResponseMessage string `json:"response_message,omitempty"`
	Status          int    `json:"status"` // 1: Pending; 2: Accepted; 3: Declined
}
