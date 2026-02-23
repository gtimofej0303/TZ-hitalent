package domain

import "time"

type Department struct {
	ID         uint 		`gorm:"primaryKey;autoIncrement" json:"id"`
	Name       string		`gorm:"type:varchar(200);not null;size:200" json:"name"`
	ParentID  *uint			`gorm:"column:parent_id" json:"parent_id"`
	CreatedAt time.Time		`json:"created_at"`
}