package models

// GroupMembership struct
type GroupMembership struct {
	ID        int64       `json:"id"`
	CreatedAt int64       `json:"created_at,omitempty"`
	UpdatedAt int64       `json:"updated_at,omitempty"`
	User      UserObject  `json:"user"`
	Group     GroupObject `json:"group"`
	Status    int         `json:"status"` // 1: Joined; 2: Blocked;
}

// IsEmpty func to check membership is null
func (membership GroupMembership) IsEmpty() bool {
	return membership == GroupMembership{}
}
