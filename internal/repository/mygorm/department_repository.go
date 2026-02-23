package mygorm

import (
	"context"
	"errors"

	"github.com/gtimofej0303/TZ-hitalent/internal/domain"
	"github.com/gtimofej0303/TZ-hitalent/internal/repository"
	"gorm.io/gorm"
)

type DepartmentRepo struct {
	db *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) repository.DepartmentRepository {
	return &DepartmentRepo{db: db}
}

func (r *DepartmentRepo) Create(ctx context.Context, dept *domain.Department) error {
	return r.db.WithContext(ctx).Create(dept).Error
}

func (r *DepartmentRepo) GetByID(ctx context.Context, id uint) (*domain.Department, error) {
	var dept domain.Department //Это будем возвращать
	err := r.db.WithContext(ctx).First(&dept, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil //Нет данных
		}
		return nil, err
	}
	return &dept, nil
}

func (r *DepartmentRepo) Update(ctx context.Context, dept *domain.Department) error {
	return r.db.WithContext(ctx).Save(dept).Error
}

func (r *DepartmentRepo) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&domain.Department{}, id).Error
}

func (r *DepartmentRepo) GetChildren(ctx context.Context, parentID uint, depth int) ([]*domain.Department, error){
	if depth <= 0{
		return nil, nil
	}

	var children []*domain.Department

	query := `
        WITH RECURSIVE department_tree AS (
            SELECT id, name, parent_id, created_at, 1 as level
            FROM departments WHERE parent_id = ?
            UNION ALL
            SELECT d.id, d.name, d.parent_id, d.created_at, dt.level + 1
            FROM departments d
            INNER JOIN department_tree dt ON d.parent_id = dt.id
            WHERE dt.level < ?
        )
        SELECT id, name, parent_id, created_at
        FROM department_tree
        ORDER BY level, id
    `
	err := r.db.WithContext(ctx).Raw(query, parentID, depth).Scan(&children).Error
    if err != nil {
        return nil, err
    }
    return children, nil
}

func (r *DepartmentRepo) ExistsByNameAndParent(ctx context.Context, name string, parentID *uint) (bool, error) {
	var count int64
	
	if parentID != nil {
		// name = ? AND parent_id = ?
		err := r.db.WithContext(ctx).
			Model(&domain.Department{}).
			Where("name = ? AND parent_id = ?", name, *parentID).
			Count(&count).Error
		if err != nil {
			return false, err
		}
	} else {
		// name = ? AND parent_id IS NULL
		err := r.db.WithContext(ctx).
			Model(&domain.Department{}).
			Where("name = ? AND parent_id IS NULL", name).
			Count(&count).Error
		if err != nil {
			return false, err
		}
	}
	
	return count > 0, nil
}