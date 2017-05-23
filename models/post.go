package models

// Post struct to include info of core Post
type Post struct {
	Owner     UserObject `json:"owner"`
	ID        int64      `json:"id"` //Id
	Message   string     `json:"message"`
	Summary   bool       `json:"summary,omitempty"` // true, if message length more than 250
	Photo     string     `json:"photo,omitempty"`
	UpdatedAt int64      `json:"updated_at,omitempty"`
	CreatedAt int64      `json:"created_at,omitempty"`
	Status    int        `json:"status,omitempty"`
	Privacy   int        `json:"privacy,omitempty"` // 1: public; 2: followers; 3: private
	Likes     int        `json:"likes,omitempty"`
	Comments  int        `json:"comments,omitempty"`
	Shares    int        `json:"shares,omitempty"`
	IsLiked   bool       `json:"is_liked"`
	CanEdit   bool       `json:"can_edit"`
	CanDelete bool       `json:"can_delete"`
}

// Posts list
type Posts []Post

// IsEmpty func to check entity empty
func (p Post) IsEmpty() bool {
	return p == Post{}
}

// InfoPost struct to include info of core Post
type InfoPost struct {
	Message string `json:"message"`
	Photo   string `json:"photo,omitempty"`
	Status  int    `json:"status,omitempty"`
	Privacy int    `json:"privacy,omitempty"` // 1: public; 2: followers; 3: private
}
