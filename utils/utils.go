package utils

import (
	"github.com/satori/go.uuid"
)

func RandToken() string {
	return uuid.Must(uuid.NewV4()).String()
}
