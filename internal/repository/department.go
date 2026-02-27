package repository

import (
	"context"

	"github.com/gtimofej0303/org-structure-api/internal/domain"
)

type DepartmentRepository interface {
	Create(ctx context.Context, dept *domain.Department) error
	GetByID(ctx context.Context, id uint) (*domain.Department, error)
	Update(ctx context.Context, dept *domain.Department) error
	Delete(ctx context.Context, id uint) error
	GetChildren(ctx context.Context, parentID uint, depth int) ([]*domain.Department, error)
	ExistsByNameAndParent(ctx context.Context, name string, parentID *uint) (bool, error)
}
