package models

type Label struct {
	ID      int    `db:"id" json:"id"`
	BoardID int    `db:"board_id" json:"board_id"`
	Name    string `db:"name" json:"name"`
	Color   string `db:"color" json:"color"`
}
