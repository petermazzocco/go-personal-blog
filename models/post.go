package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title   string `validate:"required,min=1,max=100"`
	Content string `validate:"required,min=10,max=500"`
	UserID  uint
	User    User
}
