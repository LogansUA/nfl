package db

import (
	"github.com/jinzhu/gorm"
	"github.com/logansua/nfl_app/models"
	"github.com/logansua/nfl_app/pagination"
)

type PlayerRepository interface {
	FindAllAndPaginate(paging pagination.Pagination, out *[]models.Player) error
}

type PlayerTable struct {
	DB *gorm.DB
}

func (pt *PlayerTable) FindAllAndPaginate(paging pagination.Pagination, out *[]models.Player) error {
	return pt.
		DB.
		Offset(paging.Offset).
		Limit(paging.Limit).
		Find(out).
		Error
}
