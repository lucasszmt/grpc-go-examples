package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

var (
	uri = "mongodb://root:example@localhost:27017"
	MongoClient *mongo.Client
)

func init() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	c, err :=mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}
	MongoClient = c
	conErr := c.Connect(ctx)
	if conErr != nil {
		log.Fatalln(conErr)
	}
}

func Disconnect() error{
	err := MongoClient.Disconnect(context.Background())
	if err != nil {
		return err
	}
	return nil
}
