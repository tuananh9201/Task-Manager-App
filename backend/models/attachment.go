package models

type Attachment struct {
	ID         int    `db:"id" json:"id"`
	CardID     int    `db:"card_id" json:"card_id"`
	FileName   string `db:"filename" json:"filename"`
	URL        string `db:"url" json:"url"`
	UploadedAt string `db:"uploaded_at" json:"uploaded_at"`
}
