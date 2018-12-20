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
	CreatePlayer(ctx context.Context, player *models.PlayerDTO) error
	// Get list of players
	GetPlayers(ctx context.Context, paging pagination.Pagination, players *[]models.PlayerDTO) error
	// Get single player by ID
	GetPlayer(ctx context.Context, id int, player *models.PlayerDTO) error
	// Delete player by ID
	DeletePlayer(ctx context.Context, id int) error
	// Upload player avatar by ID
	UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader, p *models.PlayerDTO) error
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

func (s *service) CreatePlayer(ctx context.Context, player *models.PlayerDTO) error {
	p := models.NewModel(player)
	err := s.DB.Repository.Create(&p)

	(*player).ID = p.ID
	(*player).CreatedAt = p.CreatedAt
	(*player).UpdatedAt = p.UpdatedAt

	return err
}

func (s *service) GetPlayers(ctx context.Context, paging pagination.Pagination, players *[]models.PlayerDTO) error {
	var p []models.Player

	err := s.DB.
		PlayerRepository.
		FindAllAndPaginate(paging, &p)

	*players = make([]models.PlayerDTO, len(p))

	for key, value := range p {
		(*players)[key] = models.NewDTO(value)
	}

	return err
}

func (s *service) GetPlayer(ctx context.Context, id int, player *models.PlayerDTO) error {
	var p models.Player

	err := s.DB.Repository.FindById(&p, id)

	*player = models.NewDTO(p)

	return err
}

func (s *service) DeletePlayer(ctx context.Context, id int) error {
	var p models.Player

	err := s.DB.Repository.FindById(&p, id)

	if err != nil {
		return err
	}

	err = s.DB.Repository.Delete(&p, id)

	return err
}

func (s *service) UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader, player *models.PlayerDTO) error {
	var p models.Player

	if err := s.DB.Repository.FindById(&p, id); err != nil {
		return errors.New("player not found")
	}

	name, err := s.BucketService.UploadPlayerAvatar(ctx, p.ID, file, fileHeader)

	if err != nil {
		return err
	}

	p.Avatar = name

	err = s.DB.Repository.Save(&p)

	*player = models.NewDTO(p)

	return err
}
