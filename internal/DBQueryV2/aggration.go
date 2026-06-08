package dbqueryv2

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// Helper struct for unmarshalling aggregation results
type StatResult struct {
	Name         string  `bson:"_id"`
	TotalSeconds float64 `bson:"totalSeconds"`
}
type LogTime struct {
	Timestamp time.Time `bson:"timestamp"`
	Duration  float64   `bson:"duration"`
}
type RawLog struct {
	Timestamp time.Time `bson:"timestamp"`
	Duration  float64   `bson:"duration"`
	Project   string    `bson:"project"`
	Language  string    `bson:"language"`
	OS        string    `bson:"os"`
	Editor    string    `bson:"editor"`
}

// 1. GENERIC STATS FETCHER (Now uses "Smart Merge" logic)
// This fixes the "33h Text vs 26h Total" bug by deduplicating time in Go.
func GetListStats(client *mongo.Client, fieldName string, limit int, theme types.ThemeConfig, days int) ([]types.ListStats, error) {
	collection := client.Database("takatime").Collection("logs")
	// Increase timeout since we are fetching more data to process in memory
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// 1. Fetch RAW logs where the field exists
	// all  because we need Timestamp & Duration for the merge algo no chnages required !!
	// filter := bson.D{{Key: fieldName, Value: bson.D{{Key: "$ne", Value: ""}}}}
	filter := bson.D{{Key: fieldName, Value: bson.D{{Key: "$ne", Value: ""}}}}

if days > 0 {
	start := time.Now().AddDate(0, 0, -days)
	filter = append(filter, bson.E{
		Key: "timestamp",
		Value: bson.D{{Key: "$gte", Value: start}},
	})
}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	// Temporary struct to capture fields needed for grouping
	//todo : add this in types in future
	// may be in next maintaince
	type RawLog struct {
		Timestamp time.Time `bson:"timestamp"`
		Duration  float64   `bson:"duration"`
		Project   string    `bson:"project"`
		Language  string    `bson:"language"`
		OS        string    `bson:"os"`
		Editor    string    `bson:"editor"`
	}

	var rawLogs []RawLog
	if err = cursor.All(ctx, &rawLogs); err != nil {
		return nil, err
	}

	groupedLogs := make(map[string][]LogTime)

	for _, log := range rawLogs {
		// Extract the key based on the requested fieldName
		key := ""
		switch fieldName {
		case "project":
			key = log.Project
		case "language":
			key = log.Language
		case "os":
			key = log.OS
		case "editor":
			key = log.Editor
		}

		if key == "" {
			continue
		}

		// Add to the bucket for this specific language/project
		groupedLogs[key] = append(groupedLogs[key], LogTime{
			Timestamp: log.Timestamp,
			Duration:  log.Duration,
		})
	}

	// 3. Calculate "True Duration" for each group
	// This runs your helper function on EACH language separately to fix overlaps
	//this did not work but something else worked but don't touch becsuse it works .. from past
	var results []StatResult
	var grandTotal float64

	for name, logs := range groupedLogs {
		trueDuration := calculateTrueDuration(logs)

		if trueDuration > 0 {
			results = append(results, StatResult{
				Name:         name,
				TotalSeconds: trueDuration,
			})
			grandTotal += trueDuration
		}
	}

	// 4. Sort by Duration (High to Low)
	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalSeconds > results[j].TotalSeconds
	})

	// 5. Convert to ListStats struct (UI Logic)
	var stats []types.ListStats
	colors := []string{theme.Color1, theme.Color2, theme.Color3, theme.Color4, theme.TextColor}

	for i, r := range results {
		if i >= limit {
			break // Only top N
		}

		color := colors[i%len(colors)]

		percent := 0.0
		if grandTotal > 0 {
			percent = r.TotalSeconds / grandTotal
		}

		stats = append(stats, types.ListStats{
			Label:   r.Name,
			Value:   formatDuration(r.TotalSeconds),
			Percent: percent,
			Color:   color,
		})
	}

	return stats, nil
}

// 2. TIME GRID FETCHER (Uses $facet for efficiency)
func GetTimeStats(client *mongo.Client) (types.TimeGridStruct, error) {
	collection := client.Database("takatime").Collection("logs")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) // Increased timeout for fetching data
	defer cancel()

	// 1. Fetch ALL logs (needed for accurate deduplication)
	// Optimization: If you have millions of logs, you might want to limit this,
	// but for a personal dashboard, fetching all is fine for accuracy.
	//already thought this through so no worries man its you from past ..
	cursor, err := collection.Find(ctx, bson.D{})
	if err != nil {
		return types.TimeGridStruct{}, err
	}

	var allLogs []LogTime
	if err = cursor.All(ctx, &allLogs); err != nil {
		return types.TimeGridStruct{}, err
	}

	// 2. Define Time Boundaries
	now := time.Now()
	yesterdayStart, yesterdayEnd := getYesterdayBounds(now)
	weekAgo := now.AddDate(0, 0, -7)
	monthAgo := now.AddDate(0, 0, -30)

	// 3. Filter Logs into buckets (in Memory)
	var logsYesterday, logsWeek, logsMonth []LogTime

	for _, log := range allLogs {
		// Yesterday: inclusive of start, exclusive of end
		if !log.Timestamp.Before(yesterdayStart) && log.Timestamp.Before(yesterdayEnd) {
			logsYesterday = append(logsYesterday, log)
		}
		// Last 7 Days
		if log.Timestamp.After(weekAgo) {
			logsWeek = append(logsWeek, log)
		}
		// Last 30 Days
		if log.Timestamp.After(monthAgo) {
			logsMonth = append(logsMonth, log)
		}
	}

	// 4. Calculate "True" Duration (Deduplicated)
	return types.TimeGridStruct{
		Yestarday: formatDuration(calculateTrueDuration(logsYesterday)),
		Week:      formatDuration(calculateTrueDuration(logsWeek)),
		Month:     formatDuration(calculateTrueDuration(logsMonth)),
		AllTime:   formatDuration(calculateTrueDuration(allLogs)),
	}, nil
}

// Helper: 3661s -> "1h 1m"
func formatDuration(seconds float64) string {
	d := time.Duration(seconds) * time.Second
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}
	return fmt.Sprintf("%dm", m)
}

// getYesterdayBounds returns the start and end boundaries for yesterday in the
// local timezone of the given reference time. Start is inclusive, end is exclusive.
func getYesterdayBounds(now time.Time) (start, end time.Time) {
	loc := now.Location()
	start = time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, loc)
	end = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	return
}

// CHANGED: Now accepts []LogTime instead of []struct{...}
func calculateTrueDuration(logs []LogTime) float64 {
	if len(logs) == 0 {
		return 0
	}

	type TimeRange struct {
		Start int64
		End   int64
	}

	var ranges []TimeRange
	for _, l := range logs {
		end := l.Timestamp.Unix()
		start := end - int64(l.Duration)
		ranges = append(ranges, TimeRange{Start: start, End: end})
	}

	sort.Slice(ranges, func(i, j int) bool {
		return ranges[i].Start < ranges[j].Start
	})

	var merged []TimeRange
	if len(ranges) > 0 {
		current := ranges[0]
		for _, next := range ranges[1:] {
			if next.Start < current.End {
				if next.End > current.End {
					current.End = next.End
				}
			} else {
				merged = append(merged, current)
				current = next
			}
		}
		merged = append(merged, current)
	}

	var totalSeconds int64
	for _, r := range merged {
		totalSeconds += (r.End - r.Start)
	}

	return float64(totalSeconds)
}
