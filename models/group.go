package models

// Group struct
type Group struct {
	ID              int64        `json:"id"`
	Name            string       `json:"name"`
	Description     string       `json:"description,omitempty"`
	Class           *ClassObject `json:"class,omitempty"`
	Type            int          `json:"type,omitempty"`
	PendingRequests int          `json:"pending_requests,omitempty"`
	Members         int          `json:"members,omitempty"`
	Posts           int          `json:"posts,omitempty"`
	Avatar          string       `json:"avatar,omitempty"`
	Cover           string       `json:"cover,omitempty"`
	Privacy         int          `json:"privacy,omitempty"`
	CreatedAt       int64        `json:"created_at,omitempty"`
	UpdatedAt       int64        `json:"updated_at,omitempty"`
	Status          int          `json:"status,omitempty"`
	CanRequest      bool         `json:"can_request,omitempty"`
	CanJoin         bool         `json:"can_join,omitempty"`
	IsPending       bool         `json:"is_pending,omitempty"`
	IsAdmin         bool         `json:"is_admin,omitempty"`
	IsMember        bool         `json:"is_member,omitempty"`
}

// IsEmpty func
func (group Group) IsEmpty() bool {
	return group == Group{}
}

// GroupJoin struct
type GroupJoin struct {
	Group
}

// InfoGroup struct to update method
type InfoGroup struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
	Privacy     int    `json:"privacy,omitempty"`
	Status      int    `json:"status,omitempty"`
}
