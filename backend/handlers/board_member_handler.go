package handlers

import (
	"net/http"

	"task-app/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type BoardMemberHandler struct {
	DB *sqlx.DB
}

// CreateBoardMember adds a new member o a board
func (h *BoardMemberHandler) CreateBoardMember(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id").(uint)

	var input struct {
		BoardID uint   `json:"board_id"`
		UserID  uint   `json:"user_id"`
		Role    string `json:"role"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Check if user is owner of the board
	var board models.Board
	if err := h.DB.Where("id = ? AND owner_id = ?", input.BoardID, userID).First(&board).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Board not found or unauthorized"})
		return
	}

	// Check if the member already exists
	var existingMember models.BoardMember
	if err := models.DB.Where("board_id = ? AND user_id = ?", input.BoardID, input.UserID).First(&existingMember).Error; err == nil {
		c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": "Member already exists"})
		return
	}

	newMember := models.BoardMember{
		BoardID: input.BoardID,
		UserID:  input.UserID,
		Role:    input.Role,
	}

	if err := models.DB.Create(&newMember).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to add member"})
		return
	}

	c.JSON(http.StatusCreated, newMember)
}

// GetBoardMembers retrieves all members of a board
func GetBoardMembers(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id").(uint)
	boardID := c.Param("board_id")

	// Check if user has access to this board
	var boardMember models.BoardMember
	if err := models.DB.Where("board_id = ? AND user_id = ?", boardID, userID).First(&boardMember).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	var members []models.BoardMember
	if err := models.DB.Where("board_id = ?", boardID).Find(&members).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch members"})
		return
	}

	c.JSON(http.StatusOK, members)
}

// DeleteBoardMember removes a member from a board
func DeleteBoardMember(c *gin.Context) {
	session := sessions.Default(c)
	currentUserID := session.Get("user_id").(uint)
	boardID := c.Param("board_id")
	targetUserID := c.Param("user_id")

	// Check if the current user is the owner or admin
	var boardMember models.BoardMember
	if err := models.DB.Where("board_id = ? AND user_id = ? AND (role = ? OR role = ?)", boardID, currentUserID, "owner", "admin").First(&boardMember).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Unauthorized to remove member"})
		return
	}

	if err := models.DB.Where("board_id = ? AND user_id = ?", boardID, targetUserID).Delete(&models.BoardMember{}).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove member"})
		return
	}

	c.JSON(http.StatusNoContent, gin.H{})
}
