package main

import (
	"context"
	"log"
	"net"

	"github.com/linkedin-sandbox/ob-proto-test/go-service/user"

	"google.golang.org/grpc"
)

type server struct {
	user.UnimplementedUserServiceServer
}

// global array of UserData
var users = []*user.UserData{
	{
		Username:  "user1",
		FirstName: "John",
		LastName:  "Doe",
		Email:     "user1@domain.com",
	},
	{
		Username:  "user2",
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "user2@domain.com",
	},
}

func print_users() {
	for _, u := range users {
		println("Username: ", u.Username)
		println("First name: ", u.FirstName)
		println("Last name: ", u.LastName)
		println("Email: ", u.Email)
		println("-----")
	}
}

func (s *server) GetUserDetails(ctx context.Context, in *user.GetUserDetailsRequest) (*user.GetUserDetailsResponse, error) {
	log.Printf("Received request: %v", in)
	// find the user in the array
	for _, u := range users {
		if u.Username == in.Username {
			println("Found user: ", u.Username)
			return &user.GetUserDetailsResponse{UserData: u}, nil
		}
	}
	return nil, nil
}

func (s *server) UpdateUserDetails(ctx context.Context, in *user.UpdateUserDetailsRequest) (*user.UpdateUserDetailsResponse, error) {
	log.Printf("Received request: %v", in)
	// find the user in the array
	for _, u := range users {
		if u.Username == in.UserData.Username {
			u.Email = in.UserData.Email
			u.FirstName = in.UserData.FirstName
			u.LastName = in.UserData.LastName
			print_users()
			return &user.UpdateUserDetailsResponse{UserData: u}, nil
		}
	}
	return nil, nil
}

func main() {
	println("Starting server...")
	print_users()
	// Listen on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register the UserService implementation with the gRPC server
	user.RegisterUserServiceServer(grpcServer, &server{})

	// Start the gRPC server
	log.Println("Server listening on port 50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
