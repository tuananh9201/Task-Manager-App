package models

type Board struct {
	ID          int    `db:"id" json:"id"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	CreatedBy   int    `db:"created_by" json:"created_by"`
	CreatedAt   string `db:"created_at" json:"created_at"`
}

type BoardMember struct {
	ID      int    `db:"id" json:"id"`
	BoardID int    `db:"board_id" json:"board_id"`
	UserID  int    `db:"user_id" json:"user_id"`
	Role    string `db:"role" json:"role"` // e.g., "owner", "member"
}

type BoardInvitation struct {
	ID        int    `db:"id" json:"id"`
	BoardID   int    `db:"board_id" json:"board_id"`
	Email     string `db:"email" json:"email"`
	Token     string `db:"token" json:"token"`
	ExpiresAt string `db:"expires_at" json:"expires_at"`
}
