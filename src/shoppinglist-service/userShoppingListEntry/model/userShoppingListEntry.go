package model

type UserShoppingListEntry struct {
	ShoppingListId uint64 `json:"shoppingListId,omitempty"`
	ProductId      uint64 `json:"productId,omitempty"`
	Count          uint16 `json:"count,omitempty"`
	Note           string `json:"note,omitempty"`
	Checked        bool   `json:"checked,omitempty"`
}
