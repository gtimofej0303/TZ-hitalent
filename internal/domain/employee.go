package domain

import "time"

type Employee struct{
	ID uint 			`gorm:"primaryKey;autoIncrement" json:"id"`
	DepartmentID uint 	`gorm:"column:department_id;not null" json:"department_id"`
	Fullname string 	`gorm:"type:varchar(200);not null;size:200" json:"fullname"`
	Position string 	`gorm:"type:varchar(200);not null;size:200" json:"position"`
	HiredAt *time.Time 	`gorm:"column:hired_at" json:"hired_at"`
	CreatedAt time.Time `json:"created_at"`
}