package tc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/jho4us/tc/test"

	"github.com/gorilla/mux"

	"net/http"

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
)

// MakeHandler returns a handler for the test constructor service.
func MakeHandler(tcs Service, logger kitlog.Logger) http.Handler {
	opts := []kithttp.ServerOption{
		kithttp.ServerErrorLogger(logger),
		kithttp.ServerErrorEncoder(encodeError),
	}

	createTestHandler := kithttp.NewServer(
		makeCreateTestEndpoint(tcs),
		decodeCreateTestRequest,
		encodeResponse,
		opts...,
	)
	putTestHandler := kithttp.NewServer(
		makePutTestEndpoint(tcs),
		decodePutTestRequest,
		encodeResponse,
		opts...,
	)
	loadTestHandler := kithttp.NewServer(
		makeLoadTestEndpoint(tcs),
		decodeLoadTestRequest,
		encodeResponse,
		opts...,
	)
	deleteTestHandler := kithttp.NewServer(
		makeDeleteTestEndpoint(tcs),
		decodeDeleteTestRequest,
		encodeResponse,
		opts...,
	)
	listTestsHandler := kithttp.NewServer(
		makeListTestsEndpoint(tcs),
		decodeListTestsRequest,
		encodeResponse,
		opts...,
	)

	r := mux.NewRouter()

	r.Handle("/tc/v1/tests", createTestHandler).Methods(http.MethodPost)
	r.Handle("/tc/v1/tests", listTestsHandler).Methods(http.MethodGet)
	r.Handle("/tc/v1/tests/{id}", loadTestHandler).Methods(http.MethodGet)
	r.Handle("/tc/v1/tests/{id}", putTestHandler).Methods(http.MethodPost)
	r.Handle("/tc/v1/tests/{id}", deleteTestHandler).Methods(http.MethodDelete)

	return r
}

var errBadRoute = errors.New("bad route")

func decodeCreateTestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var body struct {
		Name string `json:"name"`
	}
	if r.Body == nil {
		return nil, errBadRoute
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return createTestRequest{
		Name: body.Name,
	}, nil
}

func decodeLoadTestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return loadTestRequest{ID: test.ID(id)}, nil
}

func decodePutTestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok || r.Body == nil {
		return nil, errBadRoute
	}
	var body struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}

	return putTestRequest{&Test{ID: id, Name: body.Name, Content: body.Content}}, nil
}

func decodeDeleteTestRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	return deleteTestRequest{ID: test.ID(id)}, nil
}

func decodeListTestsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return listTestsRequest{}, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

type errorer interface {
	error() error
}

// encode errors from business-logic
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	case test.ErrUnknown:
		w.WriteHeader(http.StatusNotFound)
	case test.ErrInvalidArgument:
		w.WriteHeader(http.StatusBadRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
