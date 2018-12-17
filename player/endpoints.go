package player

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/logansua/nfl_app/pagination"
	"github.com/logansua/nfl_app/utils"
	"mime/multipart"
)

// Endpoints collects all of the endpoints that compose a profile service. It's
// meant to be used as a helper struct, to collect all of the endpoints into a
// single parameter.
//
// In a server, it's useful for functions that need to operate on a per-endpoint
// basis. For example, you might pass an Endpoints to a function that produces
// an http.Handler, with each method (endpoint) wired up to a specific path. (It
// is probably a mistake in design to invoke the Service methods on the
// Endpoints struct in a server.)
//
// In a client, it's useful to collect individually constructed endpoints into a
// single type that implements the Service interface. For example, you might
// construct individual endpoints using transport/http.NewClient, combine them
// into an Endpoints, and return it to the caller as a Service.
type Endpoints struct {
	CreatePlayerEndpoint           endpoint.Endpoint
	GetPlayersEndpoint             endpoint.Endpoint
	GetPlayerEndpoint              endpoint.Endpoint
	DeletePlayerEndpoint           endpoint.Endpoint
	MakeUploadPlayerAvatarEndpoint endpoint.Endpoint
}

// MakeServerEndpoints returns an Endpoints struct where each endpoint invokes
// the corresponding method on the provided service. Useful in a profilesvc
// server.
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		CreatePlayerEndpoint:           MakeCreatePlayerEndpoint(s),
		GetPlayersEndpoint:             MakeGetPlayersEndpoint(s),
		GetPlayerEndpoint:              MakeGetPlayerEndpoint(s),
		DeletePlayerEndpoint:           MakeDeletePlayerEndpoint(s),
		MakeUploadPlayerAvatarEndpoint: MakeUploadPlayerAvatarEndpoint(s),
	}
}

func (e Endpoints) CreatePlayer(ctx context.Context, p Player) error {
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

		p, err := service.CreatePlayer(ctx, req.Player)

		if err != nil {
			return nil, err
		}

		return dataResponse{Data: NewDTO(*p)}, nil
	}
}
func MakeGetPlayersEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getPlayersRequest)

		players, err := service.GetPlayers(ctx, req.Paging)

		if err != nil {
			return nil, err
		}

		dto := utils.Map(players, func(val interface{}) interface{} {
			return NewDTO(val.(Player))
		})

		return dataResponse{Data: dto}, nil
	}
}
func MakeGetPlayerEndpoint(service Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(playerIdRequest)

		player, err := service.GetPlayer(ctx, req.id)

		if err != nil {
			return nil, err
		}

		return dataResponse{Data: NewDTO(*player)}, nil
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

		player, err := service.UploadPlayerAvatar(ctx, req.id, req.file, &req.fileHeader)

		if err != nil {
			return nil, err
		}

		return dataResponse{Data: NewDTO(*player)}, nil
	}
}

type createPlayerRequest struct {
	Player Player
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
