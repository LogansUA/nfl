package db

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
)

type DB struct {
	Repository       Repository
	PlayerRepository PlayerRepository
	TeamRepository   TeamRepository
	DB               *gorm.DB
}

// Initialize connection to database
func New() (*DB, error) {
	db, err := gorm.Open("postgres", getConnectionString())

	if err != nil {
		return nil, err
	}

	return &DB{
		Repository:       &BaseRepository{DB: db},
		PlayerRepository: &PlayerTable{DB: db},
		TeamRepository:   &TeamTable{DB: db},
		DB:               db,
	}, nil
}

func getConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
	)
}
