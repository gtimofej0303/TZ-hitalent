package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/gtimofej0303/TZ-hitalent/internal/domain"
	"github.com/gtimofej0303/TZ-hitalent/internal/repository"
)

type DepartmentService interface {
	Create(ctx context.Context, name string, parentID *uint) (*domain.Department, error)
	GetByID(ctx context.Context, id uint) (*domain.Department, error)
	Update(ctx context.Context, id uint, name string, parentID *uint) (*domain.Department, error)
	Delete(ctx context.Context, id uint, mode string, reassignTo *uint) error
	GetTree(ctx context.Context, id uint, depth int, includeEmployees bool) ([]*domain.Department, []*domain.Employee, error)
}

type departmentService struct {
	repo    repository.DepartmentRepository
	empRepo repository.EmployeeRepository
}

func NewDepartmentService(
	repo repository.DepartmentRepository, empRepo repository.EmployeeRepository) DepartmentService {
	return &departmentService{
		repo:    repo,
		empRepo: empRepo,
	}
}

func (s *departmentService) Create(ctx context.Context, name string, parentID *uint) (*domain.Department, error) {
	name = strings.TrimSpace(name)
	
	if (len(name) == 0) || (len(name) > 200) {
		return nil, ErrInvalidName
	}

	if parentID != nil && *parentID == 0 {
		parentID = nil
	}

	//Проверка существования родителя
	if parentID != nil {
		if _, err := s.GetByID(ctx, *parentID); err != nil {
			return nil, ErrParentNotFound
		}
	}

	exists, err := s.repo.ExistsByNameAndParent(ctx, name, parentID)
	if err != nil {
		return nil, fmt.Errorf("check uniqueness error: %w", err)
	}
	if exists {
		return nil, ErrDuplicateName
	}

	dept := &domain.Department{Name: name, ParentID: parentID}
	if err := s.repo.Create(ctx, dept); err != nil {
		return nil, fmt.Errorf("failed to create department: %w", err)
	}
	return dept, nil
}

func (s *departmentService) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	dept, err := s.repo.GetByID(ctx, id)
	if (err != nil) || (dept == nil) {
		return nil, ErrDepartmentNotFound
	}
	return dept, nil
}

func (s *departmentService) Update(ctx context.Context, id uint, name string, parentID *uint) (*domain.Department, error) {
	name = strings.TrimSpace(name)
	
	if len(name) == 0 || len(name) > 200 {
		return nil, ErrInvalidName
	}
	if parentID != nil && *parentID == 0 {
		parentID = nil
	}

	dept, err := s.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	exists, err := s.repo.ExistsByNameAndParent(ctx, name, parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to check uniqueness: %w", err)
	}
	if exists {
		return nil, ErrDuplicateName
	}

	if parentID != nil && (dept.ParentID == nil || *parentID != *dept.ParentID) {
		if _, err := s.GetByID(ctx, *parentID); err != nil {
			return nil, ErrParentNotFound
		}
		if err := s.validateNoCycle(ctx, *parentID, &id); err != nil {
			return nil, err
		}
	}

	dept.Name = name
	dept.ParentID = parentID
	if err := s.repo.Update(ctx, dept); err != nil {
		return nil, fmt.Errorf("failed to update department: %w", err)
	}

	return dept, nil
}

func (s *departmentService) Delete(ctx context.Context, id uint, mode string, reassignTo *uint) error {
	_, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	switch mode {
	case "cascade":
		children, err := s.repo.GetChildren(ctx, id, 999)
		if err != nil {
			return fmt.Errorf("failed to get children: %w", err)
		}

		for _, child := range children {
			if err := s.empRepo.DeleteByDepartmentID(ctx, child.ID); err != nil {
				return fmt.Errorf("failed to delete employees: %w", err)
			}
			if err := s.repo.Delete(ctx, child.ID); err != nil {
				return fmt.Errorf("failed to delete department: %w", err)
			}
		}

		if err := s.empRepo.DeleteByDepartmentID(ctx, id); err != nil {
        	return fmt.Errorf("failed to delete employees of root dept: %w", err)
    	}
		return s.repo.Delete(ctx, id)


	case "reassign":
		if reassignTo == nil || *reassignTo == 0 {
			return ErrInvalidReassignTarget
		}

		if _, err := s.GetByID(ctx, *reassignTo); err != nil {
			return ErrDepartmentNotFound
		}

		if err := s.empRepo.ReassignToDepartment(ctx, id, *reassignTo); err != nil {
			return fmt.Errorf("failed to reassign department employees: %w", err)
		}

		children, err := s.repo.GetChildren(ctx, id, 1)
		if err != nil {
			return fmt.Errorf("failed to get children: %w", err)
		}
		for _, child := range children {
			if child.ID != id {
				fullChild, err := s.GetByID(ctx, child.ID)
				if err != nil {
					return fmt.Errorf("failed to load fullchild %d: %w", child.ID, err)
				}
				fullChild.ParentID = reassignTo
				if err := s.repo.Update(ctx, fullChild); err != nil { // сохраняет ВСЕ поля
					return fmt.Errorf("failed update child %d: %w", child.ID, err)
				}
			}
		}
		return s.repo.Delete(ctx, id)

	default:
		return ErrInvalidMode
	}
	return nil
}

func (s *departmentService) GetTree(ctx context.Context, id uint, depth int, includeEmployees bool) ([]*domain.Department, []*domain.Employee, error) {
	if depth < 1 || depth > 5 {
		depth = 1
	}

	children, err := s.repo.GetChildren(ctx, id, depth)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tree: %w", err)
	}

	var employees []*domain.Employee

	if includeEmployees {
		for _, child := range children {
			childEmp, err := s.empRepo.GetByDepartmentID(ctx, child.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("failed to get %d employees: %w", child.ID, err)
			}
			employees = append(employees, childEmp...)
		}
	}
	return children, employees, nil
}

func (s *departmentService) validateNoCycle(ctx context.Context, parentID uint, excludeID *uint) error {
	if excludeID == nil {
		return nil
	}

	if *excludeID == parentID {
		return ErrSelfParent
	}

	childrenPar, err := s.repo.GetChildren(ctx, parentID, 999)
	if err != nil {
		return err
	}
	for _, ch := range childrenPar {
		if ch.ID == *excludeID {
			return ErrCycleDetected
		}
	}

	childrenEx, err := s.repo.GetChildren(ctx, *excludeID, 999)
	if err != nil {
		return err
	}
	for _, ch := range childrenEx {
		if ch.ID == parentID {
			return ErrCycleDetected
		}
	}

	return nil
}
