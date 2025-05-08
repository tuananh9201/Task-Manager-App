package models

type BoardMember struct {
	BoardID int    `db:"board_id" json:"board_id"`
	UserID  int    `db:"user_id" json:"user_id"`
	Role    string `db:"role" json:"role"` // e.g. "admin", "member"
}
