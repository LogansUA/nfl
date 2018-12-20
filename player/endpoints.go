package player

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/logansua/nfl_app/models"
	"github.com/logansua/nfl_app/pagination"
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

func (e Endpoints) CreatePlayer(ctx context.Context, p models.PlayerDTO) error {
	request := createPlayerRequest{Player: p}
	response, err := e.CreatePlayerEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(dataResponse)

	return resp.Err
}
func (e Endpoints) GetPlayers(ctx context.Context, paging pagination.Pagination) error {
	request := getPlayersRequest{Paging: paging}
	response, err := e.GetPlayersEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(dataResponse)

	return resp.Err
}
func (e Endpoints) GetPlayer(ctx context.Context, id int) error {
	request := playerIdRequest{id: id}
	response, err := e.GetPlayerEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(dataResponse)

	return resp.Err
}
func (e Endpoints) DeletePlayer(ctx context.Context, id int) error {
	request := playerIdRequest{id: id}
	response, err := e.DeletePlayerEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(dataResponse)

	return resp.Err
}
func (e Endpoints) UploadPlayerAvatar(ctx context.Context, id int, file multipart.File, fileHeader multipart.FileHeader) error {
	request := uploadPlayerAvatarRequest{id: id, file: file, fileHeader: fileHeader}
	response, err := e.MakeUploadPlayerAvatarEndpoint(ctx, request)

	if err != nil {
		return err
	}

	resp := response.(dataResponse)

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

		return dataResponse{Data: p}, nil
	}
}
func MakeGetPlayersEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPlayersRequest)

		var players []models.PlayerDTO

		err = service.GetPlayers(ctx, req.Paging, &players)

		if err != nil {
			return nil, err
		}

		return dataResponse{Data: players}, nil
	}
}
func MakeGetPlayerEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(playerIdRequest)

		var player models.PlayerDTO

		err = service.GetPlayer(ctx, req.id, &player)

		if err != nil {
			return nil, err
		}

		return dataResponse{Data: player}, nil
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
		req := request.(uploadPlayerAvatarRequest)

		var player models.PlayerDTO
		err = service.UploadPlayerAvatar(ctx, req.id, req.file, &req.fileHeader, &player)

		if err != nil {
			return nil, err
		}

		return dataResponse{Data: player}, nil
	}
}

type createPlayerRequest struct {
	Player models.PlayerDTO
}

type getPlayersRequest struct {
	Paging pagination.Pagination
}

type playerIdRequest struct {
	id int
}

type uploadPlayerAvatarRequest struct {
	id         int
	file       multipart.File
	fileHeader multipart.FileHeader
}

type dataResponse struct {
	Data interface{} `json:"data"`
	Err  error       `json:"error,omitempty"`
}
