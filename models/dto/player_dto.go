package dto

import (
	"time"
)

type PlayerDTO struct {
	ID uint `json:"id"`

	Name   string `json:"name"`
	Avatar string `json:"avatar"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
