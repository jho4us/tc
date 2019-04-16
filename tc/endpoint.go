package tc

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/jho4us/tc/test"
)

type createTestRequest struct {
	Name string
}

type createTestResponse struct {
	ID  test.ID `json:"id,omitempty"`
	Err error   `json:"error,omitempty"`
}

func (r createTestResponse) error() error { return r.Err }

func makeCreateTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(createTestRequest)
		id, err := s.CreateTest(req.Name)
		return createTestResponse{ID: id, Err: err}, nil
	}
}

type loadTestRequest struct {
	ID test.ID
}

type loadTestResponse struct {
	Test *Test `json:"test,omitempty"`
	Err  error `json:"error,omitempty"`
}

func (r loadTestResponse) error() error { return r.Err }

func makeLoadTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(loadTestRequest)
		t, err := s.LoadTest(req.ID)
		return loadTestResponse{Test: &t, Err: err}, nil
	}
}

type putTestRequest struct {
	Test *Test
}

type putTestResponse struct {
	Err error `json:"error,omitempty"`
}

func (r putTestResponse) error() error { return r.Err }

func makePutTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(putTestRequest)
		err := s.PutTest(req.Test)
		return putTestResponse{Err: err}, nil
	}
}

type deleteTestRequest struct {
	ID test.ID
}

type deleteTestResponse struct {
	Err error `json:"error,omitempty"`
}

func (r deleteTestResponse) error() error { return r.Err }

func makeDeleteTestEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(deleteTestRequest)
		err := s.DeleteTest(req.ID)
		return deleteTestResponse{Err: err}, nil
	}
}

type listTestsRequest struct{}

type listTestsResponse struct {
	Tests []Test `json:"tests,omitempty"`
	Err   error  `json:"error,omitempty"`
}

func (r listTestsResponse) error() error { return r.Err }

func makeListTestsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(listTestsRequest)
		return listTestsResponse{Tests: s.Tests(), Err: nil}, nil
	}
}
