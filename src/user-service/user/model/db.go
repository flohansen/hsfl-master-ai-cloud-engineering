package model

type DbUser struct {
	Email    string
	Password []byte
	Name     string
	Balance  int64
}
