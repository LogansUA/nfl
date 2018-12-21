package player

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/logansua/nfl_app/bucket"
	apperrors "github.com/logansua/nfl_app/errors"
	"github.com/logansua/nfl_app/pagination"
	"github.com/logansua/nfl_app/router"
	"net/http"
	"strconv"
)

func GetServiceOptions(logger log.Logger) []httptransport.ServerOption {
	return []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
		httptransport.ServerAfter(func(ctx context.Context, writer http.ResponseWriter) context.Context {
			writer.Header().Set("Content-type", "application/json")

			return ctx
		}),
		httptransport.ServerFinalizer(func(ctx context.Context, code int, r *http.Request) {
			logger.Log("METHOD", r.Method, "PATH", r.URL.Path, "CODE", code)
		}),
	}
}

func CreateRoutes(s Service, logger log.Logger) []router.Route {
	endpoints := MakeServerEndpoints(s)

	options := GetServiceOptions(logger)

	return []router.Route{
		{
			Name:        "Create good",
			Method:      http.MethodPost,
			Path:        "/players",
			StrictSlash: false,
			Handler: httptransport.NewServer(
				endpoints.CreatePlayerEndpoint,
				decodeCreatePlayerRequest,
				encodeResponse,
				options...,
			),
		},
		{
			Name:        "Get goods",
			Method:      http.MethodGet,
			Path:        "/players",
			StrictSlash: true,
			Handler: httptransport.NewServer(
				endpoints.GetPlayersEndpoint,
				decodeGetPlayersRequest,
				encodeResponse,
				options...,
			),
		},
		{
			Name:        "Update good",
			Method:      http.MethodGet,
			Path:        "/players/{id}",
			StrictSlash: true,
			Handler: httptransport.NewServer(
				endpoints.GetPlayerEndpoint,
				decodeGetPlayerRequest,
				encodeResponse,
				options...,
			),
		},
		{
			Name:        "Delete good",
			Method:      http.MethodDelete,
			Path:        "/players/{id}",
			StrictSlash: true,
			Handler: httptransport.NewServer(
				endpoints.DeletePlayerEndpoint,
				decodeDeletePlayerRequest,
				encodeDeletePlayerResponse,
				options...,
			),
		},
		{
			Name:        "Upload good photo",
			Method:      http.MethodPut,
			Path:        "/players/{id}/avatar",
			StrictSlash: false,
			Handler: httptransport.NewServer(
				endpoints.MakeUploadPlayerAvatarEndpoint,
				decodeUploadPlayerAvatarRequest,
				encodeResponse,
				options...,
			),
		},
	}
}

func decodeCreatePlayerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req createPlayerRequest

	if e := json.NewDecoder(r.Body).Decode(&req.Player); e != nil {
		return nil, e
	}

	return req, nil
}
func decodeGetPlayersRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req getPlayersRequest

	params := r.URL.Query()

	paging := pagination.New(params)

	req.Paging = paging

	return req, nil
}
func decodeGetPlayerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req playerIdRequest

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		return nil, err
	}

	req.id = id

	return req, nil
}
func decodeDeletePlayerRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req playerIdRequest

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		return nil, err
	}

	req.id = id

	return req, nil
}
func decodeUploadPlayerAvatarRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	const (
		maxUploadSize = 2 * 1024 * 1024 // 2 mb
	)

	//r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return nil, errors.New("file is too big")
	}

	file, fileHeader, err := r.FormFile("image")
	if err == http.ErrMissingFile {
		return
	}
	if err != nil {
		return nil, err
	}

	var req bucket.UploadFileToBucketRequest

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		return nil, err
	}

	req.ID = id
	req.File = file
	req.FileHeader = *fileHeader

	return req, nil
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)

		return nil
	}

	return json.NewEncoder(w).Encode(response)
}

func encodeDeletePlayerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)

		return nil
	}

	w.WriteHeader(http.StatusNoContent)

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(codeFrom(err))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case apperrors.ErrNotFound:
		return http.StatusNotFound
	case apperrors.ErrAlreadyExists, apperrors.ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
