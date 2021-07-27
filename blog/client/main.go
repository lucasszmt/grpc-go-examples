package main

import (
	"context"
	pb "github.com/lucasszmt/grpcTraining/blog/gen/blog"
	"google.golang.org/grpc"
	"log"
)

func main() {
	conn, err := grpc.Dial("0.0.0.0:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	client := pb.NewBlogServiceClient(conn)
	res, err := client.ReadBlog(context.Background(), &pb.ReadBlogRequest{BlogId: "123455"})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
