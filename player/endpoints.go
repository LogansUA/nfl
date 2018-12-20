package player

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
	CreatePlayerEndpoint           endpoint.Endpoint
	GetPlayersEndpoint             endpoint.Endpoint
	GetPlayerEndpoint              endpoint.Endpoint
	DeletePlayerEndpoint           endpoint.Endpoint
	MakeUploadPlayerAvatarEndpoint endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CreatePlayerEndpoint:           MakeCreatePlayerEndpoint(s),
		GetPlayersEndpoint:             MakeGetPlayersEndpoint(s),
		GetPlayerEndpoint:              MakeGetPlayerEndpoint(s),
		DeletePlayerEndpoint:           MakeDeletePlayerEndpoint(s),
		MakeUploadPlayerAvatarEndpoint: MakeUploadPlayerAvatarEndpoint(s),
	}
}

func (e Endpoints) CreatePlayer(ctx context.Context, p dto.PlayerDTO) error {
	request := createPlayerRequest{Player: p}
	response, err := e.CreatePlayerEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}
func (e Endpoints) GetPlayers(ctx context.Context, paging pagination.Pagination) error {
	request := getPlayersRequest{Paging: paging}
	response, err := e.GetPlayersEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}
func (e Endpoints) GetPlayer(ctx context.Context, id int) error {
	request := playerIdRequest{id: id}
	response, err := e.GetPlayerEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}
func (e Endpoints) DeletePlayer(ctx context.Context, id int) error {
	request := playerIdRequest{id: id}
	response, err := e.DeletePlayerEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}
func (e Endpoints) UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader multipart.FileHeader) error {
	request := bucket.UploadFileToBucketRequest{ID: id, File: file, FileHeader: fileHeader}
	response, err := e.MakeUploadPlayerAvatarEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(utils.DataResponse)

	return resp.Err
}

func MakeCreatePlayerEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(createPlayerRequest)

		p := req.Player

		err = service.CreatePlayer(ctx, &p)

		if err != nil {
			return nil, err
		}

		return utils.DataResponse{Data: p}, nil
	}
}
func MakeGetPlayersEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPlayersRequest)

		var players []dto.PlayerDTO

		err = service.GetPlayers(ctx, req.Paging, &players)

		if err != nil {
			return nil, err
		}

		return utils.DataResponse{Data: players}, nil
	}
}
func MakeGetPlayerEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(playerIdRequest)

		var player dto.PlayerDTO

		err = service.GetPlayer(ctx, req.id, &player)

		if err != nil {
			return nil, err
		}

		return utils.DataResponse{Data: player}, nil
	}
}
func MakeDeletePlayerEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(playerIdRequest)

		err = service.DeletePlayer(ctx, req.id)

		if err != nil {
			return nil, err
		}

		return nil, nil
	}
}
func MakeUploadPlayerAvatarEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(bucket.UploadFileToBucketRequest)

		var player dto.PlayerDTO
		err = service.UploadPlayerAvatar(ctx, req.ID, req.File, &req.FileHeader, &player)

		if err != nil {
			return nil, err
		}

		return utils.DataResponse{Data: player}, nil
	}
}

type createPlayerRequest struct {
	Player dto.PlayerDTO
}

type getPlayersRequest struct {
	Paging pagination.Pagination
}

type playerIdRequest struct {
	id int
}
