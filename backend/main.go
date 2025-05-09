package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"task-app/handlers"

	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `db:"id" json:"id"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"-"`
}

// Board struct moved to models/board.go

type App struct {
	DB        *sqlx.DB
	Clients   map[*websocket.Conn]bool
	ClientsMu sync.Mutex
	Broadcast chan handlers.TaskEvent
	JWTSecret string
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (app *App) register(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hashed, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashed)
	_, err := app.DB.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func (app *App) login(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var user User
	err := app.DB.Get(&user, "SELECT * FROM users WHERE email=$1", input.Email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user_id": user.ID})
	tokenString, _ := token.SignedString([]byte(app.JWTSecret))
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (app *App) handleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}
	app.ClientsMu.Lock()
	app.Clients[conn] = true
	app.ClientsMu.Unlock()

	defer func() {
		app.ClientsMu.Lock()
		delete(app.Clients, conn)
		app.ClientsMu.Unlock()
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}
	}
}

func (app *App) broadcastTasks() {
	for task := range app.Broadcast {
		data, _ := json.Marshal(task)
		app.ClientsMu.Lock()
		for client := range app.Clients {
			err := client.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println("WebSocket write error:", err)
				client.Close()
				delete(app.Clients, client)
			}
		}
		app.ClientsMu.Unlock()
	}
}

func jwtMiddleware(app *App) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString[:7] == "Bearer " {
			tokenString = tokenString[7:]
		}
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(app.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Next()
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
func runMigrations(db *sqlx.DB) error {
	migration, err := os.ReadFile("migrations/001_create_board_tables.sql")
	if err != nil {
		return err
	}

	_, err = db.Exec(string(migration))
	return err
}
func main() {
	db, err := sqlx.Connect("postgres", "user=admin password=admin123 dbname=task_manager sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	if err := runMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	app := &App{
		DB:        db,
		Clients:   make(map[*websocket.Conn]bool),
		Broadcast: make(chan handlers.TaskEvent),
		JWTSecret: "your-secret-key",
	}

	boardHandler := &handlers.BoardHandler{DB: db}

	r := gin.Default()
	r.Use(corsMiddleware()) // Add CORS middleware
	r.POST("/register", app.register)
	r.POST("/login", app.login)
	protected := r.Group("/api", jwtMiddleware(app))
	taskHandler := &handlers.TaskHandler{
		DB:        db,
		Broadcast: app.Broadcast,
	}
	protected.POST("/tasks", taskHandler.CreateTask)
	protected.GET("/tasks", taskHandler.GetTasks)
	protected.PATCH("/tasks/:id", taskHandler.UpdateTask)
	protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

	// Board routes
	protected.POST("/boards", boardHandler.CreateBoard)
	protected.GET("/boards", boardHandler.GetBoards)
	protected.POST("/boards/invite", boardHandler.InviteToBoard)

	r.GET("/ws", app.handleWebSocket)

	go app.broadcastTasks()

	r.Run(":8080")
}
