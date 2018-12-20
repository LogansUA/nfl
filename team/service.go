package team

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
	CreateTeam(ctx context.Context, team *dto.TeamDTO) error
	// Get list of players
	GetTeams(ctx context.Context, paging pagination.Pagination, players *[]dto.TeamDTO) error
	// Get single player by ID
	GetTeam(ctx context.Context, id int, player *dto.TeamDTO) error
	// Delete player by ID
	DeleteTeam(ctx context.Context, id int) error
	// Upload player avatar by ID
	UploadTeamLogo(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader, p *dto.TeamDTO) error
}

type service struct {
	DB            *db.DB
	BucketService bucket.Service
}

func New(dbService *db.DB, bucketService bucket.Service) Service {
	return &service{DB: dbService, BucketService: bucketService}
}

func (s *service) CreateTeam(ctx context.Context, team *dto.TeamDTO) error {
	t := models.NewTeamModel(team)

	err := s.DB.Repository.Create(&t)

	*team = models.NewTeamDTO(t)

	return err
}

func (s *service) GetTeams(ctx context.Context, paging pagination.Pagination, teams *[]dto.TeamDTO) error {
	var t []models.Team

	err := s.DB.
		TeamRepository.
		FindAllAndPaginate(paging, &t)

	*teams = make([]dto.TeamDTO, len(t))

	for key, value := range t {
		(*teams)[key] = models.NewTeamDTO(value)
	}

	return err
}

func (s *service) GetTeam(ctx context.Context, id int, team *dto.TeamDTO) error {
	var t models.Team

	err := s.DB.Repository.FindById(&t, id)

	*team = models.NewTeamDTO(t)

	return err
}

func (s *service) DeleteTeam(ctx context.Context, id int) error {
	var t models.Team

	err := s.DB.Repository.FindById(&t, id)

	if err != nil {
		return err
	}

	err = s.DB.Repository.Delete(&t, id)

	return err
}

func (s *service) UploadTeamLogo(ctx context.Context, id int, file multipart.File, fileHeader *multipart.FileHeader, team *dto.TeamDTO) error {
	var t models.Team

	if err := s.DB.Repository.FindById(&t, id); err != nil {
		return errors.New("player not found")
	}

	name, err := s.BucketService.UploadTeamLogo(ctx, t.ID, file, fileHeader)

	if err != nil {
		return err
	}

	t.Logo = name

	err = s.DB.Repository.Save(&t)

	*team = models.NewTeamDTO(t)

	return err
}
