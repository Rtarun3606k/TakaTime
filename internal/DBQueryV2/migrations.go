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

// HasMigrationRun checks the migrations collection to see if a specific migration ID exists.
func HasMigrationRun(ctx context.Context, db *mongo.Database, migrationID string) (bool, error) {
	migrations := db.Collection("migrations")
	count, err := migrations.CountDocuments(ctx, bson.D{{Key: "_id", Value: migrationID}})
	if err != nil {
		return false, fmt.Errorf("failed to check status for %s: %w", migrationID, err)
	}
	return count > 0, nil
}

// MarkMigrationDone saves the checkpoint so the migration never runs again.
func MarkMigrationDone(ctx context.Context, db *mongo.Database, migrationID string, metadata bson.D) error {
	migrations := db.Collection("migrations")

	// Always save the exact time it ran
	updateFields := bson.D{{Key: "ran_at", Value: time.Now()}}

	// Append any extra metadata passed in
	if len(metadata) > 0 {
		updateFields = append(updateFields, metadata...)
	}

	_, err := migrations.UpdateOne(
		ctx,
		bson.D{{Key: "_id", Value: migrationID}},
		bson.D{{Key: "$set", Value: updateFields}},
		options.UpdateOne().SetUpsert(true),
	)
	return err
}

// RunMigrations checks if we need to clean up bad data
func RunMigrations(client *mongo.Client) error {
	db := client.Database("takatime")
	logs := db.Collection("logs")
	migrationID := "v2_1_cleanup"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	alreadyRan, err := HasMigrationRun(ctx, db, migrationID)
	if err != nil {
		return err
	}
	if alreadyRan {
		return nil // Skip
	}

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

	fmt.Printf("Deleted %d buggy logs.\n", res.DeletedCount)

	//markmigration
	metadata := bson.D{{
		Key: "deleted_count", Value: res.DeletedCount,
	}}

	return MarkMigrationDone(ctx, db, migrationID, metadata)

}

// RunLanguageMigration scans historical records and normalizes the "language" field using enry/overrides
func RunLanguageMigration(mongoClient *mongo.Client, dbName, collectionName string) error {
	db := mongoClient.Database(dbName)
	collection := db.Collection(collectionName)
	migrationID := "language_normalization_v1"

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	alreadyRan, err := HasMigrationRun(ctx, db, migrationID)
	if err != nil {
		return err
	}
	if alreadyRan {
		return nil // Skip silently
	}

	log.Println("TakaTime: Running one-time language normalization...")

	// 1. Fetch all distinct values from the "language" field (e.g., ["go", "javascriptreact", "JAVA"])
	results := collection.Distinct(ctx, "language", bson.D{})

	var languages []any

	if err := results.Decode(&languages); err != nil {
		return fmt.Errorf("failed to decode distinct languages: %w", err)
	}

	log.Printf("Found %d distinct language variations in MongoDB. Starting batch normalization...", len(languages))

	updatedCategories := 0
	totalLogsUpdated := int64(0)

	// 2. Iterate through the distinct names found in your database
	for _, rawLang := range languages {
		oldLangStr, ok := rawLang.(string)
		if !ok {
			continue // Skip if null or structurally unexpected
		}

		// 3. Pass it to the function that uses enry.GetLanguageByAlias / Extension / Overrides
		cleanLang := CleanTelemetryLanguage(oldLangStr)

		// 4. Update documents in bulk if the string changed (e.g., "go" -> "Go")
		if oldLangStr != cleanLang {
			filter := bson.M{"language": oldLangStr}
			update := bson.M{"$set": bson.M{"language": cleanLang}}

			updateResult, err := collection.UpdateMany(ctx, filter, update)
			if err != nil {
				log.Printf("Failed to migrate '%s' to '%s': %v", oldLangStr, cleanLang, err)
				continue
			}

			log.Printf("Consolidated: '%s' ➔ '%s' (%d logs updated)", oldLangStr, cleanLang, updateResult.ModifiedCount)
			updatedCategories++
			totalLogsUpdated += updateResult.ModifiedCount
		}
	}

	log.Printf("Migration finished cleanly. Consolidated %d fragmented categories.", updatedCategories)

	//mark migration
	metadata := bson.D{
		{Key: "categories_consolidated", Value: updatedCategories},
		{Key: "total_logs_updated", Value: totalLogsUpdated},
	}
	return MarkMigrationDone(ctx, db, migrationID, metadata)
}
