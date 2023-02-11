package main

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/carloscontrerasruiz/gogrpc/testspb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(
		"localhost:5070",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)

	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	defer conn.Close()

	c := testspb.NewTestServiceClient(conn)
	//DoUnary(c)
	//DoClientStreaming(c)
	//DoServerStreaming(c)
	DoBidireccionalStreaming(c)
}

func DoUnary(c testspb.TestServiceClient) {
	req := &testspb.GetTestRequest{
		Id: "t1",
	}

	res, err := c.GetTest(context.Background(), req)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	log.Printf("Response: %v", res)
}

func DoClientStreaming(c testspb.TestServiceClient) {
	questions := []*testspb.Question{
		{
			Id:       "q5",
			Answer:   "answer",
			Question: "the question",
			TestId:   "t1",
		},
		{
			Id:       "q6",
			Answer:   "answer",
			Question: "the question",
			TestId:   "t1",
		},
	}

	stream, err := c.SetQuestions(context.Background())

	if err != nil {
		log.Fatalf("Error %v", err)
	}

	for _, question := range questions {
		log.Println("sending question: ", question.Id)
		stream.Send(question)
		time.Sleep(2 * time.Second)
	}

	msg, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	log.Println(msg)

}

func DoServerStreaming(c testspb.TestServiceClient) {
	req := &testspb.GetStudentsPerTestRequest{
		TestId: "t1",
	}

	stream, err := c.GetStudentsPerTest(context.Background(), req)
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error %v", err)
		}

		log.Println(msg)
	}
}

func DoBidireccionalStreaming(c testspb.TestServiceClient) {
	answer := testspb.TakeTestRequest{
		Answer: "The answer",
	}

	numberOfquestions := 4

	waitChannel := make(chan chan struct{})

	stream, err := c.TakeTest(context.Background())
	if err != nil {
		log.Fatalf("Error %v", err)
	}

	go func() {
		for i := 0; i < numberOfquestions; i++ {
			stream.Send(&answer)
			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error %v", err)
				break
			}

			log.Println(res)
		}
		log.Println("Close channel")
		close(waitChannel)
	}()

	<-waitChannel
}
