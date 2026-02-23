package dto

import "time"

type DepartmentTree struct {
	ID        	uint				`json:"id"`
	Name      	string				`json:"name"`
	ParentID  	*uint				`json:"parent_id"`
	CreatedAT 	time.Time			`json:"created_at"`
	Children 	[]*DepartmentTree	`json:"children"`
	Employees	[]*EmployeeList		`json:"employees,omitempty"`
}

type EmployeeList struct {
	Fullname 	string		`json:"fullname"`
	CreatedAt 	time.Time	`json:"created_at"`

}