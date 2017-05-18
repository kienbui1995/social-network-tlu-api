package models

//Subscription struct
type Subscription struct {
	SubscriberID int64 `json:"id"`
	UserID       int64 `json:"user_id"`
	ObjectID     int64 `json:"object_id"`
	CreatedAt    int64 `json:"created_at"`
}
