package utils

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFromHeuristics(t *testing.T) {
	tests := []struct {
		name        string
		ext         string
		content     []byte
		wantMatched bool
	}{
		{
			name: "Perfect Match",
			ext:  "m",
			// Matches the Objective-C #import regex
			content:     []byte("#import <Foundation/Foundation.h>\n"),
			wantMatched: true,
		},
		{
			name:        "Negative Pattern",
			ext:         "pl",
			content:     []byte("<?xml version=\"1.0\"?>"),
			wantMatched: false,
		},
		{
			name: "Best Partial Match",
			ext:  "pm",
			// Matches the Perl "use strict" regex
			content:     []byte("use strict;\n"),
			wantMatched: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			langs, matched := DetectFromHeuristics(tt.ext, tt.content)

			if matched != tt.wantMatched {
				t.Errorf("DetectFromHeuristics() matched = %v, want %v", matched, tt.wantMatched)
			}

			if tt.wantMatched && len(langs) == 0 {
				t.Errorf("DetectFromHeuristics() expected language candidates, got none")
			}

			if !tt.wantMatched && tt.name == "Unknown Extension" && langs != nil {
				t.Errorf("DetectFromHeuristics() expected nil langs for unknown extension, got %v", langs)
			}
		})
	}
}

func TestDetectLanguage(t *testing.T) {
	dir := t.TempDir()

	tests := []struct {
		name          string
		filename      string
		content       []byte
		setupFile     bool
		wantDetected  bool
		wantCandidate bool
		wantErr       bool
	}{
		{
			name:          "Layer 1: Exact Filename Match",
			filename:      "Makefile", // Change this from Dockerfile
			content:       []byte("build:\n\tgo build"),
			setupFile:     true,
			wantDetected:  true,
			wantCandidate: false,
			wantErr:       false,
		},
		{
			name:          "Layer 2: Unique Extension Match",
			filename:      "main.go",
			content:       []byte("package main"),
			setupFile:     true,
			wantDetected:  true,
			wantCandidate: false,
			wantErr:       false,
		},
		{
			name:          "Layer 3: Unknown Extension",
			filename:      "a.xyz123",
			content:       []byte("abc"),
			setupFile:     true,
			wantDetected:  false,
			wantCandidate: false,
			wantErr:       false,
		},
		{
			name:          "Layer 4: Heuristics Match",
			filename:      "script.m",
			content:       []byte("#import <Foundation/Foundation.h>\n"), // Valid Obj-C
			setupFile:     true,
			wantDetected:  true,
			wantCandidate: true,
			wantErr:       false,
		},
		{
			name:          "I/O Error: Unreadable Ambiguous File",
			filename:      "missing.m",
			content:       nil,
			setupFile:     false,
			wantDetected:  false,
			wantCandidate: false,
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(dir, tt.filename)

			if tt.setupFile {
				err := os.WriteFile(path, tt.content, 0644)
				if err != nil {
					t.Fatalf("Failed to create temp file: %v", err)
				}
			}

			detected, _, candidates, err := DetectLanguage(path)

			if (err != nil) != tt.wantErr {
				t.Errorf("DetectLanguage() error = %v, wantErr %v", err, tt.wantErr)
			}
			if detected != tt.wantDetected {
				t.Errorf("DetectLanguage() detected = %v, want %v", detected, tt.wantDetected)
			}
			if tt.wantCandidate && len(candidates) == 0 {
				t.Errorf("DetectLanguage() expected candidates array, got empty/nil")
			}
		})
	}
}

func TestReadFile(t *testing.T) {
	dir := t.TempDir()

	t.Run("Successful Read", func(t *testing.T) {
		path := filepath.Join(dir, "test.txt")
		expected := []byte("hello world")

		if err := os.WriteFile(path, expected, 0644); err != nil {
			t.Fatalf("Failed to write temp file: %v", err)
		}

		got, err := ReadFile(path)
		if err != nil {
			t.Errorf("ReadFile() unexpected error: %v", err)
		}
		if !bytes.Equal(got, expected) {
			t.Errorf("ReadFile() got = %s, want %s", got, expected)
		}
	})

	t.Run("File Not Found", func(t *testing.T) {
		_, err := ReadFile(filepath.Join(dir, "missing.txt"))
		if err == nil {
			t.Errorf("ReadFile() expected error for missing file, got nil")
		}
	})

	t.Run("Enforces 4096 Byte Limit", func(t *testing.T) {
		path := filepath.Join(dir, "large.txt")
		data := bytes.Repeat([]byte("a"), 6000)

		if err := os.WriteFile(path, data, 0644); err != nil {
			t.Fatalf("Failed to write large temp file: %v", err)
		}

		got, err := ReadFile(path)
		if err != nil {
			t.Errorf("ReadFile() unexpected error: %v", err)
		}
		if len(got) != 4096 {
			t.Errorf("ReadFile() read %d bytes, want exactly 4096 bytes", len(got))
		}
	})
}

