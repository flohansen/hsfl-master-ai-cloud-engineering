package model

type User struct {
	Id       uint64
	Email    string
	Password []byte
	Name     string
	Role     Role
}
