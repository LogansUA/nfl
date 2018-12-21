package models

import "time"

type Player struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time

	Name   string
	Avatar string
	TeamID int
	Team   Team `gorm:"foreignkey:TeamID" sql:"type:int REFERENCES teams(id)"`
}
