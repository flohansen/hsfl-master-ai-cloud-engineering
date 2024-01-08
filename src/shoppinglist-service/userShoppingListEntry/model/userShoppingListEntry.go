package model

type UserShoppingListEntry struct {
	ShoppingListId uint64 `json:"shoppingListId,omitempty" db:"shoppingListId"`
	ProductId      uint64 `json:"productId,omitempty" db:"productId"`
	Count          uint16 `json:"count,omitempty" db:"count"`
	Note           string `json:"note,omitempty" db:"note"`
	Checked        bool   `json:"checked,omitempty" db:"checked"`
}
