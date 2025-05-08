package models

type User struct {
	ID           int    `db:"id" json:"id"`
	Email        string `db:"email" json:"email"`
	HashPassword string `db:"password_hash" json:"-"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	Name         string `db:"name" json:"name"`
}
