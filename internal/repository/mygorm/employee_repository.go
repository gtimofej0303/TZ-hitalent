package mygorm

import (
	"context"

	"github.com/gtimofej0303/TZ-hitalent/internal/domain"
	"github.com/gtimofej0303/TZ-hitalent/internal/repository"
	"gorm.io/gorm"
)

type EmployeeRepo struct {
	db *gorm.DB
}

func NewEmployeeRepository(db *gorm.DB) repository.EmployeeRepository {
	return &EmployeeRepo{db: db}
}

func (r *EmployeeRepo) Create(ctx context.Context, emp *domain.Employee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}

func (r *EmployeeRepo) GetByDepartmentID(ctx context.Context, deptID uint) ([]*domain.Employee, error) {
	var employees []*domain.Employee
	err := r.db.WithContext(ctx).
		Where("department_id = ?", deptID).
		Find(&employees).Error
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *EmployeeRepo) ReassignToDepartment(ctx context.Context, deptID uint, newDeptID uint) error {
	if deptID == newDeptID {
		return nil
	}
	result := r.db.WithContext(ctx).
		Model(&domain.Employee{}).
		Where("department_id = ?", deptID).
		Update("department_id", newDeptID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (r *EmployeeRepo) DeleteByDepartmentID(ctx context.Context, deptID uint) error {
	result := r.db.WithContext(ctx).
		Where("department_id = ?", deptID).
		Delete(&domain.Employee{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}
