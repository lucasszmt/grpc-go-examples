package main

import (
	"context"
	pb "github.com/lucasszmt/grpcTraining/calculator/gen/calculator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"log"
	"time"
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

func FindMaximum(client pb.CalculatorClient, nums []int32) {
	done := make(chan struct{}, 1)
	stream, err := client.FindMaximum(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	go func() {
		for _, num := range nums {
			err := stream.Send(&pb.FindMaxRequest{Number: num})
			if err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
		close(done)
	}()
	go func() {
		for {
			resp, err := stream.Recv()
			if err != nil {
				log.Fatal(err)
			}
			log.Println(resp.Number)
		}
	}()
	<-done
}

func main() {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	client := pb.NewCalculatorClient(conn)
	_, respErr := client.Sum(context.Background(), &pb.SumRequest{Values: &pb.Values{ValA: -1, ValB: 2}})
	if respErr != nil {
		s, _ := status.FromError(respErr)
		log.Print("Message",s.Message())
		log.Print("Code",s.Code())
		log.Print("Details",s.Err())
	}
}
