package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	PassportNumber string `json:"passportNumber" gorm:"unique;not null;column:passportNumber"`
	Surname        string `json:"surname" gorm:"column:surname"`
	Name           string `json:"name" gorm:"column:name"`
	Patronymic     string `json:"patronymic" gorm:"column:patronymic"`
	Address        string `json:"address" gorm:"column:address"`
	Tasks          []Task `gorm:"foreignKey:userId;constraint:OnDelete:CASCADE"`
}

type Task struct {
	gorm.Model
	UserID      uint        `json:"userId" gorm:"column:userId"`
	Description string      `json:"description" gorm:"column:description"`
	Start       time.Time   `json:"start" gorm:"column:start"`
	Status      bool        `json:"status" gorm:"column:status"`
	TimeSpents  []TimeSpent `gorm:"foreignKey:taskId;constraint:OnDelete:CASCADE"`
}

type TimeSpent struct {
	gorm.Model
	TaskID        uint      `json:"taskId" gorm:"column:taskId"`
	BeginInterval time.Time `json:"start" gorm:"column:begin_interval"`
	EndInterval   time.Time `json:"end" gorm:"column:end_interval"`
}
