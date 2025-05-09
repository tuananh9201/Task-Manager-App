package models

type Board struct {
	ID        int    `db:"id" json:"id"`
	Name      string `db:"name" json:"name"`
	OwnerID   int    `db:"owner_id" json:"owner_id"`
	IsPublic  bool   `db:"is_public" json:"is_public"`
	CreatedAt string `db:"created_at" json:"created_at"`
}
