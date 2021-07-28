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
	//res, err := client.UpdateBlog(context.Background(), &pb.UpdateBlogRequest{Blog: &pb.Blog{
	//	Id:       "6100b91fa3730f4b4f0cd2d0",
	//	AuthorId: "1",
	//	Title:    "The Sword of Destiny",
	//	Content:  "Book written by Andrey Popovich",
	//}})
	res, err := client.DeleteBlog(context.Background(), &pb.DeleteBlogRequest{Id: "6100b91fa3730f4b4f0cd2d0"})
	if err != nil {
		log.Println(err)
	}
	log.Println(res)
}
