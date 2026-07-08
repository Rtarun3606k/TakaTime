package exporter

import (
	"bytes"
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

type failingWriter struct{}

func (f failingWriter) Write(p []byte) (int, error) {
	return 0, errors.New("write failed")
}

func TestWriteCSV(t *testing.T) {
	rows := []LogRow{
		{
			Timestamp: "2025-01-01T00:00:00Z",
			Date:      "2025-01-01",
			FileName:  "main.go",
			Project:   "takatime",
			Language:  "Go",
			Editor:    "Neovim",
			OS:        "Linux",
			GitBranch: "main",
			Duration:  12.34,
		},
	}

	var buf bytes.Buffer

	err := WriteCSV(&buf, rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	out := buf.String()

	if !strings.Contains(out, "timestamp,date,file_name") {
		t.Fatal("missing csv header")
	}

	if !strings.Contains(out, "main.go") {
		t.Fatal("missing row data")
	}

	if !strings.Contains(out, "12.34") {
		t.Fatal("missing duration")
	}
}

func TestWriteCSVError(t *testing.T) {
	rows := []LogRow{
		{FileName: "main.go"},
	}

	err := WriteCSV(failingWriter{}, rows)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestWriteJSON(t *testing.T) {
	rows := []LogRow{
		{
			FileName: "main.go",
			Project:  "takatime",
		},
	}

	var buf bytes.Buffer

	if err := WriteJSON(&buf, rows); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var decoded []LogRow

	if err := json.Unmarshal(buf.Bytes(), &decoded); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if len(decoded) != 1 {
		t.Fatalf("expected 1 row, got %d", len(decoded))
	}

	if decoded[0].FileName != "main.go" {
		t.Fatalf("expected filename main.go, got %s", decoded[0].FileName)
	}

	if decoded[0].Project != "takatime" {
		t.Fatalf("expected project takatime, got %s", decoded[0].Project)
	}
}

func TestWriteJSONError(t *testing.T) {
	rows := []LogRow{
		{FileName: "main.go"},
	}

	err := WriteJSON(failingWriter{}, rows)
	if err == nil {
		t.Fatal("expected error")
	}
}

