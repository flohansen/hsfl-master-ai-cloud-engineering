package model

type UserShoppingList struct {
	Id          uint64 `json:"id,omitempty"`
	UserId      uint64 `json:"userId,omitempty"`
	Description string `json:"description,omitempty"`
	Completed   bool   `json:"completed,omitempty"`
}
