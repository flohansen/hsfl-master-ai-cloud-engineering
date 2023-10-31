package model

type Product struct {
	Id          uint64 `json:"id,omitempty"`
	Description string `json:"description,omitempty"`
	Ean         uint64 `json:"ean,omitempty"`
}
