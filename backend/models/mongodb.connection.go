package models

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx2 = context.TODO()

func GetMongoModel() *mongo.Client {
	clientOpts := options.Client().ApplyURI(os.Getenv("MONGODB_HOST"))
	_mongo, errConnect := mongo.Connect(ctx2, clientOpts)
	if errConnect != nil {
		log.Fatal(errConnect)
	}
	return _mongo
}
