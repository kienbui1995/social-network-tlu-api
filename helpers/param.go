package helpers

// ParamsGetAll struct
type ParamsGetAll struct {
	Filter     string                 `json:"filter,omitempty"`
	Sort       string                 `json:"sort,omitempty"`
	Skip       int                    `json:"skip,omitempty"`
	Limit      int                    `json:"limit,omitempty"`
	Type       string                 `json:"type,omitempty"`
	Properties map[string]interface{} `json:"properties,omitempty"`
}
