package model

import (
	"gorm.io/gorm"
	"time"
)

// Goal struct
type Goal struct {
	gorm.Model
	UserID          uint            `gorm:"not null" json:"user_id"`
	Name            string          `gorm:"not null;size:255" json:"name"`
	Deadline        time.Time       `json:"deadline"`
	Subgoals        []Subgoal       `gorm:"foreignKey:GoalID" json:"subgoals"`
	Habits          []Habit         `gorm:"foreignKey:GoalID" json:"habits"`
	Description     string          `json:"description"`
	What            string          `json:"what"`
	HowMuch         string          `json:"how_much"`
	Resources       string          `json:"resources"`
	Alignment       string          `json:"alignment"`
	Completed       bool            `gorm:"default:false" json:"completed"`
	SubgoalProgress map[string]bool `gorm:"-" json:"subgoal_progress"` // Not stored in DB, but handled in code
}

// Subgoal struct
type Subgoal struct {
	gorm.Model
	GoalID    uint   `gorm:"not null" json:"goal_id"`
	Name      string `gorm:"not null;size:255" json:"name"`
	Completed bool   `gorm:"default:false" json:"completed"`
}

// Habit struct
type Habit struct {
	gorm.Model
	GoalID    uint   `gorm:"not null" json:"goal_id"`
	Name      string `gorm:"not null;size:255" json:"name"`
	Frequency string `json:"frequency"`
}
