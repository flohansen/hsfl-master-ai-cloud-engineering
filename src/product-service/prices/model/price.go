package model

type Price struct {
	UserId    uint64  `json:"userId,omitempty" db:"userId"`
	ProductId uint64  `json:"productId,omitempty" db:"productId"`
	Price     float32 `json:"price,omitempty" db:"price"`
}
