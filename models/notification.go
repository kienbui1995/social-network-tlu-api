package models

// Notification struct
type Notification struct {
	ID          int64          `json:"id"`
	Actor       *ActorObject   `json:"actor,omitempty"`
	Action      int            `json:"action,omitempty"` //liked/commented/posted/mentioned/followed
	TotalAction int            `json:"total_action,omitempty"`
	LastPost    *PostObject    `json:"last_post,omitempty"`
	LastComment *CommentObject `json:"last_comment,omitempty"`
	LastMention *MentionObject `json:"last_mention,omitempty"`
	LastUser    *UserObject    `json:"last_user,omitempty"`
	Group       *GroupObject   `json:"group,omitempty"`
	Title       string         `json:"title,omitempty"`
	Message     string         `json:"message,omitempty"`
	UpdatedAt   int64          `json:"updated_at,omitempty"`
	CreatedAt   int64          `json:"created_at,omitempty"`
	SeenAt      int64          `json:"seen_at,omitempty"`
}
