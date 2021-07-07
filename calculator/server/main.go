package main

import (
	"context"
	pb "github.com/lucasszmt/grpcTraining/calculator/gen/calculator"
	"github.com/lucasszmt/grpcTraining/calculator/server/utils"
	"google.golang.org/grpc"
	"log"
	"net"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

func (s *server) Sum(ctx context.Context, in *pb.SumRequest) (*pb.ResultResponse, error) {
	result := in.Values.ValA + in.Values.ValB
	response := &pb.ResultResponse{Response: result}
	return response, nil
}

func (s *server) PrimeNumberDecomposition(req *pb.NRequest, stream pb.Calculator_PrimeNumberDecompositionServer) error {
	log.Println("Received : ", req.Number)

	primNum, _ := utils.PrimeNumberDecompose(req.Number)
	for _, i := range primNum {
		err := stream.Send(&pb.DecomposedNumber{Number: i})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterCalculatorServer(s, &server{})
	log.Println("listening: ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
