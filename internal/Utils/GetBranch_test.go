package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetGitBranch(t *testing.T) {
	tmpDir := t.TempDir()

	// 1. Initialize the repository
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skip("git not installed")
	}

	// 2. Set dummy config to prevent CI failures where git user is not set
	exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()
	exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()

	// 3. Create an empty initial commit so HEAD actually exists
	commitCmd := exec.Command("git", "commit", "--allow-empty", "-m", "Initial commit")
	commitCmd.Dir = tmpDir
	if err := commitCmd.Run(); err != nil {
		t.Fatalf("failed to create initial commit: %v", err)
	}

	// 4. Create and checkout the test branch
	checkoutCmd := exec.Command("git", "checkout", "-b", "test-branch")
	checkoutCmd.Dir = tmpDir
	if err := checkoutCmd.Run(); err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}

	// 5. Run the actual test
	branch, err := GetGitBranch(tmpDir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if branch != "test-branch" {
		t.Fatalf("expected test-branch, got %s", branch)
	}
}

func TestGetGitBranch_InvalidDirectory(t *testing.T) {
	_, err := GetGitBranch("/path/that/does/not/exist")

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetGitBranch_NotGitRepo(t *testing.T) {
	tmpDir := t.TempDir()

	file := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(file, []byte("hello"), 0644); err != nil {
		t.Fatal(err)
	}

	_, err := GetGitBranch(tmpDir)

	if err == nil {
		t.Fatal("expected error for non-git directory")
	}
}
