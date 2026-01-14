package services

import (
	"context"
	"employee-system/internal/models"
	"employee-system/internal/repositories"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, emp *models.Employee) error
	GetEmployeeByID(ctx context.Context, id int) (*models.Employee, error)
	GetEmployees(ctx context.Context, limit, offset int, departmentID *int) ([]*models.Employee, int, error)
	UpdateEmployee(ctx context.Context, id int, emp *models.Employee) error
	DeleteEmployee(ctx context.Context, id int) error
}

type employeeService struct {
	repo repositories.EmployeeRepository
}

func NewEmployeeService(repo repositories.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

func (s *employeeService) CreateEmployee(ctx context.Context, emp *models.Employee) error {
	return s.repo.Create(ctx, emp)
}

func (s *employeeService) GetEmployeeByID(ctx context.Context, id int) (*models.Employee, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *employeeService) GetEmployees(ctx context.Context, limit, offset int, departmentID *int) ([]*models.Employee, int, error) {
	return s.repo.GetAll(ctx, limit, offset, departmentID)
}

func (s *employeeService) UpdateEmployee(ctx context.Context, id int, emp *models.Employee) error {
	return s.repo.Update(ctx, id, emp)
}

func (s *employeeService) DeleteEmployee(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
