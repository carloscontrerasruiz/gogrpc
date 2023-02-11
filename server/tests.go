package server

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/carloscontrerasruiz/gogrpc/models"
	"github.com/carloscontrerasruiz/gogrpc/repository"
	"github.com/carloscontrerasruiz/gogrpc/studentpb"
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

func (s *TestServer) EnrollStudents(stream testspb.TestService_EnrollStudentsServer) error {
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

		enrollment := &models.Enrollment{
			StudentId: msg.StudentId,
			TestId:    msg.TestId,
		}

		err = s.repo.SetEnrollment(context.Background(), enrollment)
		if err != nil {
			return stream.SendAndClose(&testspb.SetQuestionResponse{
				Ok: false,
			})
		}
	}
}

func (s *TestServer) GetStudentsPerTest(req *testspb.GetStudentsPerTestRequest, stream testspb.TestService_GetStudentsPerTestServer) error {
	students, err := s.repo.GetStudentsPerTest(context.Background(), req.GetTestId())
	if err != nil {
		return err
	}

	for _, student := range students {
		studentProto := &studentpb.Student{
			Id:   student.Id,
			Name: student.Name,
			Age:  student.Age,
		}

		err := stream.Send(studentProto)
		time.Sleep(2 * time.Second)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *TestServer) TakeTest(stream testspb.TestService_TakeTestServer) error {
	questions, err := s.repo.GetQuestionPerTest(context.Background(), "t1")
	if err != nil {
		return err
	}

	i := 0
	var currentQuestion = &models.Question{}

	for {
		if i < len(questions) {
			currentQuestion = questions[i]
		}

		if i <= len(questions) {
			questionToSend := &testspb.Question{
				Id:       currentQuestion.Id,
				Question: currentQuestion.Question,
			}
			err := stream.Send(questionToSend)
			if err != nil {
				return err
			}
			i++
		}

		answer, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		log.Println(answer)
	}
}
