package db

import (
	"context"
	"log"
	"time"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func AddEntryToMongo(logs []types.LogEntry, mongoURI string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOption := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(clientOption)

	if err != nil {
		log.Printf("Cannot connect to mongodb %s", err)
		return err
	}

	defer client.Disconnect(ctx)

	collection := client.Database("takatime").Collection("logs")
	var documents []any
	//add logs
	for _, log := range logs {
		documents = append(documents, log)
	}

	_, err = collection.InsertMany(ctx, documents)
	if err != nil {
		log.Printf("error: failed to insert logs into mongodb: %v", err)
		return err
	}

	return nil
}
