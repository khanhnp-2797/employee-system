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

type EmployeeHandler struct {
	service services.EmployeeService
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data interface{} `json:"data"`
}

type ListResponse struct {
	TotalCount int                `json:"totalCount"`
	Employees  []*models.Employee `json:"employees"`
}

func NewEmployeeHandler(service services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	log.Println("CreateEmployee handler called")

	var req struct {
		Name         string   `json:"name"`
		Age          *int     `json:"age"`
		Position     *string  `json:"position"`
		DepartmentID int      `json:"departmentId"`
		Salary       *float64 `json:"salary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	employee := &models.Employee{
		Name:         req.Name,
		DepartmentID: req.DepartmentID,
	}

	if req.Age != nil {
		employee.Age = *req.Age
	}
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.Salary != nil {
		employee.Salary = *req.Salary
	}

	if err := h.service.CreateEmployee(r.Context(), employee); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(employee)
}

func (h *EmployeeHandler) GetEmployeeByID(w http.ResponseWriter, r *http.Request) {
	log.Println("GetEmployeeByID handler called")

	id, err := extractIDFromPath(r.URL.Path, "/employees/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid employee id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	emp, err := h.service.GetEmployeeByID(ctx, id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emp)
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	log.Println("GetEmployees handler called")

	limit := 10
	offset := 0
	var departmentID *int

	if l := r.URL.Query().Get("limit"); l != "" {
		if parsedLimit, err := strconv.Atoi(l); err == nil && parsedLimit > 0 {
			limit = parsedLimit
		}
	}

	if o := r.URL.Query().Get("offset"); o != "" {
		if parsedOffset, err := strconv.Atoi(o); err == nil && parsedOffset >= 0 {
			offset = parsedOffset
		}
	}

	if d := r.URL.Query().Get("departmentId"); d != "" {
		if parsedDept, err := strconv.Atoi(d); err == nil {
			departmentID = &parsedDept
		}
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	employees, totalCount, err := h.service.GetEmployees(ctx, limit, offset, departmentID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ListResponse{TotalCount: totalCount, Employees: employees})
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	log.Println("UpdateEmployee handler called")

	id, err := extractIDFromPath(r.URL.Path, "/employees/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid employee id")
		return
	}

	var req struct {
		Name         *string  `json:"name"`
		Age          *int     `json:"age"`
		Position     *string  `json:"position"`
		DepartmentID *int     `json:"departmentId"`
		Salary       *float64 `json:"salary"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	employee := &models.Employee{}

	if req.Name != nil {
		employee.Name = *req.Name
	}
	if req.Age != nil {
		employee.Age = *req.Age
	}
	if req.Position != nil {
		employee.Position = *req.Position
	}
	if req.DepartmentID != nil {
		employee.DepartmentID = *req.DepartmentID
	}
	if req.Salary != nil {
		employee.Salary = *req.Salary
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.UpdateEmployee(ctx, id, employee); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	employee.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(employee)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	log.Println("DeleteEmployee handler called")

	id, err := extractIDFromPath(r.URL.Path, "/employees/")
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid employee id")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	if err := h.service.DeleteEmployee(ctx, id); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

func extractIDFromPath(path, prefix string) (int, error) {
	idStr := strings.TrimPrefix(path, prefix)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, err
	}
	return id, nil
}
