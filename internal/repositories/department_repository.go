package repositories

import (
	"context"
	"database/sql"
	"employee-system/internal/models"
	"fmt"
)

type DepartmentRepository interface {
	Create(ctx context.Context, dept *models.Department) error
	GetByID(ctx context.Context, id int) (*models.Department, error)
	GetAll(ctx context.Context) ([]*models.Department, error)
	GetEmployeesByDepartment(ctx context.Context, departmentID int) ([]*models.Employee, error)
}

type departmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

func (r *departmentRepository) Create(ctx context.Context, dept *models.Department) error {
	query := `INSERT INTO departments (name)
			VALUES ($1)
			RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, dept.Name).
		Scan(&dept.ID, &dept.CreatedAt, &dept.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create department: %v", err)
	}
	return nil
}

func (r *departmentRepository) GetByID(ctx context.Context, id int) (*models.Department, error) {
	query := `SELECT id, name, created_at, updated_at
			FROM departments WHERE id = $1`

	dept := &models.Department{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&dept.ID, &dept.Name, &dept.CreatedAt, &dept.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("department not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get department: %v", err)
	}

	return dept, nil
}

func (r *departmentRepository) GetAll(ctx context.Context) ([]*models.Department, error) {
	query := `SELECT id, name, created_at, updated_at
			FROM departments ORDER BY id ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get departments: %v", err)
	}
	defer rows.Close()

	var departments []*models.Department
	for rows.Next() {
		dept := &models.Department{}
		if err := rows.Scan(&dept.ID, &dept.Name, &dept.CreatedAt, &dept.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan department: %v", err)
		}
		departments = append(departments, dept)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating departments: %v", err)
	}

	return departments, nil
}

// GetEmployeesByDepartment retrieves all employees in a department using a JOIN query
func (r *departmentRepository) GetEmployeesByDepartment(ctx context.Context, departmentID int) ([]*models.Employee, error) {
	query := `SELECT e.id, e.name, e.age, e.position, e.department_id, e.salary, e.created_at, e.updated_at
			FROM employees e
			INNER JOIN departments d ON e.department_id = d.id
			WHERE e.department_id = $1
			ORDER BY e.id ASC`

	rows, err := r.db.QueryContext(ctx, query, departmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees by department: %v", err)
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		emp := &models.Employee{}
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Age, &emp.Position,
			&emp.DepartmentID, &emp.Salary, &emp.CreatedAt, &emp.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan employee: %v", err)
		}
		employees = append(employees, emp)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating employees: %v", err)
	}

	return employees, nil
}
