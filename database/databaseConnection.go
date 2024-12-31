package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client {
    err := godotenv.Load(".env")
    
    if err != nil {
        log.Fatal("Error loading the .env file")
    }

    MongoDb := os.Getenv("MONGODB_URL")

    clientOptions := options.Client().ApplyURI(MongoDb)

    client, err := mongo.Connect(context.Background(), clientOptions)

    if err != nil { 
        log.Fatal("Cannot connect to the database", err)
    }

    err = client.Ping(context.Background(), nil)
    if err != nil {
        log.Fatal("Cannot connect to the database", err)
    }
    fmt.Println("Connected to MongoDB successfully")


    return client

}


var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client, collectionName string) *mongo.Collection{
    var collection *mongo.Collection = client.Database("cluster0").Collection(collectionName)
    return collection
}
