package player

import (
	"context"
	"errors"
	"github.com/logansua/nfl_app/bucket"
	"github.com/logansua/nfl_app/db"
	"github.com/logansua/nfl_app/models"
	"github.com/logansua/nfl_app/pagination"
	"mime/multipart"
)

// Service is a simple CRUD interface for players.
type Service interface {
	// Create player
	CreatePlayer(ctx context.Context, p models.Player) (*models.Player, error)
	// Get list of players
	GetPlayers(ctx context.Context, paging pagination.Pagination, players *[]models.Player) error
	// Get single player by ID
	GetPlayer(ctx context.Context, id int, player *models.Player) error
	// Delete player by ID
	DeletePlayer(ctx context.Context, id int, player *models.Player) error
	// Upload player avatar by ID
	UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader, p *models.Player) error
}

type service struct {
	DB            *db.DB
	BucketService bucket.Service
}

func New(dbService *db.DB, bucketService bucket.Service) Service {
	return &service{DB: dbService, BucketService: bucketService}
}

var (
	ErrInconsistentIDs = errors.New("inconsistent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

func (s *service) CreatePlayer(ctx context.Context, p models.Player) (*models.Player, error) {
	err := s.DB.Repository.Create(&p)

	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (s *service) GetPlayers(ctx context.Context, paging pagination.Pagination, players *[]models.Player) error {
	var p []models.Player

	err := s.DB.
		PlayerRepository.
		FindAllAndPaginate(paging, &p)

	*players = make([]models.Player, len(p))

	for key, value := range p {
		(*players)[key] = value
	}

	return err
}

func (s *service) GetPlayer(ctx context.Context, id int, player *models.Player) error {
	err := s.DB.Repository.FindById(&player, id)

	return err
}

func (s *service) DeletePlayer(ctx context.Context, id int, player *models.Player) error {
	err := s.DB.Repository.FindById(&player, id)

	if err != nil {
		return err
	}

	err = s.DB.Repository.Delete(&player, id)

	return err
}

func (s *service) UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader, p *models.Player) error {
	if err := s.DB.Repository.FindById(&p, id); err != nil {
		return errors.New("player not found")
	}

	name, err := s.BucketService.UploadPlayerAvatar(ctx, p.ID, file, fileHeader)

	if err != nil {
		return err
	}

	p.Avatar = name

	err = s.DB.Repository.Save(&p)

	return err
}
