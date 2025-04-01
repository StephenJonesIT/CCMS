package config

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB() *mongo.Client {
	dns := fmt.Sprintf("mongodb+srv://%s:%s@%s",
			os.Getenv("MONGODB_USERNAME"),
			os.Getenv("MONGODB_PASSWORD"),
			os.Getenv("MONGGODB_URL"))
	client, err := mongo.NewClient(options.Client().ApplyURI(dns))

	if err != nil {
		log.Error(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Error(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Error(err)
	}

	log.Info("Connected to MongoDB")
	return client
}


func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(os.Getenv("MONGODB_DB")).Collection(collectionName)
}