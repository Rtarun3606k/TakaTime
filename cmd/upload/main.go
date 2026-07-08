package main

import (
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

// pes2ug23cs645
func main() {
	uri := flag.String("uri", "", "MongoDB Atlas Connection URI")
	project := flag.String("project", "unknown", "Project Name")
	file := flag.String("file", "", "File Name")
	duration := flag.Float64("duration", 0, "Duration in seconds")
	language := flag.String("language", "unknown", "Language (Deprecated)")
	editor := flag.String("editor", "unknown", "Editor Name NeoVim/VsCode")
	versionFlag := flag.Bool("version", false, "show Version")

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

	if *language != "unknown" {
		log.Println("The flag language is deprecated. The lang now is detected from the file. Lang provided:", *language)
	}

	*language = utils.DetectLanguage(*file)

	// Fallback
	if *language == "" {
		*language = "text"
	}
	log.Println("Detected language:", *language, "from file:", *file)

	var errr error
	types.DB, errr = db.InitSQLite()
	if errr != nil {
		log.Fatal("Could not initialize local DB:", errr)
	}
	defer types.DB.Close()

	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	//
	// client, err := db.ConnectToDataBase(*uri)
	// if err != nil {
	// 	log.Fatalln("Counld not connect to mongo db", err)
	// }
	//
	// collection := client.Database("takatime").Collection("logs")
	//
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
	}
	log.Printf("file: %s , project : %s , duration : %f , Language : %s , Gitbrach : %s , editor : %s ", *file, *project, duration, *language, gitBranch, *editor)

	// _, err = collection.InsertOne(ctx, entry)

	// 6. STEP 1: Always Save to Local DB First (Safety Net)
	if err := db.Enqueue(entry, types.DB); err != nil {
		log.Printf("Failed to save offline: %v", err)
		// If we can't save to disk, we probably shouldn't continue
		return
	}
	log.Printf("Saved log for '%s' to offline queue.", *file)

	// 7. STEP 2: The Sync Loop (Drain the Queue)
	// We assume *uri is valid here. If empty, we just skip syncing.
	if *uri != "" {
		db.SyncQueue(*uri, types.DB)
	}

}
