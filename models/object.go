package models

// UserObject struct
type UserObject struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Avatar   string `json:"avatar"`
}

//UserFollowObject struct for a sub user for get all user; search user, list user any where
type UserFollowObject struct {
	UserObject
	IsFollowed bool `json:"is_followed"`
}

//UserLikedObject struct for a user liked things
type UserLikedObject struct {
	UserObject
	LikedAt int64 `json:"liked_at"`
}

// MentionObject struct
type MentionObject struct {
	ID     int64 `json:"id"`
	Length int64 `json:"length"`
	Offset int64 `json:"offset"`
}

// GroupObject struct
type GroupObject struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// PlaceObject struct
type PlaceObject struct {
	GroupObject
}

// IsEmpty func
func (object PlaceObject) IsEmpty() bool {
	return object == PlaceObject{}
}
