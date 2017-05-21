package models

// GroupMembershipRequest struct
type GroupMembershipRequest struct {
	ID        int64       `json:"id"`
	CreatedAt int64       `json:"created_at"`
	UpdatedAt int64       `json:"updated_at"`
	User      UserObject  `json:"user"`
	Group     GroupObject `json:"group"`
	Status    int         `json:"status"` // 1: Pending; 2: Accepted; 3: Declined
}
