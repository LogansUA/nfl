package dto

import (
	"time"
)

type PlayerDTO struct {
	ID uint `json:"id"`

	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	TeamID int    `json:"team_id"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
