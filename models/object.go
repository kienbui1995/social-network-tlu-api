package models

// UserObject struct
type UserObject struct {
	ID       int64  `json:"id"`
	Username string `json:"username,omitempty"`
	FullName string `json:"full_name,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
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

// PostObject struct
type PostObject struct {
	ID      int64        `json:"id"`
	Message string       `json:"message"`
	Photo   string       `json:"photo,omitempty"`
	Owner   *UserObject  `json:"owner,omitempty"`
	Place   *GroupObject `json:"place,omitempty"`
}

// CommentObject struct
type CommentObject struct {
	ID       int64            `json:"id"`
	Message  string           `json:"message"`
	Mentions []*MentionObject `json:"mentions,omitempty"`
	Owner    *UserObject      `json:"owner,omitempty"`
	Post     *PostObject      `json:"post,omitempty"`
}

// ChannelObject struct
type ChannelObject struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// ActorObject struct
type ActorObject struct {
	*UserObject `json:",omitempty"`
	// *ChannelObject `json:",omitempty"`
}

// IsEmpty func
func (object ActorObject) IsEmpty() bool {
	return object == ActorObject{}
}
