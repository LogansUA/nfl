package player

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"github.com/logansua/nfl_app/pagination"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting is returned when an expected path variable is missing.
	// It always indicates programmer error.
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programmer error)")
)

// MakeHTTPHandler mounts all of the service endpoints into an http.Handler.
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	router := mux.NewRouter()

	endpoints := MakeServerEndpoints(s)

	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.
		Methods(http.MethodPost).
		Path("/players").
		Handler(httptransport.NewServer(
			endpoints.CreatePlayerEndpoint,
			decodeCreatePlayerRequest,
			encodeResponse,
			options...,
		))
	router.
		Methods(http.MethodGet).
		Path("/players").
		Handler(httptransport.NewServer(
			endpoints.GetPlayersEndpoint,
			decodeGetPlayersRequest,
			encodeResponse,
			options...,
		))
	router.
		Methods(http.MethodGet).
		Path("/players/{id}").
		Handler(httptransport.NewServer(
			endpoints.GetPlayerEndpoint,
			decodeGetPlayerRequest,
			encodeResponse,
			options...,
		))
	router.
		Methods(http.MethodDelete).
		Path("/players/{id}").
		Handler(httptransport.NewServer(
			endpoints.DeletePlayerEndpoint,
			decodeDeletePlayerRequest,
			encodeDeletePlayerResponse,
			options...,
		))
	router.
		Methods(http.MethodPut).
		Path("/players/{id}/avatar").
		Handler(httptransport.NewServer(
			endpoints.MakeUploadPlayerAvatarEndpoint,
			decodeUploadPlayerAvatarRequest,
			encodeResponse,
			options...,
		))

	return router
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

const (
	maxUploadSize = 2 * 1024 * 1024 // 2 mb
)

func decodeUploadPlayerAvatarRequest(ctx context.Context, r *http.Request) (request interface{}, err error) {
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

	var req uploadPlayerAvatarRequest

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		return nil, err
	}

	req.id = id
	req.file = file
	req.fileHeader = *fileHeader

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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	return json.NewEncoder(w).Encode(response)
}

func encodeDeletePlayerResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		// Not a Go kit transport error, but a business-logic error.
		// Provide those as HTTP errors.
		encodeError(ctx, e.error(), w)

		return nil
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusNoContent)

	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err == nil {
		panic("encodeError with nil error")
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))

	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}