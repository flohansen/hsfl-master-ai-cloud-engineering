package model

type Product struct {
	Id          uint64 `json:"id,omitempty" db:"id" fieldtag:"pk"`
	Description string `json:"description,omitempty" db:"description"`
	Ean         string `json:"ean,omitempty" db:"ean"`
}
