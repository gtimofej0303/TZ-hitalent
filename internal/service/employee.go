package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gtimofej0303/TZ-hitalent/internal/domain"
	"github.com/gtimofej0303/TZ-hitalent/internal/repository"
)

type EmployeeService interface {
	Create(ctx context.Context, deptID uint, fullName string, position string, hiredAt *time.Time) (*domain.Employee, error)
	GetByDepartmentID(ctx context.Context, deptID uint) ([]*domain.Employee, error)
}

type employeeService struct {
	repo     repository.EmployeeRepository
	deptRepo repository.DepartmentRepository
}

func NewEmployeeService(
	repo repository.EmployeeRepository,
	deptRepo repository.DepartmentRepository,
) EmployeeService {
	return &employeeService{
		repo:     repo,
		deptRepo: deptRepo,
	}
}

func (s *employeeService) Create(ctx context.Context, deptID uint, fullName string, position string, hiredAt *time.Time) (*domain.Employee, error){
	fullName = strings.TrimSpace(fullName)
	position = strings.TrimSpace(position)

	if len(fullName) == 0 || len(fullName) > 200 {
		return nil, ErrInvalidFullName
	}

	if len(position) == 0 || len(position) > 200 {
		return nil, ErrInvalidPosition
	}

	dept, err := s.deptRepo.GetByID(ctx, deptID)
	if err != nil || dept == nil {
		return nil, ErrDepartmentNotFound
	}

	emp := &domain.Employee{
		DepartmentID: deptID,
		Fullname:     fullName,
		Position:     position,
		HiredAt:      hiredAt,
	}

	if err := s.repo.Create(ctx, emp); err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}
	return emp, nil
}


func (s *employeeService) GetByDepartmentID(ctx context.Context, deptID uint) ([]*domain.Employee, error) {
	dept, err := s.deptRepo.GetByID(ctx, deptID)
	if err != nil || dept == nil {
		return nil, ErrDepartmentNotFound
	}

	employees, err := s.repo.GetByDepartmentID(ctx, deptID)
	if err != nil {
		return nil, fmt.Errorf("failed to get employees: %w", err)
	}
	return employees, nil
}