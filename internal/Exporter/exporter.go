// Package exporter provides data-fetching and serialisation helpers for the
// taka-export CLI binary.  It is intentionally kept separate from the rest of
// the TakaTime internals so the export concern stays self-contained and easy
// to test without pulling in image-generation or dashboard dependencies.
package exporter

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sort"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ── Data types ────────────────────────────────────────────────────────────────

// LogRow is a single raw log entry in export-friendly form.
// The Timestamp field is serialised as RFC3339 UTC so downstream tools
// (pandas, Excel) parse it unambiguously.
type LogRow struct {
	Timestamp string  `json:"timestamp"        csv:"timestamp"`
	Date      string  `json:"date"             csv:"date"`
	FileName  string  `json:"file_name"        csv:"file_name"`
	Project   string  `json:"project"          csv:"project"`
	Language  string  `json:"language"         csv:"language"`
	Editor    string  `json:"editor"           csv:"editor"`
	OS        string  `json:"os"               csv:"os"`
	GitBranch string  `json:"git_branch"       csv:"git_branch"`
	Duration  float64 `json:"duration_seconds" csv:"duration_seconds"`
}

// rawDoc mirrors the MongoDB document structure; used only during the fetch.
type rawDoc struct {
	Timestamp time.Time `bson:"timestamp"`
	Date      string    `bson:"date"`
	FileName  string    `bson:"name"`
	Project   string    `bson:"project"`
	Language  string    `bson:"language"`
	Editor    string    `bson:"editor"`
	OS        string    `bson:"os"`
	GitBranch string    `bson:"gitBranch"`
	Duration  float64   `bson:"duration"`
}

// FilterOptions controls which documents are returned by FetchAllLogs.
type FilterOptions struct {
	// From and To are inclusive date boundaries (zero value = unbounded).
	From time.Time
	To   time.Time
}

// ── Fetch ─────────────────────────────────────────────────────────────────────

// FetchAllLogs queries the takatime.logs collection and returns every matching
// document as a slice of LogRow values, ordered by timestamp ascending.
//
// Both FilterOptions boundaries are inclusive and optional:
//   - From zero  → no lower bound
//   - To   zero  → no upper bound
func FetchAllLogs(ctx context.Context, client *mongo.Client, opts FilterOptions) ([]LogRow, error) {
	collection := client.Database("takatime").Collection("logs")

	// Build the date-range filter only when bounds are actually set.
	tsFilter := bson.D{}
	if !opts.From.IsZero() {
		tsFilter = append(tsFilter, bson.E{Key: "$gte", Value: opts.From})
	}
	if !opts.To.IsZero() {
		// Make --to inclusive by advancing to the end of that calendar day.
		inclusive := opts.To.Truncate(24 * time.Hour).Add(24*time.Hour - time.Nanosecond)
		tsFilter = append(tsFilter, bson.E{Key: "$lte", Value: inclusive})
	}

	filter := bson.D{}
	if len(tsFilter) > 0 {
		filter = append(filter, bson.E{Key: "timestamp", Value: tsFilter})
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Printf("Error fetching logs from MongoDB: %v", err)
		return nil, fmt.Errorf("mongo find: %w", err)
	}
	defer cursor.Close(ctx)

	var docs []rawDoc
	if err := cursor.All(ctx, &docs); err != nil {
		log.Printf("Error decoding logs from MongoDB: %v", err)

		return nil, fmt.Errorf("cursor decode: %w", err)
	}

	// Sort ascending by timestamp so CSV/JSON output is chronological.
	sort.Slice(docs, func(i, j int) bool {
		return docs[i].Timestamp.Before(docs[j].Timestamp)
	})

	rows := make([]LogRow, len(docs))
	for i, d := range docs {
		rows[i] = LogRow{
			Timestamp: d.Timestamp.UTC().Format(time.RFC3339),
			Date:      d.Date,
			FileName:  d.FileName,
			Project:   d.Project,
			Language:  d.Language,
			Editor:    d.Editor,
			OS:        d.OS,
			GitBranch: d.GitBranch,
			Duration:  d.Duration,
		}
	}
	return rows, nil
}

// ── Writers ───────────────────────────────────────────────────────────────────

// csvHeader defines the canonical column order, matching the issue spec.
var csvHeader = []string{
	"timestamp",
	"date",
	"file_name",
	"project",
	"language",
	"editor",
	"os",
	"git_branch",
	"duration_seconds",
}

// WriteCSV serialises rows into RFC 4180-compliant CSV, safe to open directly
// in Excel and load with pandas (pd.read_csv).
func WriteCSV(w io.Writer, rows []LogRow) error {
	cw := csv.NewWriter(w)

	if err := cw.Write(csvHeader); err != nil {
		log.Printf("Error writing CSV header: %v", err)
		return fmt.Errorf("csv header write: %w", err)
	}

	for _, r := range rows {
		record := []string{
			r.Timestamp,
			r.Date,
			r.FileName,
			r.Project,
			r.Language,
			r.Editor,
			r.OS,
			r.GitBranch,
			strconv.FormatFloat(r.Duration, 'f', 2, 64),
		}
		if err := cw.Write(record); err != nil {
			return fmt.Errorf("csv row write: %w", err)
		}
	}

	cw.Flush()
	return cw.Error()
}

// WriteJSON serialises rows into a pretty-printed JSON array.
func WriteJSON(w io.Writer, rows []LogRow) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rows); err != nil {
		return fmt.Errorf("json encode: %w", err)
	}
	return nil
}
