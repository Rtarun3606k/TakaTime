package debugger

import (
	"log"
	"os"
	"path/filepath"
)

func SetupLog() error {
	// 1. Get Home Dir
	homeDirLocation, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		return err
	}

	// 2. Construct Paths
	logDir := filepath.Join(homeDirLocation, ".takatime")
	logPath := filepath.Join(logDir, "debug-logs.log")

	// 3. Ensure Directory Exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Error creating log directory: %v", err)
		return err
	}

	// 4. Open/Create File
	// O_CREATE = Create if missing
	// O_APPEND = Add to bottom if exists
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return err
	}

	// 5. Hook up standard logger
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	return nil
}

// debugger logging frog dashboard
func SetupDashboardLog() (*os.File, error) {
	// 1. Get Home Dir
	homeDirLocation, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Error getting user home directory: %v", err)
		return nil, err
	}

	// 2. Construct Paths
	logDir := filepath.Join(homeDirLocation, ".takatime")
	logPath := filepath.Join(logDir, "debug-logs-dashboard.log")

	// 3. Ensure Directory Exists
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Error creating log directory: %v", err)
		return nil, err
	}

	// 4. Open/Create File
	// O_CREATE = Create if missing
	// O_APPEND = Add to bottom if exists
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("Error opening log file: %v", err)
		return nil, err
	}

	// 5. Hook up standard logger
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	log.SetPrefix("[DASHBOARD] ")

	return file, nil
}
