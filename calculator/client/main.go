package main

import (
	"context"
	pb "github.com/lucasszmt/grpcTraining/calculator/gen/calculator"
	"google.golang.org/grpc"
	"log"
)

const (
	addr        = "localhost:50051"
	defaultName = "world"
)

func ComputeAverage(client pb.CalculatorClient, nums []float64) error {
	stream, _ := client.ComputeAverage(context.Background())
	for _, num := range nums {
		err := stream.Send(&pb.NumberStream{Number: num})
		if err != nil {
			return err
		}
	}
	reply, err := stream.CloseAndRecv()
	log.Println(reply)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := pb.NewCalculatorClient(conn)
	cerr := ComputeAverage(client, []float64{4, 5, 7, 8, 10})
	log.Println(cerr)
}