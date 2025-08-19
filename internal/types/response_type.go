package types

type Metadata struct {
	Total  int64  `json:"total"`
	Count  int    `json:"count"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Order  string `json:"order"`
	Search string `json:"search"`
}
