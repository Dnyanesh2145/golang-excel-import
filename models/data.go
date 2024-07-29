package models

import (
	"time"

	"gorm.io/gorm"
)

type Records struct {
	gorm.Model

	RecordID  int       `gorm:"primaryKey" json:"record_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Gender    string    `json:"gender"`
	Country   string    `json:"country"`
	Age       int       `json:"age"`
	Date      time.Time `json:"date"`
}
