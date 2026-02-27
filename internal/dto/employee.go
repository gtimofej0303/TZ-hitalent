package dto

import "time"

type CreateEmployeeRequest struct {
	Fullname     string     `json:"fullname"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
}

type EmployeeResponse struct {
	ID           uint       `json:"id"`
	DepartmentID uint       `json:"department_id"`
	Fullname     string     `json:"fullname"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at"`
}