package main

import (
	"PR_reviewer_assign_service/internal/handlers"
	"PR_reviewer_assign_service/internal/service"
	"PR_reviewer_assign_service/internal/storage"
	"log"
	"net/http"
	"os"
)

func main() {
	connection := os.Getenv("DATABASE_URL")

	// Create new storage: PostgreSQL
	storage, err := storage.NewPostgresStorage(connection)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}
	defer storage.Close()

	// Create new Service
	svc := service.NewService(storage)

	// Create handlers for requests
	userHandler := handlers.NewUserHandler(svc)
	teamHandler := handlers.NewTeamHandler(svc)
	prHandler := handlers.NewPRHandler(svc)

	// Handle functions
	http.HandleFunc("/team/add", teamHandler.AddTeam)
	http.HandleFunc("/team/get", teamHandler.GetTeam)
	http.HandleFunc("/users/setIsActive", userHandler.SetIsActive)
	http.HandleFunc("/pullRequest/create", prHandler.CreatePR)
	http.HandleFunc("/pullRequest/merge", prHandler.Merge)
	http.HandleFunc("/pullRequest/reassign", prHandler.Reassign)
	http.HandleFunc("/users/getReview", userHandler.GetReview)
	// Additional functions
	http.HandleFunc("/users/statistics", userHandler.GetUserStatistics)
	http.HandleFunc("/users/get", userHandler.GetStatistics)
	http.HandleFunc("/team/statistics", teamHandler.GetTeamStatistics)
	http.HandleFunc("/team/count", teamHandler.GetStatistics)
	http.HandleFunc("/pullRequest/statistics", prHandler.GetStatistics)

	// Handle HealthCheck
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
