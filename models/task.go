package models

import (
	"time"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title       string    `gorm:"not null"`
	Description string    
	DueTime     time.Time `gorm:"not null"`
	Completed   bool      `gorm:"default:false"`
	Priority    int       `gorm:"default:1"`
	UserID      uint      `gorm:"not null"`
	User        User      `gorm:"foreignKey:UserID"`
}
