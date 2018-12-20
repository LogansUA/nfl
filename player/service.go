package player

import (
	"context"
	"errors"
	"github.com/logansua/nfl_app/bucket"
	"github.com/logansua/nfl_app/db"
	"github.com/logansua/nfl_app/models"
	"github.com/logansua/nfl_app/models/dto"
	"github.com/logansua/nfl_app/pagination"
	"mime/multipart"
)

// Service is a simple CRUD interface for players.
type Service interface {
	// Create player
	CreatePlayer(ctx context.Context, player *dto.PlayerDTO) error
	// Get list of players
	GetPlayers(ctx context.Context, paging pagination.Pagination, players *[]dto.PlayerDTO) error
	// Get single player by ID
	GetPlayer(ctx context.Context, id int, player *dto.PlayerDTO) error
	// Delete player by ID
	DeletePlayer(ctx context.Context, id int) error
	// Upload player avatar by ID
	UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader, player *dto.PlayerDTO) error
}

type service struct {
	DB            *db.DB
	BucketService bucket.Service
}

func New(dbService *db.DB, bucketService bucket.Service) Service {
	return &service{DB: dbService, BucketService: bucketService}
}

func (s *service) CreatePlayer(ctx context.Context, player *dto.PlayerDTO) error {
	p := models.NewPlayerModel(player)
	err := s.DB.Repository.Create(&p)

	(*player).ID = p.ID
	(*player).CreatedAt = p.CreatedAt
	(*player).UpdatedAt = p.UpdatedAt

	return err
}

func (s *service) GetPlayers(ctx context.Context, paging pagination.Pagination, players *[]dto.PlayerDTO) error {
	var p []models.Player

	err := s.DB.
		PlayerRepository.
		FindAllAndPaginate(paging, &p)

	*players = make([]dto.PlayerDTO, len(p))

	for key, value := range p {
		(*players)[key] = models.NewPlayerDTO(value)
	}

	return err
}

func (s *service) GetPlayer(ctx context.Context, id int, player *dto.PlayerDTO) error {
	var p models.Player

	err := s.DB.Repository.FindById(&p, id)

	*player = models.NewPlayerDTO(p)

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

func (s *service) UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader, player *dto.PlayerDTO) error {
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

	*player = models.NewPlayerDTO(p)

	return err
}
