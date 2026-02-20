package domain

import "time"

type Employee struct{
	id int
	department_id int
	full_name string //Not null
	position string //Not null
	hired_at time.Time
	created_at time.Time
}