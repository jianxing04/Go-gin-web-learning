package models

import "gorm.io/gorm"

type Employee struct {
	gorm.Model
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
	Dept  string `json:"dept"`
}
