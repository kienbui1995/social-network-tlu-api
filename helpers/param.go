package helpers

// ParamsGetAll struct
type ParamsGetAll struct {
	Filter string `json:"filter,omitempty"`
	Skip   int    `json:"skip,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Type   string `json:"type,omitempty"`
}
