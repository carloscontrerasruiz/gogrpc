package server

import (
	"context"

	"github.com/carloscontrerasruiz/gogrpc/models"
	"github.com/carloscontrerasruiz/gogrpc/repository"
	"github.com/carloscontrerasruiz/gogrpc/testspb"
)

type TestServer struct {
	repo repository.Repository
	testspb.UnimplementedTestServiceServer
}

func NewTestServer(repo repository.Repository) *TestServer {
	return &TestServer{repo: repo}
}

func (s *TestServer) GetTest(ctx context.Context, req *testspb.GetTestRequest) (*testspb.Test, error) {
	test, err := s.repo.GetTest(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &testspb.Test{
		Id:   test.Id,
		Name: test.Name,
	}, nil
}

func (s *TestServer) SetTest(ctx context.Context, req *testspb.Test) (*testspb.SetTestResponse, error) {
	test := &models.Test{
		Id:   req.Id,
		Name: req.Name,
	}

	err := s.repo.SetTest(ctx, test)
	if err != nil {
		return nil, err
	}

	return &testspb.SetTestResponse{
		Id: test.Id,
	}, nil
}
