package models

import (
	"github.com/jinzhu/gorm"
)

type Player struct {
	gorm.Model

	Name   string
	Avatar string
	TeamID int
	Team   Team `gorm:"foreignkey:TeamID"`
}
