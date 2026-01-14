package repositories

import (
	"context"
	"database/sql"
	"employee-system/internal/models"
	"fmt"
)

type EmployeeRepository interface {
	Create(ctx context.Context, emp *models.Employee) error
	GetByID(ctx context.Context, id int) (*models.Employee, error)
	GetAll(ctx context.Context, limit, offset int, departmentID *int) ([]*models.Employee, int, error)
	Update(ctx context.Context, id int, emp *models.Employee) error
	Delete(ctx context.Context, id int) error
}

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) Create(ctx context.Context, emp *models.Employee) error {
	query := `INSERT INTO employees (name, age, position, department_id, salary) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, 
		emp.Name, emp.Age, emp.Position, emp.DepartmentID, emp.Salary).
		Scan(&emp.ID, &emp.CreatedAt, &emp.UpdatedAt)

	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id int) (*models.Employee, error) {
	query := `SELECT id, name, age, position, department_id, salary, created_at, updated_at 
			FROM employees WHERE id = $1`

	emp := &models.Employee{}
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&emp.ID, &emp.Name, &emp.Age, &emp.Position, 
		&emp.DepartmentID, &emp.Salary, &emp.CreatedAt, &emp.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee not found")
	}
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}
	
	return emp, nil
}

func (r *employeeRepository) GetAll(ctx context.Context, limit, offset int, departmentID *int) ([]*models.Employee, int, error) {
	query := `SELECT id, name, age, position, department_id, salary, created_at, updated_at 
			FROM employees`
	countQuery := `SELECT COUNT(*) FROM employees`
	
	args := []interface{}{}
	argIndex := 1

	if departmentID != nil {
		query += fmt.Sprintf(` WHERE department_id = $%d`, argIndex)
		countQuery += fmt.Sprintf(` WHERE department_id = $%d`, argIndex)
		args = append(args, *departmentID)
		argIndex++
	}

	query += fmt.Sprintf(` LIMIT $%d OFFSET $%d`, argIndex, argIndex+1)
	args = append(args, limit, offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf(err.Error())
	}
	defer rows.Close()

	var employees []*models.Employee
	for rows.Next() {
		emp := &models.Employee{}
		if err := rows.Scan(&emp.ID, &emp.Name, &emp.Age, &emp.Position, 
			&emp.DepartmentID, &emp.Salary, &emp.CreatedAt, &emp.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf(err.Error())
		}
		employees = append(employees, emp)
	}

	var totalCount int
	if err := r.db.QueryRowContext(ctx, countQuery, args[:argIndex-1]...).Scan(&totalCount); err != nil {
		return nil, 0, fmt.Errorf(err.Error())
	}

	return employees, totalCount, nil
}

func (r *employeeRepository) Update(ctx context.Context, id int, emp *models.Employee) error {
	query := `UPDATE employees 
			SET name = $1, age = $2, position = $3, department_id = $4, salary = $5, 
				updated_at = CURRENT_TIMESTAMP 
			WHERE id = $6 
			RETURNING updated_at`

	err := r.db.QueryRowContext(ctx, query, 
		emp.Name, emp.Age, emp.Position, emp.DepartmentID, emp.Salary, id).
		Scan(&emp.UpdatedAt)

	if err == sql.ErrNoRows {
		return fmt.Errorf("employee not found")
	}
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	
	return nil
}

func (r *employeeRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM employees WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}
