package dto

import (
	"time"

	"github.com/gtimofej0303/org-structure-api/internal/domain"
)

type DepartmentTree struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	ParentID  *uint             `json:"parent_id"`
	CreatedAt time.Time         `json:"created_at"`
	Children  []*DepartmentTree `json:"children"`
	Employees []*EmployeeList   `json:"employees,omitempty"`
}

type EmployeeList struct {
	Fullname  string    `json:"fullname"`
	CreatedAt time.Time `json:"created_at"`
}

// ---------------------------------------------
type CreateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type UpdateDepartmentRequest struct {
	Name     string `json:"name"`
	ParentID *uint  `json:"parent_id"`
}

type DepartmentResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

func BuildDepartmentTree(rootID uint, departments []*domain.Department, employees []*domain.Employee) *DepartmentTree {
	empByDept := make(map[uint][]*EmployeeList)
	for _, e := range employees {
		empByDept[e.DepartmentID] = append(empByDept[e.DepartmentID], &EmployeeList{
			Fullname:  e.Fullname,
			CreatedAt: e.CreatedAt,
		})
	}

	nodes := make(map[uint]*DepartmentTree)
	for _, d := range departments {
		nodes[d.ID] = &DepartmentTree{
			ID:        d.ID,
			Name:      d.Name,
			ParentID:  d.ParentID,
			CreatedAt: d.CreatedAt,
			Children:  []*DepartmentTree{},
			Employees: empByDept[d.ID],
		}
	}

	var roots []*DepartmentTree
	for _, d := range departments {
		node := nodes[d.ID]
		if *d.ParentID == rootID {
			roots = append(roots, node)
		} else {
			if parent, ok := nodes[*d.ParentID]; ok {
				parent.Children = append(parent.Children, node)
			}
		}
	}

	return &DepartmentTree{
		ID:       rootID,
		Children: roots,
	}
}
