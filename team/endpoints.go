package team

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/logansua/nfl_app/bucket"
	"github.com/logansua/nfl_app/models/dto"
	"github.com/logansua/nfl_app/pagination"
	"github.com/logansua/nfl_app/utils"
	"mime/multipart"
)

type Endpoints struct {
	CreateTeamEndpoint         endpoint.Endpoint
	GetTeamsEndpoint           endpoint.Endpoint
	GetTeamEndpoint            endpoint.Endpoint
	DeleteTeamEndpoint         endpoint.Endpoint
	MakeUploadTeamLogoEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CreateTeamEndpoint:         MakeCreateTeamEndpoint(s),
		GetTeamsEndpoint:           MakeGetTeamsEndpoint(s),
		GetTeamEndpoint:            MakeGetTeamEndpoint(s),
		DeleteTeamEndpoint:         MakeDeleteTeamEndpoint(s),
		MakeUploadTeamLogoEndpoint: MakeUploadTeamLogoEndpoint(s),
	}
}

func (e Endpoints) CreateTeam(ctx context.Context, p dto.TeamDTO) error {
	request := createTeamRequest{Team: p}
	response, err := e.CreateTeamEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}
func (e Endpoints) GetTeams(ctx context.Context, paging pagination.Pagination) error {
	request := getTeamsRequest{Paging: paging}
	response, err := e.GetTeamsEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}
func (e Endpoints) GetTeam(ctx context.Context, id int) error {
	request := teamIdRequest{id: id}
	response, err := e.GetTeamEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}
func (e Endpoints) DeleteTeam(ctx context.Context, id int) error {
	request := teamIdRequest{id: id}
	response, err := e.DeleteTeamEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}
func (e Endpoints) UploadTeamAvatar(ctx context.Context, id int, file multipart.File, fileHeader multipart.FileHeader) error {
	request := bucket.UploadFileToBucketRequest{ID: id, File: file, FileHeader: fileHeader}
	response, err := e.MakeUploadTeamLogoEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}

func MakeCreateTeamEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createTeamRequest)

		p := req.Team

		err = service.CreateTeam(ctx, &p)

		if err != nil {
			return nil, err
		}

		return utils.DataResponse{Data: p}, nil
	}
}
func MakeGetTeamsEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getTeamsRequest)

		var teams []dto.TeamDTO

		err = service.GetTeams(ctx, req.Paging, &teams)

		if err != nil {
			return nil, err
		}

		return utils.DataResponse{Data: teams}, nil
	}
}
func MakeGetTeamEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(teamIdRequest)

		var team dto.TeamDTO

		err = service.GetTeam(ctx, req.id, &team)

		if err != nil {
			return nil, err
		}

		return utils.DataResponse{Data: team}, nil
	}
}
func MakeDeleteTeamEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(teamIdRequest)

		err = service.DeleteTeam(ctx, req.id)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
func MakeUploadTeamLogoEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bucket.UploadFileToBucketRequest)

		var team dto.TeamDTO
		err = service.UploadTeamLogo(ctx, req.ID, req.File, &req.FileHeader, &team)

		if err != nil {
			return nil, err
		}

		return utils.DataResponse{Data: team}, nil
	}
}

type createTeamRequest struct {
	Team dto.TeamDTO
}

type getTeamsRequest struct {
	Paging pagination.Pagination
}

type teamIdRequest struct {
	id int
}
