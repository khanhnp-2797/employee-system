package services

import (
	"context"
	"employee-system/internal/models"
	"employee-system/internal/repositories"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, dept *models.Department) error
	GetDepartmentByID(ctx context.Context, id int) (*models.Department, error)
	GetDepartments(ctx context.Context) ([]*models.Department, error)
	GetEmployeesByDepartment(ctx context.Context, departmentID int) ([]*models.Employee, error)
}

type departmentService struct {
	repo repositories.DepartmentRepository
}

func NewDepartmentService(repo repositories.DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

func (s *departmentService) CreateDepartment(ctx context.Context, dept *models.Department) error {
	return s.repo.Create(ctx, dept)
}

func (s *departmentService) GetDepartmentByID(ctx context.Context, id int) (*models.Department, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *departmentService) GetDepartments(ctx context.Context) ([]*models.Department, error) {
	return s.repo.GetAll(ctx)
}

func (s *departmentService) GetEmployeesByDepartment(ctx context.Context, departmentID int) ([]*models.Employee, error) {
	return s.repo.GetEmployeesByDepartment(ctx, departmentID)
}
