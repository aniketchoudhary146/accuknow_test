package models

type UserNotes struct {
	ID      uint   `json:"id"`
	UserId  uint   `json:"user_id"`
	Message string `json:"note"`
}
