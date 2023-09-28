package main

import (
	"context"
	"log"

	"github.com/linkedin-sandbox/ob-proto-test/go-client/user"
	"google.golang.org/grpc"
)

func main() {
	// Connect to Service B (Server) at localhost:50051
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create a UserService client
	client := user.NewUserServiceClient(conn)

	// Create a UserData message with username field
	request := &user.GetUserDetailsRequest{
		Username: "user2",
	}

	// Make a GetUserDetails RPC call to Service B
	response, err := client.GetUserDetails(context.Background(), request)
	if err != nil {
		log.Fatalf("Could not get user details: %v", err)
	}

	// Print the received response
	log.Printf("Response received: %v", response)

	new_user := response.UserData

	new_user.LastName = "updated"

	// Make a UpdateUserDetails RPC call to Service B
	update_request := &user.UpdateUserDetailsRequest{
		UserData: new_user,
	}

	update_response, err := client.UpdateUserDetails(context.Background(), update_request)
	if err != nil {
		log.Fatalf("Could not update user details: %v", err)
	}
	log.Printf("Response received: %v", update_response)
}
