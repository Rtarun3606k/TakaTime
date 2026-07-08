package utils

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetGitBranch(t *testing.T) {
	tmpDir := t.TempDir()

	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Skip("git not installed")
	}

	cmd = exec.Command("git", "checkout", "-b", "test-branch")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to create branch: %v", err)
	}

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
