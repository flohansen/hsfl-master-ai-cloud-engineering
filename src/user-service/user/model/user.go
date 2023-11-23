package model

type User struct {
	Id       uint64 `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password []byte `json:"password,omitempty"`
	Name     string `json:"name,omitempty"`
	Role     Role   `json:"role,omitempty"`
}
