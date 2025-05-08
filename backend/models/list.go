package models

type List struct {
	ID        int    `db:"id" json:"id"`
	BoardID   int    `db:"board_id" json:"board_id"`
	Name      string `db:"name" json:"name"`
	Position  int    `db:"position" json:"position"`
	CreatedAt string `db:"created_at" json:"created_at"`
}
