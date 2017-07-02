package models

// Channel struct
type Channel struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	ShortName     string `json:"short_name"`
	Description   string `json:"description,omitempty"`
	Followers     int    `json:"followers,omitempty"`
	Notifications int    `json:"notifications,omitempty"`
	Avatar        string `json:"avatar,omitempty"`
	Cover         string `json:"cover,omitempty"`
	CreatedAt     int64  `json:"created_at,omitempty"`
	UpdatedAt     int64  `json:"updated_at,omitempty"`
	Status        int    `json:"status,omitempty"`
	IsAdmin       bool   `json:"is_admin,omitempty"`
	IsFollowed    bool   `json:"is_followed,omitempty"`
}

// IsEmpty func
func (channel Channel) IsEmpty() bool {
	return channel == Channel{}
}

// InfoChannel struct
type InfoChannel struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ShortName   string `json:"short_name"`
	Description string `json:"description,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Cover       string `json:"cover,omitempty"`
	Status      int    `json:"status,omitempty"`
}
