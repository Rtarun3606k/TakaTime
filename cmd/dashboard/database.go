package main

import (
	"context"
	"log"
	"time"

	dbqueryv2 "github.com/Rtarun3606k/TakaTime/internal/DBQueryV2"
	"github.com/Rtarun3606k/TakaTime/internal/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var labels = [4]string{"language", "project", "os", "editor"}

func (m Model) GetData(URI string) (Model, *mongo.Client, error) {
	Client, err := db.ConnectToDataBase(URI)
	if err != nil {
		log.Println("Database connection failed:", err)
		return m, nil, err // Return the unmodified model on error
	}

	for _, value := range labels {
		data, err := dbqueryv2.GetListStats(Client, value, 30, m.TUITheme, 0)
		if err != nil {
			log.Println("Failed to fetch stats for", value, ":", err)
			return m, nil, err
		}


		// Assign the data directly to the model's fields
		switch value {
		case "language":
			m.LanguageListStats = data
		case "project":
			m.ProjectListStats = data
		case "os":
			m.OsListStats = data
		case "editor":
			m.editorListStats = data
		}

	}

	context, cancle := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancle()
	collection := Client.Database("takatime").Collection("logs")

	streak, todayHours, avgHours, dailyHistory, err := dbqueryv2.FetchStreakAndToday(context, collection)
	if err != nil {
		// handle error or just log it
		log.Printf("Error fetching streak: %v", err)
	} else {
		m.Streak = streak
		m.TodayHours = todayHours
		m.AverageHours = avgHours
		m.DailyHistory = dailyHistory
        log.Println("Streak fetched!")
	}

	// 2. Fetch Activity Distribution (Coder Persona)
	activityDist, err := dbqueryv2.FetchActivityDistribution(context, collection)
	if err != nil {
		log.Printf("Error fetching activity: %v", err)
	} else {
		m.ActivityData = activityDist
		log.Println("Activity Distribution fetched!")
	}

	//get time grid stats today yestarday all that

	timeStats, err := dbqueryv2.GetTimeStats(Client)
	if err != nil {
		log.Println("could not fetch timestats ", err)
	} else {
		log.Println("Timestats fetched!")
	}
	

	m.DataFetchedTime = time.Now().Add(-3 * time.Minute)
	m.TimeStats = timeStats
	// Return model
	return m, Client, nil
}
