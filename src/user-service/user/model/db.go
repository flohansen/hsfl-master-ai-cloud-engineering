package model

type UpdateUser struct {
	ProfileName string `json:"profile_name"`
	Balance     int64  `json:"balance"`
}

type DbUser struct {
	ID          uint64
	Email       string
	Password    []byte
	Username    string
	ProfileName string
	Balance     int64
}
