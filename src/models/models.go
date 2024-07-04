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
	UserID      uint      `json:"userId"`
	Description string    `json:"description"`
	Duration    time.Time `json:"duration"`
	Start       time.Time `json:"start"`
	Status      bool      `json:"status"`
}
