package dbqueryv2

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// RunMigrations checks if we need to clean up bad data
func RunMigrations(client *mongo.Client) error {
	db := client.Database("takatime")
	migrations := db.Collection("migrations")
	logs := db.Collection("logs")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 1. Check if we already ran the cleanup
	// We look for a document with _id: "v2_1_cleanup"
	count, err := migrations.CountDocuments(ctx, bson.D{{Key: "_id", Value: "v2_1_cleanup"}})
	if err != nil {
		return err
	}

	if count > 0 {
		// Migration already done. Skip.
		return nil
	}

	fmt.Println("🧹 TakaTime: Running one-time database cleanup...")

	// 2. THE CLEANUP LOGIC
	// Delete all logs caused by the "Exit Handler" bug
	deleteFilter := bson.D{
		{Key: "project", Value: "unknown"},
		{Key: "language", Value: "text"},
	}

	res, err := logs.DeleteMany(ctx, deleteFilter)
	if err != nil {
		log.Printf("Error deleting buggy logs: %v", err)
		return fmt.Errorf("failed to delete bad logs: %v", err)
	}

	fmt.Printf("✨ Deleted %d buggy logs.\n", res.DeletedCount)

	// 3. Mark as Done
	// Insert the marker so we don't run this again
	_, err = migrations.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: "v2_1_cleanup"}}, // Query
		bson.D{{Key: "$set", Value: bson.D{ // Update
			{Key: "ran_at", Value: time.Now()},
			{Key: "deleted_count", Value: res.DeletedCount},
		}}},
		options.UpdateOne().SetUpsert(true), // Create if not exists
	)

	return err
}
