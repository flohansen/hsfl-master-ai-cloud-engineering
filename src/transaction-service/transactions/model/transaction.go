package model

type Transaction struct {
	ID           uint64 `json:"id"`
	ChapterID    uint64 `json:"chapterID"`
	PayingUserID uint64 `json:"payingUserID"`
	Amount       uint64 `json:"amount"`
}
