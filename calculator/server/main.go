package main

import (
	"context"
	pb "github.com/lucasszmt/grpcTraining/calculator/gen/calculator"
	"github.com/lucasszmt/grpcTraining/calculator/server/utils"
	"google.golang.org/grpc"
	"io"
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

func (s *server) ComputeAverage(stream pb.Calculator_ComputeAverageServer) error {
	log.Println("Receiving Stream")
	count, avg := 0, 0.
	for {
		val, err := stream.Recv()
		if err == io.EOF {
			avg = avg / float64(count)
			log.Println("Sending: ", avg)
			return stream.SendAndClose(&pb.NumberStream{Number: avg})
		}
		if err != nil {
			return err
		}
		log.Println(val.Number)
		count++
		avg += val.Number
	}
}

func (*server) FindMaximum(stream pb.Calculator_FindMaximumServer) error {
	var max int32
	for {
		req, recvEr := stream.Recv()
		if recvEr == io.EOF {
			log.Println("End of stream")
			return nil
		}
		if recvEr != nil {
			return recvEr
		}
		if req.Number > max {
			max = req.Number
			sendErr := stream.Send(&pb.FindMaxResponse{Number: max})
			if sendErr != nil {
				log.Fatal(sendErr)
			}
		}
	}
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
