package models

type Card struct {
	ID          int     `db:"id" json:"id"`
	ListID      int     `db:"list_id" json:"list_id"`
	Title       string  `db:"title" json:"title"`
	Description *string `db:"description" json:"description,omitempty"`
	Position    int     `db:"position" json:"position"`
	DueDate     *string `db:"due_date" json:"due_date,omitempty"`
	CreatedAt   string  `db:"created_at" json:"created_at"`
}
