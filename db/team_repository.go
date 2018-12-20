package db

import (
	"github.com/jinzhu/gorm"
	"github.com/logansua/nfl_app/models"
	"github.com/logansua/nfl_app/pagination"
)

type TeamRepository interface {
	FindAllAndPaginate(paging pagination.Pagination, out *[]models.Team) error
}

type TeamTable struct {
	DB *gorm.DB
}

func (pt *TeamTable) FindAllAndPaginate(paging pagination.Pagination, out *[]models.Team) error {
	return pt.
		DB.
		Offset(paging.Offset).
		Limit(paging.Limit).
		Find(out).
		Order("id ASC").
		Error
}
