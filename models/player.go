package models

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Player struct {
	gorm.Model

	Name   string
	Avatar string
}

type PlayerDTO struct {
	ID uint `json:"id"`

	Name   string `json:"name"`
	Avatar string `json:"avatar"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	//DeletedAt *time.Time `json:"deleted_at"`
}

func NewDTO(data Player) PlayerDTO {
	return PlayerDTO{
		ID:        data.ID,
		Name:      data.Name,
		Avatar:    data.Avatar,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
		//DeletedAt: data.DeletedAt,
	}
}

func NewModel(data *PlayerDTO) Player {
	return Player{
		//ID:        data.ID,
		Name:   data.Name,
		Avatar: data.Avatar,
		//CreatedAt: data.CreatedAt,
		//UpdatedAt: data.UpdatedAt,
		//DeletedAt: data.DeletedAt,
	}
}
