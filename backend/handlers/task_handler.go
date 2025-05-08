package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"task-app/models"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

type TaskHandler struct {
	DB        *sqlx.DB
	Broadcast chan TaskEvent
}

type TaskEvent struct {
	Type string
	Task models.Task
}

type PaginatedTasksResponse struct {
	Tasks       []models.Task `json:"tasks"`
	CurrentPage int           `json:"current_page"`
	TotalPages  int           `json:"total_pages"`
	TotalItems  int           `json:"total_items"`
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")
	boardID := c.Query("board_id")

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
	fmt.Printf("user_id: %d\n", userID)
	fmt.Printf("board_id: %s\n", boardID)
	// Get total count of tasks
	var totalItems int
	query := "SELECT COUNT(*) FROM tasks WHERE user_id = $1"
	args := []interface{}{userID}
	if boardID != "" {
		query += " AND board_id = $2"
		var boardIDNum int
		boardIDNum, err = strconv.Atoi(boardID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board_id"})
			return
		}
		args = append(args, boardIDNum)
	}

	err = h.DB.Get(&totalItems, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks count"})
		return
	}

	totalPages := (totalItems + limitNum - 1) / limitNum

	// Get paginated tasks
	var tasks []models.Task
	query = "SELECT * FROM tasks WHERE user_id = $1"
	args = []interface{}{userID}
	if boardID != "" {
		query += " AND board_id = $2"
		boardIDNum, _ := strconv.Atoi(boardID)
		args = append(args, boardIDNum)
	}
	query += " ORDER BY id LIMIT $" + strconv.Itoa(len(args)+1) + " OFFSET $" + strconv.Itoa(len(args)+2)
	args = append(args, limitNum, offset)

	err = h.DB.Select(&tasks, query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tasks " + err.Error()})
		return
	}
	if tasks == nil {
		tasks = []models.Task{}
	}
	response := PaginatedTasksResponse{
		Tasks:       tasks,
		CurrentPage: pageNum,
		TotalPages:  totalPages,
		TotalItems:  totalItems,
	}

	c.JSON(http.StatusOK, response)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	taskID := c.Param("id")
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Check if task exists and belongs to user
	userID := c.GetInt("user_id")
	var existingTask models.Task
	err = h.DB.Get(&existingTask, "SELECT * FROM tasks WHERE id = $1 AND user_id = $2", taskIDInt, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or unauthorized"})
		return
	}

	// Bind update data
	var task models.Task
	if bindErr := c.ShouldBindJSON(&task); bindErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
		return
	}

	// Update task
	setClauses := []string{}
	args := []interface{}{}
	argIdx := 1

	if task.Title != "" {
		setClauses = append(setClauses, "title = $"+strconv.Itoa(argIdx))
		args = append(args, task.Title)
		argIdx++
	}

	if task.Description != "" {
		setClauses = append(setClauses, "description = $"+strconv.Itoa(argIdx))
		args = append(args, task.Description)
		argIdx++
	}

	if task.Status != "" {
		setClauses = append(setClauses, "status = $"+strconv.Itoa(argIdx))
		args = append(args, task.Status)
		argIdx++
	}

	if task.DueDate != "" {
		setClauses = append(setClauses, "due_date = $"+strconv.Itoa(argIdx))
		args = append(args, task.DueDate)
		argIdx++
	}

	if len(setClauses) == 0 {
		c.JSON(http.StatusOK, task)
		return
	}

	query := "UPDATE tasks SET " + strings.Join(setClauses, ", ") + " WHERE id = $" + strconv.Itoa(argIdx) + " AND user_id = $" + strconv.Itoa(argIdx+1)
	args = append(args, taskIDInt, userID)

	_, err = h.DB.Exec(query, args...)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update task"})
		return
	}

	// Get updated task
	h.DB.Get(&task, "SELECT * FROM tasks WHERE id = $1", taskIDInt)
	fmt.Printf("Updated task: %+v\n", task.BoardID)
	// Broadcast the updated task
	h.Broadcast <- TaskEvent{Type: "update", Task: task}

	c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	taskID := c.Param("id")
	taskIDInt, err := strconv.Atoi(taskID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	userID := c.GetInt("user_id")

	// Check if task exists and belongs to user
	var task models.Task
	err = h.DB.Get(&task, "SELECT * FROM tasks WHERE id = $1 AND user_id = $2", taskIDInt, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found or unauthorized"})
		return
	}

	// Delete task
	_, err = h.DB.Exec("DELETE FROM tasks WHERE id = $1 AND user_id = $2", taskIDInt, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete task"})
		return
	}

	// Broadcast the deleted task with a specific action type
	h.Broadcast <- TaskEvent{Type: "delete", Task: models.Task{ID: task.ID}}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	var task models.Task
	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetInt("user_id")
	boardID, err := strconv.Atoi(fmt.Sprintf("%v", task.BoardID))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid board_id format"})
		return
	}

	_, err = h.DB.Exec(
		"INSERT INTO tasks (board_id, title, description, status, due_date, user_id) VALUES ($1, $2, $3, $4, $5, $6)",
		boardID, task.Title, task.Description, task.Status, task.DueDate, userID,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create task"})
		return
	}

	h.DB.Get(&task, "SELECT * FROM tasks WHERE title=$1 AND user_id=$2", task.Title, userID)

	// Broadcast the new task to all connected clients
	h.Broadcast <- TaskEvent{Type: "create", Task: task}

	c.JSON(http.StatusOK, task)
}
