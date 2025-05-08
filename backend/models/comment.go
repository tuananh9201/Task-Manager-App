package models

type Comment struct {
	ID        int    `db:"id" json:"id"`
	CardID    int    `db:"card_id" json:"card_id"`
	UserID    int    `db:"user_id" json:"user_id"`
	Content   string `db:"content" json:"content"`
	CreatedAt string `db:"created_at" json:"created_at"`
}
