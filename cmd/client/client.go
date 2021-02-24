package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/patricksalmeida/fc2-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	connection, err := grpc.Dial("localhost:50051", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("Could not connect to gRpc server: %v", err)
	}

	defer connection.Close()

	client := pb.NewUserServiceClient(connection)

	// AddUser(client)
	// AddUserVerbose(client)
	// AddUsers(client)
	AddUserStreamBoth(client)
}

func AddUser(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Airton",
		Email: "airton@teste.com",
	}

	res, err := client.AddUser(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRpc request: %v", err)
	}

	fmt.Println(res)
}

func AddUserVerbose(client pb.UserServiceClient) {

	req := &pb.User{
		Id:    "0",
		Name:  "Airton",
		Email: "airton@teste.com",
	}

	responseStream, err := client.AddUserVerbose(context.Background(), req)

	if err != nil {
		log.Fatalf("Could not make gRpc request: %v", err)
	}

	for {
		stream, err := responseStream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Could not receive gRpc message: %v", err)
		}

		fmt.Println("Status: ", stream.Status, " - ", stream.GetUser())
	}
}

func AddUsers(client pb.UserServiceClient) {
	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Will 1",
			Email: "will1@test.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Will 2",
			Email: "will2@test.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Will 3",
			Email: "will3@test.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Will 4",
			Email: "will4@test.com",
		},
		&pb.User{
			Id:    "w5",
			Name:  "Will 5",
			Email: "will5@test.com",
		},
	}

	stream, err := client.AddUsers(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	for _, req := range reqs {
		stream.Send(req)
		time.Sleep(time.Second * 3)
	}

	res, err := stream.CloseAndRecv()

	if err != nil {
		log.Fatalf("Error receiving response: %v", err)
	}

	fmt.Println(res)
}

func AddUserStreamBoth(client pb.UserServiceClient) {
	stream, err := client.AddUserStreamBoth(context.Background())

	if err != nil {
		log.Fatalf("Error creating request: %v", err)
	}

	reqs := []*pb.User{
		&pb.User{
			Id:    "w1",
			Name:  "Will 1",
			Email: "will1@test.com",
		},
		&pb.User{
			Id:    "w2",
			Name:  "Will 2",
			Email: "will2@test.com",
		},
		&pb.User{
			Id:    "w3",
			Name:  "Will 3",
			Email: "will3@test.com",
		},
		&pb.User{
			Id:    "w4",
			Name:  "Will 4",
			Email: "will4@test.com",
		},
		&pb.User{
			Id:    "w5",
			Name:  "Will 5",
			Email: "will5@test.com",
		},
	}

	wait := make(chan int)

	go func() {
		for _, req := range reqs {
			fmt.Println("Sending user: ", req.Name)
			stream.Send(req)
			time.Sleep(time.Second * 2)
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()

			if err == io.EOF {
				break
			}

			if err != nil {
				log.Fatalf("Error receiving data: %v", err)
			}

			fmt.Printf("Receiving user %v with status %v \n", res.GetUser().Name, res.GetStatus())
		}

		close(wait)
	}()

	<-wait
}
