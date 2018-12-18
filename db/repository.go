package db

import "github.com/jinzhu/gorm"

type Repository interface {
	FindById(model interface{}, id int) error
	FindAll(model interface{}) error
	Delete(model interface{}, id int) error
	Create(model interface{}) error
	Save(model interface{}) error
}

type BaseRepository struct {
	DB *gorm.DB
}

func (r *BaseRepository) FindById(model interface{}, id int) error {
	return r.DB.First(model, id).Error
}

func (r *BaseRepository) FindAll(model interface{}) error {
	return r.DB.Find(model).Error
}

func (r *BaseRepository) Delete(model interface{}, id int) error {
	return r.DB.Delete(model, id).Error
}

func (r *BaseRepository) Create(model interface{}) error {
	return r.DB.Create(model).Error
}

func (r *BaseRepository) Save(model interface{}) error {
	return r.DB.Save(model).Error
}
