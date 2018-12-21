package models

import (
	"time"
)

type Team struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name    string
	Logo    string
	Players []Player `gorm:"foreignkey:TeamID"`
}
