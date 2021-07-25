package api

import (
	"context"
	"fmt"
	pb "github.com/lucasszmt/grpcTraining/blog/gen/blog"
	"github.com/lucasszmt/grpcTraining/blog/server/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

var db *mongo.Database

func init() {
	db = database.MongoClient.Database("blog")
}

type blogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	AuthorId string             `json:"author_id,omitempty"`
	Content  string             `json:"content,omitempty"`
	Title    string             `json:"title,omitempty"`
}

type Server struct {
	pb.UnimplementedBlogServiceServer
}

func (*Server) CreateBlog(ctx context.Context, in *pb.CreateBlogRequest) (*pb.CreateBlogResponse, error) {
	log.Println("Creating new Blog")
	blog := in.GetBlog()

	data := blogItem{
		AuthorId: blog.AuthorId,
		Content:  blog.Content,
		Title:    blog.Title,
	}

	c := db.Collection("blog")
	res, err := c.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Cannot convert to OID"))
	}
	blogResp := &pb.Blog{
		Id:       id.Hex(),
		AuthorId: blog.AuthorId,
		Title:    blog.Title,
		Content:  blog.Content,
	}
	return &pb.CreateBlogResponse{Blog: blogResp}, nil
}
