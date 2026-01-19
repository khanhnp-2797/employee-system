package services

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"employee-system/internal/models"
	"employee-system/internal/repositories"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, emp *models.Employee) error
	GetEmployeeByID(ctx context.Context, id int) (*models.Employee, error)
	GetEmployees(ctx context.Context, limit, offset int, departmentID *int) ([]*models.Employee, int, error)
	UpdateEmployee(ctx context.Context, id int, emp *models.Employee) error
	DeleteEmployee(ctx context.Context, id int) error
	SearchEmployees(ctx context.Context, keyword string) ([]*models.Employee, error)
	ExportEmployees(ctx context.Context) error
}

type employeeService struct {
	repo repositories.EmployeeRepository
	mu   sync.Mutex
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

func (s *employeeService) SearchEmployees(ctx context.Context, keyword string) ([]*models.Employee, error) {
	return s.repo.Search(ctx, keyword)
}

func (s *employeeService) ExportEmployees(ctx context.Context) error {
	employees, _, err := s.repo.GetAll(ctx, 10000, 0, nil)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.exportJSON(employees); err != nil {
			errChan <- fmt.Errorf("JSON export error: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.exportCSV(employees); err != nil {
			errChan <- fmt.Errorf("CSV export error: %v", err)
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *employeeService) exportJSON(employees []*models.Employee) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create("employees.json")
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(employees)
}

func (s *employeeService) exportCSV(employees []*models.Employee) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, err := os.Create("employees.csv")
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"ID", "Name", "Age", "Position", "DepartmentID", "Salary", "CreatedAt", "UpdatedAt"})

	for _, emp := range employees {
		writer.Write([]string{
			fmt.Sprintf("%d", emp.ID),
			emp.Name,
			fmt.Sprintf("%d", emp.Age),
			emp.Position,
			fmt.Sprintf("%d", emp.DepartmentID),
			fmt.Sprintf("%.2f", emp.Salary),
			emp.CreatedAt.String(),
			emp.UpdatedAt.String(),
		})
	}

	return nil
}
