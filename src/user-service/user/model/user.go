package model

type User struct {
	Email    string
	Password []byte
	Name     string
	Role     Role
}
