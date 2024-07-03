package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PassportNumber string `json:"passportNumber" gorm:"unique;not null"`
	Surname        string `json:"surname"`
	Name           string `json:"name"`
	Patronymic     string `json:"patronymic"`
	Address        string `json:"address"`
	Tasks          []Task `gorm:"foreignKey:UserID"`
}

type Task struct {
	gorm.Model
	UserID       int       `json:"userId"`
	TaskName     string    `json:"taskName"`
	TaskDuration time.Time `json:"taskDuration"`
	TaskStart    time.Time `json:"taskStart"`
	Status       bool      `json:"status"`
}
