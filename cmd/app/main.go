package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"employee-system/config"
	"employee-system/internal/database"
	"employee-system/internal/handlers"
	"employee-system/internal/repositories"
	"employee-system/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg.GetDSN())
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	log.Printf("Connected to database: %s\n", cfg.DBName)

	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo)
	employeeHandler := handlers.NewEmployeeHandler(employeeService)

	mux := http.NewServeMux()

	mux.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			employeeHandler.CreateEmployee(w, r)
		} else if r.Method == http.MethodGet {
			if len(r.URL.Path) == len("/employees") {
				employeeHandler.GetEmployees(w, r)
			} else {
				employeeHandler.GetEmployeeByID(w, r)
			}
		} else if r.Method == http.MethodPut {
			employeeHandler.UpdateEmployee(w, r)
		} else if r.Method == http.MethodDelete {
			employeeHandler.DeleteEmployee(w, r)
		}
	})

	server := &http.Server{
		Addr:         ":" + cfg.AppPort,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		<-sigChan

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
	}()

	log.Printf("Server running on http://localhost:%s\n", cfg.AppPort)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
