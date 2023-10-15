package model

type Price struct {
	UserId    uint64  `json:"userId,omitempty"`
	ProductId uint64  `json:"productId,omitempty"`
	Price     float32 `json:"price,omitempty"`
}
