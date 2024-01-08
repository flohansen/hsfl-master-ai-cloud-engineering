package model

type UserShoppingList struct {
	Id          uint64 `json:"id,omitempty" db:"id" fieldtag:"pk"`
	UserId      uint64 `json:"userId,omitempty" db:"userId"`
	Description string `json:"description,omitempty" db:"description"`
	Completed   bool   `json:"completed,omitempty" db:"completed"`
}
