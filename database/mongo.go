package database

import (
    "context"
    "fmt"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var MongoDB *mongo.Database

func ConnectMongo(uri string, dbName string) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    clientOpts := options.Client().ApplyURI(uri)
    client, err := mongo.Connect(ctx, clientOpts)
    if err != nil {
        log.Fatalf("error connecting mongodb: %v", err)
    }

    err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatalf("error ping mongodb: %v", err)
    }

    MongoClient = client
    MongoDB = client.Database(dbName)

    fmt.Println("✅ MongoDB connected:", uri, "db:", dbName)
}

func CloseMongo() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if MongoClient != nil {
        _ = MongoClient.Disconnect(ctx)
    }
}
