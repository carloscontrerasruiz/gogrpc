package server

import (
	"context"
	"io"

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

func (s *TestServer) SetQuestions(stream testspb.TestService_SetQuestionsServer) error {
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&testspb.SetQuestionResponse{
				Ok: true,
			})
		}
		if err != nil {
			return err
		}

		question := &models.Question{
			Id:       msg.GetId(),
			Answer:   msg.Answer,
			Question: msg.Question,
			TestId:   msg.GetTestId(),
		}

		err = s.repo.SetQuestion(context.Background(), question)
		if err != nil {
			return stream.SendAndClose(&testspb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}
