package main

import (
	"context"
	"log"
	"time"

	pb "app/api/user/v1"

	"google.golang.org/grpc"
)

const (
	address = "localhost:18031"
	id      = 1
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Connect Failed: %v", err)
		return
	}
	defer conn.Close()
	client := pb.NewUserClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := client.Find(ctx, &pb.FindRequest{Id: id})
	if err != nil {
		log.Fatalf("Find User Failed: %v", err)
	}
	log.Printf("Find User's Nickname: %s", r.GetNickname())
}
