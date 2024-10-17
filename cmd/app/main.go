package main

import (
	"context"
	pup_service "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	c "gitlab.ozon.dev/akugnerevich/homework.git/internal/cli"
	"log"
)

const grpcHost = "localhost:7002"

func main() {
	ctx, _ := context.WithCancel(context.Background())

	conn, err := grpc.NewClient(grpcHost, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to create grpc client: %v", err)
	}
	defer conn.Close()

	pupServiceClient := pup_service.NewPupServiceClient(conn)

	err = c.Run(ctx, pupServiceClient)
	if err != nil {
		log.Printf("Error in Run: %v", err)
	}

}
