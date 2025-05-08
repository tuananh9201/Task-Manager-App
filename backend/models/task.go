package models

type Task struct {
	ID          int    `db:"id" json:"id"`
	BoardID     int    `db:"board_id" json:"board_id"`
	Title       string `db:"title" json:"title"`
	Description string `db:"description" json:"description"`
	Status      string `db:"status" json:"status"`
	DueDate     string `db:"due_date" json:"due_date"`
	UserID      int    `db:"user_id" json:"user_id"`
}
