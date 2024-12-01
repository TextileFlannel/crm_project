package grpc_test

import (
	"context"
	"log"
	"testing"
	"time"

	pb "http-server/danilkovalev/internal/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUnsubscribeAccount(t *testing.T) {
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAccountServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.UnsubscribeAccount(ctx, &pb.UnsubscribeRequest{AccountId: 1})
	if err != nil {
		t.Fatalf("UnsubscribeAccount failed: %v", err)
	}

	if !resp.Success {
		t.Errorf("Expected successful unsubscription, but got failure")
	}

	log.Printf("Unsubscription successful: %v", resp.Success)
}