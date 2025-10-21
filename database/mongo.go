package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Database

func ConnectMongoDB() {
	uri := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	if uri == "" {
		uri = "mongodb://localhost:27017"
		log.Println("⚠️  MONGO_URI tidak diset, gunakan default:", uri)
	}
	if dbName == "" {
		dbName = "alumni_db"
		log.Println("⚠️  MONGO_DB tidak diset, gunakan default:", dbName)
	}

	clientOpts := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("❌ Gagal koneksi MongoDB: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ Gagal ping MongoDB: %v", err)
	}

	MongoDB = client.Database(dbName)
	log.Println("✅ Berhasil konek MongoDB:", dbName)
}
