package player

import (
	"context"
	"github.com/logansua/nfl_app/db"
	"github.com/logansua/nfl_app/mocks"
	"github.com/logansua/nfl_app/models"
	"github.com/logansua/nfl_app/models/dto"
	"github.com/logansua/nfl_app/pagination"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/url"
	"testing"
	"time"
)

func TestService_GetPlayer(t *testing.T) {
	team := models.Team{
		ID:        1,
		Name:      "TEST_TEAM",
		Logo:      "TEST_LOGO",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	player := models.Player{
		ID:        1,
		Name:      "TEST_PLAYER",
		Avatar:    "TEST_AVATAR",
		Team:      team,
		TeamID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repository := &mocks.Repository{}
	repository.On("FindById", mock.AnythingOfType("*models.Player"), mock.AnythingOfType("int")).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.Player)

			*arg = player
		}).
		Return(nil)

	playerService := New(&db.DB{Repository: repository}, nil, nil)
	var actualPlayer dto.PlayerDTO

	err := playerService.GetPlayer(context.Background(), int(player.ID), &actualPlayer)

	assert.Nil(t, err)
	assert.NotEmpty(t, actualPlayer)
	assert.Equal(t, models.NewPlayerDTO(player), actualPlayer)

	repository.AssertExpectations(t)
}

func TestService_GetPlayers(t *testing.T) {
	team := models.Team{
		ID:        1,
		Name:      "TEST_TEAM",
		Logo:      "TEST_LOGO",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	players := []models.Player{
		{
			ID:        1,
			Name:      "TEST_PLAYER",
			Avatar:    "TEST_AVATAR",
			Team:      team,
			TeamID:    1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	playerDTOS := make([]dto.PlayerDTO, len(players))

	for key, value := range players {
		playerDTOS[key] = models.NewPlayerDTO(value)
	}

	playerRepository := &mocks.PlayerRepository{}
	playerRepository.On(
		"FindAllAndPaginate",
		mock.AnythingOfType("pagination.Pagination"),
		mock.AnythingOfType("*[]models.Player"),
	).
		Run(func(args mock.Arguments) {
			arg := args.Get(1).(*[]models.Player)

			*arg = players
		}).
		Return(nil)

	playerService := New(&db.DB{PlayerRepository: playerRepository}, nil, nil)

	values := url.Values{
		"page":     []string{"1"},
		"per_page": []string{"10"},
	}
	paging := pagination.New(values)
	var actualPlayers []dto.PlayerDTO

	err := playerService.GetPlayers(context.Background(), paging, &actualPlayers)

	assert.Nil(t, err)
	assert.NotEmpty(t, actualPlayers)
	assert.Equal(t, playerDTOS, actualPlayers)

	playerRepository.AssertExpectations(t)
}

func TestService_DeletePlayer(t *testing.T) {
	team := models.Team{
		ID:        1,
		Name:      "TEST_TEAM",
		Logo:      "TEST_LOGO",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	player := models.Player{
		ID:        1,
		Name:      "TEST_PLAYER",
		Avatar:    "TEST_AVATAR",
		Team:      team,
		TeamID:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	repository := &mocks.Repository{}
	repository.On("FindById", mock.AnythingOfType("*models.Player"), mock.AnythingOfType("int")).
		Run(func(args mock.Arguments) {
			arg := args.Get(0).(*models.Player)

			*arg = player
		}).
		Return(nil)
	repository.On("Delete", mock.AnythingOfType("*models.Player"), mock.AnythingOfType("int")).
		Return(nil)

	playerService := New(&db.DB{Repository: repository}, nil, nil)

	err := playerService.DeletePlayer(context.Background(), int(player.ID))

	assert.Nil(t, err)

	repository.AssertExpectations(t)
}
