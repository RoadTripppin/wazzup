package interfaces

import "github.com/jinzhu/gorm"

type User struct {
	gorm.Model
	Name string
	Email string
	Password string
}

type ResponseUser struct {
	ID uint
	Name string
	Email string
}

type Validation struct {
	Value string
	Valid string
}