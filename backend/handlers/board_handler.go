package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strconv"
	"time"

	"task-app/models"

	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type BoardHandler struct {
	DB *sqlx.DB
}

type PaginatedBoardsResponse struct {
	Boards      []models.Board `json:"boards"`
	CurrentPage int            `json:"current_page"`
	TotalPages  int            `json:"total_pages"`
	TotalItems  int            `json:"total_items"`
}

func generateInviteToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (h *BoardHandler) CreateBoard(c *gin.Context) {
	var board models.Board
	if err := c.ShouldBindJSON(&board); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	fmt.Println("userID : ", userID)
	board.CreatedBy = userID

	tx, err := h.DB.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Insert the board
	var boardID int
	err = tx.QueryRow(
		"INSERT INTO boards (name, description, created_by) VALUES ($1, $2, $3) RETURNING id",
		board.Name, board.Description, board.CreatedBy,
	).Scan(&boardID)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create board" + err.Error()})
		return
	}

	// Add creator as board member with 'owner' role
	_, err = tx.Exec(
		"INSERT INTO board_members (board_id, user_id, role) VALUES ($1, $2, $3)",
		boardID, userID, "owner",
	)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add board member"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	board.ID = boardID
	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) InviteToBoard(c *gin.Context) {
	type InviteRequest struct {
		BoardID int    `json:"board_id"`
		Email   string `json:"email"`
	}

	var req InviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")

	// Check if user has permission to invite
	var role string
	err := h.DB.Get(&role,
		"SELECT role FROM board_members WHERE board_id = $1 AND user_id = $2",
		req.BoardID, userID,
	)

	if err != nil || (role != "owner" && role != "member") {
		c.JSON(http.StatusForbidden, gin.H{"error": "No permission to invite"})
		return
	}

	// Generate invitation token
	token := generateInviteToken()
	expires := time.Now().Add(24 * time.Hour)

	// Create invitation
	_, err = h.DB.Exec(
		"INSERT INTO board_invitations (board_id, email, token, expires_at) VALUES ($1, $2, $3, $4)",
		req.BoardID, req.Email, token, expires,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create invitation"})
		return
	}

	// TODO: Send invitation email

	c.JSON(http.StatusOK, gin.H{"message": "Invitation sent"})
}

func (h *BoardHandler) GetBoards(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, err := strconv.Atoi(page)
	if err != nil || pageNum < 1 {
		pageNum = 1
	}

	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 1 {
		limitNum = 10
	}

	offset := (pageNum - 1) * limitNum
	userID := c.GetInt("user_id")

	// Get total count of boards
	var totalItems int
	err = h.DB.Get(&totalItems, "SELECT COUNT(*) FROM boards b JOIN board_members bm ON b.id = bm.board_id WHERE bm.user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get boards count"})
		return
	}

	totalPages := (totalItems + limitNum - 1) / limitNum

	// Get paginated boards
	var boards []models.Board
	err = h.DB.Select(&boards,
		"SELECT b.* FROM boards b JOIN board_members bm ON b.id = bm.board_id WHERE bm.user_id = $1 ORDER BY b.id LIMIT $2 OFFSET $3",
		userID, limitNum, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get boards " + err.Error()})
		return
	}

	response := PaginatedBoardsResponse{
		Boards:      boards,
		CurrentPage: pageNum,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
	}

	c.JSON(http.StatusOK, response)
}

func (h *BoardHandler) AcceptInvitation(c *gin.Context) {
	type AcceptRequest struct {
		Token string `json:"token"`
	}

	var req AcceptRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")

	// Start transaction
	tx, err := h.DB.Beginx()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
		return
	}

	// Get and validate invitation
	var invitation models.BoardInvitation
	err = tx.Get(&invitation,
		"SELECT * FROM board_invitations WHERE token = $1 AND expires_at > NOW()",
		req.Token,
	)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or expired invitation"})
		return
	}

	// Add user as board member
	_, err = tx.Exec(
		"INSERT INTO board_members (board_id, user_id, role) VALUES ($1, $2, $3)",
		invitation.BoardID, userID, "member",
	)

	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add board member"})
		return
	}

	// Delete the invitation
	_, err = tx.Exec("DELETE FROM board_invitations WHERE token = $1", req.Token)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete invitation"})
		return
	}

	if err := tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to commit transaction"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Successfully joined board"})
}
