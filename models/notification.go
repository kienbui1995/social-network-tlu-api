package models

// Notification struct
type Notification struct {
	UserID     int64  `json:"user_id"`
	ObjectID   int64  `json:"object_id"`
	ObjectType string `json:"object_type"`
	Title      string `json:"title"`
	Message    string `json:"message"`
}
