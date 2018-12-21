package models

import "github.com/logansua/nfl_app/models/dto"

func NewTeamDTO(data Team) dto.TeamDTO {
	return dto.TeamDTO{
		ID:        data.ID,
		Name:      data.Name,
		Logo:      data.Logo,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func NewPlayerDTO(data Player) dto.PlayerDTO {
	return dto.PlayerDTO{
		ID:        data.ID,
		Name:      data.Name,
		Avatar:    data.Avatar,
		TeamID:    data.TeamID,
		CreatedAt: data.CreatedAt,
		UpdatedAt: data.UpdatedAt,
	}
}

func NewPlayerModel(data *dto.PlayerDTO) Player {
	return Player{
		Name:   data.Name,
		Avatar: data.Avatar,
		TeamID: data.TeamID,
	}
}

func NewTeamModel(data *dto.TeamDTO) Team {
	return Team{
		Name: data.Name,
		Logo: data.Logo,
	}
}
