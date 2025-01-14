package model

import "gorm.io/gorm"

// List struct
type TaskList struct {
	gorm.Model
	UserID uint   `gorm:"not null" json:"user_id"`
	Name   string `gorm:"not null" json:"name"`
	Tasks  []Task `gorm:"foreignKey:TaskListID" json:"tasks"`
}

// Task struct
type Task struct {
	gorm.Model
	TaskListID uint   `gorm:"not null" json:"task_list_id"`
	Text       string `gorm:"not null" json:"text"`
	Completed  bool   `gorm:"default:false" json:"completed"`
}
