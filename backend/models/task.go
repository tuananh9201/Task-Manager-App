package models

type Task struct {
	ID          int       `db:"id" json:"id"`
	BoardID     int       `db:"board_id" json:"board_id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	Status      string    `db:"status" json:"status"`
	DueDate     string    `db:"due_date" json:"due_date"`
	UserID      int       `db:"user_id" json:"user_id"`
	Comments    []Comment `db:"comments" json:"comments"`
	CreatedAt   string    `db:"created_at" json:"created_at"`
	UpdatedAt   string    `db:"updated_at" json:"updated_at"`
	ListID      int       `db:"list_id" json:"list_id"`
}

// Add this new struct to hold tasks with list information
type TaskWithList struct {
	Task
	ListName string `db:"list_name" json:"list_name"`
}
