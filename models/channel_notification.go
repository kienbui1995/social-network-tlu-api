package models

// ChannelNotification struct
type ChannelNotification struct {
	ID        int64          `json:"id"`
	Owner     *ChannelObject `json:"owner,omitempty"`
	Title     string         `json:"title"`
	Message   string         `json:"message"`
	Photo     string         `json:"photo,omitempty"`
	Time      string         `json:"time,omitempty"`
	Place     string         `json:"place,omitempty"`
	UpdatedAt int64          `json:"updated_at,omitempty"`
	CreatedAt int64          `json:"created_at,omitempty"`
	Status    int            `json:"status,omitemtpy"`
	SeenAt    int64          `json:"seen_at,omitempty"`
}

// IsEmpty func
func (channel ChannelNotification) IsEmpty() bool {
	return channel == ChannelNotification{}
}
