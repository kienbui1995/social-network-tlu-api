package models

// GroupMembership struct
type GroupMembership struct {
	ID        int64        `json:"id"`
	CreatedAt int64        `json:"created_at,omitempty"`
	UpdatedAt int64        `json:"updated_at,omitempty"`
	User      *UserObject  `json:"user,omitempty"`
	Group     *GroupObject `json:"group,omitempty"`
	Role      int          `json:"role,omitempty"` // 1:member; 2: admin; 3: creator 4: block
	Status    int          `json:"status,omitempty"`
	CanEdit   bool         `json:"can_edit,omitempty"`
	CanDelete bool         `json:"can_delete,omitempty"`
}

// IsEmpty func to check membership is null
func (membership GroupMembership) IsEmpty() bool {
	return membership == GroupMembership{}
}
