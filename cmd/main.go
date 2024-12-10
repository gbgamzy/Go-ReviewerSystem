package main

import (
	"ReviewerService/pkg/handlers"
	//"ReviewerService/pkg/models"
	"ReviewerService/pkg/repository"
	"database/sql"
	//"encoding/json"
	"fmt"
	"log"
	"net/http"
	//"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var (
	db           *sql.DB
	userRepo     *repository.UserRepository
	taskRepo     *repository.TaskRepository
	approverRepo *repository.ApproverRepository
	commentRepo  *repository.CommentRepository
	approvalRepo *repository.ApprovalRepository
)

func initDB() {
	connStr := "postgres://client:password@localhost:5432/ReviewerDB?sslmode=disable"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Connected to the database successfully!")
}

func main() {
	// Initialize the database
	initDB()

	// Initialize repositories
	userRepo = &repository.UserRepository{DB: db}
	taskRepo = &repository.TaskRepository{DB: db}
	approverRepo = &repository.ApproverRepository{DB: db}
	commentRepo = &repository.CommentRepository{DB: db}
	approvalRepo = &repository.ApprovalRepository{DB: db}

	userHandler := handlers.UserHandler{UserRepo: userRepo}
	taskHandler := handlers.TaskHandler{TaskRepo: taskRepo}
	approverHandler := &handlers.ApproverHandler{
		ApproverRepo: approverRepo,
		ApprovalRepo: approvalRepo,
		UserRepo:     userRepo,
		TaskRepo:     taskRepo,
	}
	commentHandler := handlers.CommentHandler{CommentRepo: commentRepo}

	// Set up routes
	r := mux.NewRouter()

	r.HandleFunc("/users/login", userHandler.LoginHandler).Methods("POST")
	r.HandleFunc("/tasks", taskHandler.CreateTaskHandler).Methods("POST")
	r.HandleFunc("/tasks/{task_id}/publish", taskHandler.PublishTaskHandler).Methods("PUT")
	r.HandleFunc("/tasks/{task_id}/approvers", approverHandler.AssignApproversHandler).Methods("POST")
	r.HandleFunc("/tasks/{task_id}/comments", commentHandler.AddCommentHandler).Methods("POST")
	r.HandleFunc("/tasks/{task_id}/approve", approverHandler.ApproveTaskHandler).Methods("POST")
	r.HandleFunc("/users", userHandler.GetUsersHandler).Methods("GET")
	r.Use(loggingMiddleware)

	// Start the server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
