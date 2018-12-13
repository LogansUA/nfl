package player

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/logansua/nfl_app/bucket"
	"github.com/logansua/nfl_app/pagination"
	"mime/multipart"
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

type Player struct {
	Model

	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

// Service is a simple CRUD interface for user profiles.
type Service interface {
	PostProfile(ctx context.Context, p Player) (player *Player, err error)
	CreatePlayer(ctx context.Context, p Player) (player *Player, err error)
	GetPlayers(ctx context.Context, paging pagination.Pagination) ([]Player, error)
	GetPlayer(ctx context.Context, id int) (player *Player, err error)
	DeletePlayer(ctx context.Context, id int) error
	UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader multipart.FileHeader) (player *Player, err error)
}

type service struct {
	DBService     *gorm.DB
	BucketService bucket.Service
}

func New(dbService *gorm.DB, bucketService bucket.Service) Service {
	return &service{DBService: dbService, BucketService: bucketService}
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

func (s *service) PostProfile(ctx context.Context, p Player) (player *Player, err error) {
	return &Player{Name: "Hello"}, nil
}

func (s *service) CreatePlayer(ctx context.Context, p Player) (player *Player, err error) {
	s.DBService.Create(&p)

	return &p, nil
}

func (s *service) GetPlayers(ctx context.Context, paging pagination.Pagination) ([]Player, error) {
	var players []Player

	s.DBService.
		Offset(paging.Offset).
		Limit(paging.Limit).
		Order("id ASC").
		Find(&players)

	return players, nil
}

func (s *service) GetPlayer(ctx context.Context, id int) (player *Player, err error) {
	var p Player

	if result := s.DBService.First(&p, id); result.RecordNotFound() {
		return nil, errors.New("player not found")
	}

	return &p, nil
}

func (s *service) DeletePlayer(ctx context.Context, id int) error {
	var player Player

	if s.DBService.First(&player, id).RecordNotFound() {
		return errors.New("player not found")
	}

	s.DBService.Delete(&player)

	return nil
}

func (s *service) UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader multipart.FileHeader) (player *Player, err error) {
	var p Player

	if result := s.DBService.First(&p, id); result.Error != nil {
		return nil, errors.New("player not found")
	}

	url, err := s.BucketService.UploadPlayerAvatar(p.ID, file, &fileHeader)

	if err != nil {
		return nil, err
	}

	p.Avatar = url

	if result := s.DBService.Save(&p); result.Error != nil {
		return nil, result.Error
	}

	return &p, nil
}
