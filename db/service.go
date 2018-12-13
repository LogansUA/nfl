package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func New(dialect, connectionString string) *gorm.DB {
	db, err := gorm.Open(dialect, connectionString)

	if err != nil {
		panic(err)
	}

	return db
}
