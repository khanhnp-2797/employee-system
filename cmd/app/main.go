package main

import (
	"log"
	"net/http"
	"strings"

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

	departmentRepo := repositories.NewDepartmentRepository(db)
	departmentService := services.NewDepartmentService(departmentRepo)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

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

	mux.HandleFunc("/employees/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			employeeHandler.SearchEmployees(w, r)
		}
	})

	mux.HandleFunc("/employees/export", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			employeeHandler.ExportEmployees(w, r)
		}
	})

	mux.HandleFunc("/departments", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			departmentHandler.CreateDepartment(w, r)
		} else if r.Method == http.MethodGet {
			if len(r.URL.Path) == len("/departments") {
				departmentHandler.GetDepartments(w, r)
			} else if strings.Contains(r.URL.Path, "/employees") {
				departmentHandler.GetEmployeesByDepartment(w, r)
			}
		}
	})

	log.Printf("Server running on http://localhost:%s\n", cfg.AppPort)
	if err := http.ListenAndServe(":"+cfg.AppPort, mux); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
