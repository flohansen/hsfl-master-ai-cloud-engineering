package model

type Role int64

const (
	Customer Role = iota
	Merchant
	Administrator
)
