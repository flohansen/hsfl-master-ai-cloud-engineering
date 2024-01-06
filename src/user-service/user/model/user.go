package model

type User struct {
	Id       uint64 `json:"id,omitempty" db:"id" fieldtag:"pk"`
	Email    string `json:"email,omitempty" db:"email"`
	Password []byte `json:"password,omitempty"  db:"password"`
	Name     string `json:"name,omitempty" db:"name"`
	Role     Role   `json:"role,omitempty" db:"role"`
}
