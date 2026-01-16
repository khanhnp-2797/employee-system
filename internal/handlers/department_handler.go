package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"employee-system/internal/models"
	"employee-system/internal/services"
)

type DepartmentHandler struct {
	service services.DepartmentService
}

type DepartmentsListResponse struct {
	Departments []*models.Department `json:"departments"`
}

type DepartmentEmployeesResponse struct {
	DepartmentID int                `json:"departmentId"`
	Employees    []*models.Employee `json:"employees"`
}

func NewDepartmentHandler(service services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateDepartment handler called")

	var req struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "department name is required")
		return
	}

	dept := &models.Department{
		Name: req.Name,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.CreateDepartment(ctx, dept); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dept)
}

func (h *DepartmentHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	log.Println("GetDepartments handler called")

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	departments, err := h.service.GetDepartments(ctx)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if departments == nil {
		departments = []*models.Department{}
	}

	json.NewEncoder(w).Encode(DepartmentsListResponse{Departments: departments})
}

func (h *DepartmentHandler) GetEmployeesByDepartment(w http.ResponseWriter, r *http.Request) {
	log.Println("GetEmployeesByDepartment handler called")

	deptID, err := extractDepartmentIDFromPath(r.URL.Path)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid department id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Verify department exists
	_, err = h.service.GetDepartmentByID(ctx, deptID)
	if err != nil {
		writeError(w, http.StatusNotFound, "department not found")
		return
	}

	// Get employees in department
	employees, err := h.service.GetEmployeesByDepartment(ctx, deptID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if employees == nil {
		employees = []*models.Employee{}
	}

	json.NewEncoder(w).Encode(DepartmentEmployeesResponse{
		DepartmentID: deptID,
		Employees:    employees,
	})
}

func extractDepartmentIDFromPath(path string) (int, error) {
	// Path format: /departments/{id}/employees
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return 0, strconv.ErrSyntax
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}
	return id, nil
}
