package dbqueryv2

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DailyStat struct {
	Date         string  `bson:"_id"`          // Format: "YYYY-MM-DD"
	TotalSeconds float64 `bson:"totalSeconds"` // Sum of duration
}

// HourStat represents the aggregation result for a specific hour of the day
type HourStat struct {
	Hour         int     `bson:"_id"` // 0 to 23
	TotalSeconds float64 `bson:"totalSeconds"`
}

func FetchActivityDistribution(ctx context.Context, collection *mongo.Collection) (types.ActivityDistribution, error) {
	var dist types.ActivityDistribution
	// thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	//local time
	loc := time.Now().Location().String()
	if loc == "Local" || loc == "" {
		loc = "Asia/Kolkata"
	}
	// pipeline := mongo.Pipeline{
	// 	{{"$match", bson.D{{"timestamp", bson.D{{"$gte", thirtyDaysAgo}}}}}},
	// 	// Extract just the hour (0-23) in local time
	// 	{{"$group", bson.D{
	// 		{"_id", bson.D{
	// 			{"$hour", bson.D{
	// 				{"date", "$timestamp"},
	// 				{"timezone", loc},
	// 			}},
	// 		}},
	// 		{"totalSeconds", bson.D{{"$sum", "$duration"}}},
	// 	}}},
	// }
	pipeline := mongo.Pipeline{
		// NO $match stage here! It reads the whole database.
		{{"$group", bson.D{
			{"_id", bson.D{
				{"$hour", bson.D{
					{"date", "$timestamp"},
					{"timezone", loc},
				}},
			}},
			{"totalSeconds", bson.D{{"$sum", "$duration"}}},
		}}},
	}
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("Error running aggregation for time stats: %v", err)
		return dist, err
	}
	defer cursor.Close(ctx)

	var hourStats []HourStat
	if err = cursor.All(ctx, &hourStats); err != nil {
		log.Printf("Error decoding aggregation results for time stats: %v", err)
		return dist, err
	}

	// Group the 24 hours into our 4 buckets
	for _, stat := range hourStats {
		hoursCoded := stat.TotalSeconds / 3600.0

		if stat.Hour >= 6 && stat.Hour < 12 {
			dist.Morning += hoursCoded
		} else if stat.Hour >= 12 && stat.Hour < 18 {
			dist.Afternoon += hoursCoded
		} else if stat.Hour >= 18 && stat.Hour < 24 {
			dist.Evening += hoursCoded
		} else {
			// 00:00 to 05:59
			dist.Night += hoursCoded
		}
	}

	// Calculate the MaxVal so the UI progress bars scale perfectly
	dist.MaxVal = dist.Morning
	if dist.Afternoon > dist.MaxVal {
		dist.MaxVal = dist.Afternoon
	}
	if dist.Evening > dist.MaxVal {
		dist.MaxVal = dist.Evening
	}
	if dist.Night > dist.MaxVal {
		dist.MaxVal = dist.Night
	}
	log.Println(dist.MaxVal, dist.Afternoon, dist.Night, dist.Morning)
	fmt.Println(dist.MaxVal, dist.Afternoon, dist.Night, dist.Morning)
	return dist, nil
}

// Notice we now return (int, float64, float64, error) -> Streak, Today, Average, Error
func FetchStreakAndToday(ctx context.Context, collection *mongo.Collection) (int, float64, float64, map[string]float64, error) {
	loc := time.Now().Location().String()
	if loc == "Local" || loc == "" {
		loc = "Asia/Kolkata"
	}

	pipeline := mongo.Pipeline{
		// Removed the $match stage so it fetches ALL TIME history!
		{{"$group", bson.D{
			{"_id", bson.D{
				{"$dateToString", bson.D{
					{"format", "%Y-%m-%d"},
					{"date", "$timestamp"},
					{"timezone", loc},
				}},
			}},
			{"totalSeconds", bson.D{{"$sum", "$duration"}}},
		}}},
		{{"$sort", bson.D{{"_id", -1}}}},
	}

	cursor, err := collection.Aggregate(ctx, pipeline)
	dailyMap := make(map[string]float64)
	if err != nil {
		log.Printf("Error running aggregation for daily stats: %v", err)
		return 0, 0, 0, dailyMap, err
	}
	defer cursor.Close(ctx)

	var dailyStats []DailyStat
	if err = cursor.All(ctx, &dailyStats); err != nil {
		log.Printf("Error decoding aggregation results for daily stats: %v", err)
		return 0, 0, 0, dailyMap, err
	}

	todayStr := time.Now().Format("2006-01-02")
	yesterdayStr := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	var todayHours float64
	var totalAllTimeSeconds float64 // New variable to track total time
	streak := 0
	activeDayMap := make(map[string]bool)

	// Loop through every day you've ever coded
	for _, stat := range dailyStats {
		activeDayMap[stat.Date] = true
		totalAllTimeSeconds += stat.TotalSeconds // Add to the grand total
		//count git hub type blocks hours count
		hours := stat.TotalSeconds / 3600.0
		dailyMap[stat.Date] = hours

		//
		if stat.Date == todayStr {
			todayHours = stat.TotalSeconds / 3600.0
		}
	}

	// Calculate the All-Time Average!
	var averageHours float64
	if len(dailyStats) > 0 {
		// Total seconds / Number of active days / 3600 = Average Hours per day
		averageHours = (totalAllTimeSeconds / float64(len(dailyStats))) / 3600.0
	}

	// Calculate Streak
	checkDate := time.Now()
	if !activeDayMap[todayStr] && activeDayMap[yesterdayStr] {
		checkDate = checkDate.AddDate(0, 0, -1)
	}
	for {
		dateStr := checkDate.Format("2006-01-02")
		if activeDayMap[dateStr] {
			streak++
			checkDate = checkDate.AddDate(0, 0, -1)
		} else {
			break
		}
	}

	return streak, todayHours, averageHours, dailyMap, nil
}
