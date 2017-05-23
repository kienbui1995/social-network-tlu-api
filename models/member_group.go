package models

// GroupMember struc
type GroupMember struct {
	ID         int64            `json:"id"`
	User       UserFollowObject `json:"user"`
	JoinedAt   int64            `json:"joined_at,omitempty"`
	AcceptedBy UserFollowObject `json:"accepted_by,omitempty"`
}
