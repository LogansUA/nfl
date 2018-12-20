package dto

import (
	"time"
)

type TeamDTO struct {
	ID uint `json:"id"`

	Name string `json:"name"`
	Logo string `json:"logo"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
