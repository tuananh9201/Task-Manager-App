package models

type CardLabel struct {
	CardID  int `db:"card_id" json:"card_id"`
	LabelID int `db:"label_id" json:"label_id"`
}
