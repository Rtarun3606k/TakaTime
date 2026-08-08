package debugger

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// setupIsolatedEnv creates a temporary mock home directory to prevent tests
// from overwriting the developer's actual ~/.takatime log files.
// It also saves and restores the global log state.
func setupIsolatedEnv(t *testing.T) string {
	t.Helper()

	// 1. Create a safe temporary directory
	tempDir := t.TempDir()

	// 2. Mock the environment variables os.UserHomeDir() relies on
	t.Setenv("HOME", tempDir)        // For Unix/Linux/macOS
	t.Setenv("USERPROFILE", tempDir) // For Windows

	// 3. Save the original log state so we can restore it after the test
	originalWriter := log.Writer()
	originalFlags := log.Flags()
	originalPrefix := log.Prefix()

	t.Cleanup(func() {
		log.SetOutput(originalWriter)
		log.SetFlags(originalFlags)
		log.SetPrefix(originalPrefix)
	})

	return tempDir
}

func TestSetupLog(t *testing.T) {
	// Fake the home directory
	tempDir := setupIsolatedEnv(t)

	// Run the function
	err := SetupLog()
	if err != nil {
		t.Fatalf("SetupLog() unexpected error: %v", err)
	}

	// 1. Verify the directory and file were actually created
	expectedLogPath := filepath.Join(tempDir, ".takatime", "debug-logs.log")
	if _, err := os.Stat(expectedLogPath); os.IsNotExist(err) {
		t.Fatalf("SetupLog() failed to create log file at expected path: %s", expectedLogPath)
	}

	// 2. Verify that writing to the standard logger goes to our file
	testMessage := "Test standard log entry"
	log.Println(testMessage)

	content, err := os.ReadFile(expectedLogPath)
	if err != nil {
		t.Fatalf("Failed to read generated log file: %v", err)
	}

	if !strings.Contains(string(content), testMessage) {
		t.Errorf("Log file did not contain the expected test message. Got:\n%s", string(content))
	}
}

func TestSetupDashboardLog(t *testing.T) {
	// Fake the home directory
	tempDir := setupIsolatedEnv(t)

	// Run the function
	file, err := SetupDashboardLog()
	if err != nil {
		t.Fatalf("SetupDashboardLog() unexpected error: %v", err)
	}
	if file == nil {
		t.Fatal("SetupDashboardLog() returned a nil file pointer")
	}
	defer file.Close() // Clean up the file handle

	// 1. Verify the directory and file were actually created
	expectedLogPath := filepath.Join(tempDir, ".takatime", "debug-logs-dashboard.log")
	if _, err := os.Stat(expectedLogPath); os.IsNotExist(err) {
		t.Fatalf("SetupDashboardLog() failed to create log file at expected path: %s", expectedLogPath)
	}

	// 2. Verify logging works AND uses the correct prefix
	testMessage := "Test dashboard entry"
	log.Println(testMessage)

	content, err := os.ReadFile(expectedLogPath)
	if err != nil {
		t.Fatalf("Failed to read generated log file: %v", err)
	}

	contentStr := string(content)

	if !strings.Contains(contentStr, "[DASHBOARD]") {
		t.Errorf("Log file missing the expected [DASHBOARD] prefix. Got:\n%s", contentStr)
	}

	if !strings.Contains(contentStr, testMessage) {
		t.Errorf("Log file did not contain the expected test message. Got:\n%s", contentStr)
	}
}

