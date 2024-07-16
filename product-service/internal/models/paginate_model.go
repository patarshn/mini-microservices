package models

type Paginate struct {
	From int   `json:"from"`
	To   int   `json:"to"`
	Next bool  `json:"next"`
	Prev bool  `json:"prev"`
	Page int64 `json:"page"`
	Data []any `json:"data,omitempy"`
}
