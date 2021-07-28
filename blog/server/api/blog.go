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
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorId string             `bson:"author_id,omitempty"`
	Content  string             `bson:"content,omitempty"`
	Title    string             `bson:"title,omitempty"`
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

func (*Server) UpdateBlog(ctx context.Context, in *pb.UpdateBlogRequest) (*pb.UpdateBlogResponse, error) {
	log.Println("Updating blog")
	blogId := in.Blog.GetId()
	oid, err := primitive.ObjectIDFromHex(blogId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error parsing the informed ID: %v", err))
	}

	data := &blogItem{}
	filter := bson.M{"_id": oid}
	blogItem := db.Collection("blog").FindOne(ctx, filter)
	if err := blogItem.Decode(data); err != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Could not find the blog with the specified ID: %v", err))
	}
	blog := in.GetBlog()
	data.Title = blog.Title
	data.Content = blog.Content
	data.AuthorId = blog.AuthorId

	_, updtErr := db.Collection("blog").UpdateOne(context.Background(), filter, bson.M{"$set": data})
	if updtErr != nil {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Unable to update blog with the specified ID: %v", updtErr))
	}

	return &pb.UpdateBlogResponse{Blog: dataToBlogPb(data)}, nil
}

func (*Server) DeleteBlog(ctx context.Context, in *pb.DeleteBlogRequest) (*pb.DeleteBlogResponse, error) {
	log.Println("Deleting Blog Item...")
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error parsing the informed ID: %v", err))
	}

	res, dErr := db.Collection("blog").DeleteOne(context.Background(), bson.D{{"_id", oid}})
	if dErr != nil {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Error deleting the item: %v", dErr))
	}
	log.Println("Deleted: ", res.DeletedCount)
	if res.DeletedCount == 0 {
		return nil, status.Errorf(codes.NotFound, fmt.Sprintf("Item not found: %v", dErr))
	}
	return &pb.DeleteBlogResponse{Id: oid.Hex()}, nil
}
