package api

import (
	"context"
	"fmt"
	pb "github.com/lucasszmt/grpcTraining/blog/gen/blog"
	"github.com/lucasszmt/grpcTraining/blog/server/database"
	"go.mongodb.org/mongo-driver/bson"
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

func (*Server) ReadBlog(ctx context.Context, in *pb.ReadBlogRequest) (*pb.ReadBlogResponse, error) {
	log.Println("Reading blog...")
	blogId := in.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error parsing the object ID: %v", err))
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}
	res := db.Collection("blog").FindOne(ctx, filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the blog with the specified ID: %v", err))
	}
	return &pb.ReadBlogResponse{Blog: dataToBlogPb(data)}, nil
}

func dataToBlogPb(item *blogItem) *pb.Blog {
	return &pb.Blog{
		Id:       item.ID.Hex(),
		AuthorId: item.AuthorId,
		Title:    item.Title,
		Content:  item.Content,
	}
}
