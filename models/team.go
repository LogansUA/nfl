package models

import (
	"github.com/jinzhu/gorm"
)

type Team struct {
	gorm.Model

	Name    string
	Logo    string
	Players []Player
}
