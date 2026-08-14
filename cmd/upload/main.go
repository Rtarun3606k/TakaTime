package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"time"

	utils "github.com/Rtarun3606k/TakaTime/internal/Utils"
	"github.com/Rtarun3606k/TakaTime/internal/db"
	"github.com/Rtarun3606k/TakaTime/internal/debugger"
	"github.com/Rtarun3606k/TakaTime/internal/types"
)

func main() {
	uri := flag.String("uri", "", "MongoDB Atlas Connection URI")
	project := flag.String("project", "unknown", "Project Name")
	file := flag.String("file", "", "File Name")
	duration := flag.Float64("duration", 0, "Duration in seconds")
	language := flag.String("language", "unknown", "Language (Deprecated)")
	editor := flag.String("editor", "unknown", "Editor Name NeoVim/VsCode")
	versionFlag := flag.Bool("version", false, "show Version")
	aiAcceptedCount := flag.Int("ai-accepted", 0, "Total count of accepted AI suggestions")
	aiMetadataStr := flag.String("ai-metadata", "[]", "JSON string of AI completion events")

	flag.Parse()

	if *versionFlag {
		fmt.Println(types.Version)
		return
	}

	if *uri == "" || *duration <= 0 {
		log.Fatalln("Arguments are empty MongoDB URI or Duration is less than 0")
		return
	}

	// setup debugger logs
	err := debugger.SetupLog()
	if err != nil {
		log.Panic("Failed to initialize logger: ", err)
	}

	var aiEvents []types.AIMetadata
	if *aiMetadataStr != "" {
		err:= json.Unmarshal([]byte(*aiMetadataStr), &aiEvents)
		if err!=nil {
			return 
		}
	}
	for _, event := range aiEvents {
		matched, detectedLanguage, candidates, err := utils.DetectLanguage(event.FilePath)
		if err != nil {
			return
		}
		if matched {
			event.Language = detectedLanguage
			log.Printf("Detected language: %s", *language)
		} else if len(candidates) > 0 {
			event.Language = candidates[0]
			log.Printf("Heuristics could not determine the language. Falling back to: %s", *language)
			log.Printf("Possible languages: %v", candidates)
		} else {
			event.Language = "Unknown"
			log.Printf("Unknown language")
		}

	}

	if *language != "unknown" {
		log.Println("The flag language is deprecated. The lang now is detected from the file. Lang provided:", *language)
	}

	matched, detectedLanguage, candidates, err := utils.DetectLanguage(*file)
	if err != nil {
		return
	}

	if matched {
		*language = detectedLanguage
		log.Printf("Detected language: %s", *language)
	} else if len(candidates) > 0 {
		*language = candidates[0]
		log.Printf("Heuristics could not determine the language. Falling back to: %s", *language)
		log.Printf("Possible languages: %v", candidates)
	} else {
		*language = "Unknown"
		log.Printf("Unknown language")
	}

	var errr error
	types.DB, errr = db.InitSQLite()
	if errr != nil {
		log.Fatal("Could not initialize local DB:", errr)
	}
	defer types.DB.Close()

	fileDir := filepath.Dir(*file)

	gitBranch, err := utils.GetGitBranch(fileDir)
	if err != nil {
		log.Printf("Could not get git branch there might not be one or initiated yet !! %s", err)
		gitBranch = ""
	}

	entry := types.LogEntry{
		FileName:  *file,
		Project:   *project,
		Duration:  *duration,
		TimeStamp: time.Now(),
		Date:      time.Now().Format("2006-01-02"),
		Language:  *language,
		Os:        utils.GetOS(),
		GitBranch: gitBranch,
		Editor:    *editor,
		aiData:    aiEvents,
	}

	//  Always Save to Local DB First (Safety Net)
	if err := db.Enqueue(entry, types.DB); err != nil {
		log.Printf("Failed to save offline: %v", err)
		// If we can't save to disk, we probably shouldn't continue
		return
	}
	log.Printf("Saved log for '%s' to offline queue.", *file)

	// The Sync Loop (Drain the Queue)
	// We assume *uri is valid here. If empty, we just skip syncing.
	if *uri != "" {
		db.SyncQueue(*uri, types.DB)
	}
}
