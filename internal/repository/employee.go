package repository

import (
	"context"

	"github.com/gtimofej0303/TZ-hitalent/internal/domain"
)

type EmployeeRepository interface {
	Create(ctx context.Context, emp *domain.Employee) error
	GetByDepartmentID(ctx context.Context, deptID uint) ([]*domain.Employee, error)
	ReassignToDepartment(ctx context.Context, deptID uint, newDeptID uint) error
	DeleteByDepartmentID(ctx context.Context, deptID uint) error
}