package models

// Group struct
type Group struct {
	ID              int64  `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description,omitempty"`
	PendingRequests int    `json:"pending_requests"`
	Members         int    `json:"members,omitempty"`
	Posts           int    `json:"posts,omitempty"`
	Avatar          string `json:"avatar,omitempty"`
	Cover           string `json:"cover,omitempty"`
	Privacy         int    `json:"privacy,omitempty"`
	CreatedAt       int64  `json:"created_at,omitempty"`
	UpdatedAt       int64  `json:"updated_at,omitempty"`
	Status          int    `json:"status,omitempty"`
}

// GroupObject struct
type GroupObject struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// GroupJoin struct
type GroupJoin struct {
	Group
	IsJoined bool `json:"is_joined"`
}

// GroupAdmin struct
type GroupAdmin struct {
	Group
	IsAdmin bool `json:"is_admin"`
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
