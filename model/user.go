package model

import "gorm.io/gorm"

// User struct
type User struct {
	gorm.Model
	Name       string     `gorm:"not null;size:50;" validate:"required,min=3,max=50" json:"name"`
	Username   string     `gorm:"uniqueIndex;not null;size:50;" validate:"required,min=3,max=50" json:"username"`
	Email      string     `gorm:"uniqueIndex;not null;size:255;" validate:"required,email" json:"email"`
	Password   string     `gorm:"not null;" validate:"required,min=6,max=50" json:"password"`
	Occupation string     `json:"occupation"`
	About      string     `json:"about"`
	TaskLists  []TaskList `gorm:"foreignKey:UserID" json:"lists"`
}
