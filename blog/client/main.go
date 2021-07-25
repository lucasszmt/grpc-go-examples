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
	resp, _ := client.CreateBlog(context.Background(), &pb.CreateBlogRequest{Blog: &pb.Blog{
		AuthorId: "1",
		Title:    "The Witcher",
		Content:  "The white wolf returns",
	}})
	log.Println(resp)
}
