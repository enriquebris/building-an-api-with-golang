package main

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/maddevsio/grpc-rest-api-example/pb"
	"google.golang.org/grpc"
)

const address = "localhost:8080"

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewGreeterClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "EBR"}); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(r.Message)
	}
}
