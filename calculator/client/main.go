package main

import (
	"context"
	pb "github.com/lucasszmt/grpcTraining/calculator/gen/calculator"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

const (
	addr        = "localhost:50051"
	defaultName = "world"
)

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := pb.NewCalculatorClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	stream, streamErr := client.PrimeNumberDecomposition(ctx, &pb.NRequest{Number: 5})
	if streamErr != nil {
		log.Fatal(streamErr)
	}
	for {
		num, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		log.Println(num.Number)
	}
}
