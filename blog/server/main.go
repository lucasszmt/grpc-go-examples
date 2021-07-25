package main

import (
	pb "github.com/lucasszmt/grpcTraining/blog/gen/blog"
	"github.com/lucasszmt/grpcTraining/blog/server/api"
	"github.com/lucasszmt/grpcTraining/blog/server/database"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	defer func() {
		err := database.Disconnect()
		if err != nil {
			log.Fatal(err)
		}
	}()
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal(err)
	}
	s := grpc.NewServer()
	pb.RegisterBlogServiceServer(s, &api.Server{})
	log.Println("listening: ", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
