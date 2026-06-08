package dbqueryv2

import (
	"testing"
	"time"
)

func TestCalculateTrueDuration(t *testing.T) {
	base := time.Unix(1000, 0)

	tests := []struct {
		name     string
		logs     []LogTime
		expected float64
	}{
		{
			name: "single log",
			logs: []LogTime{
				{Timestamp: base, Duration: 60},
			},
			expected: 60,
		},
		{
			name: "overlapping logs",
			logs: []LogTime{
				{Timestamp: base, Duration: 60},
				{Timestamp: base.Add(30 * time.Second), Duration: 60},
			},
			expected: 90,
		},
		{
			name: "non-overlapping logs",
			logs: []LogTime{
				{Timestamp: base, Duration: 60},
				{Timestamp: base.Add(120 * time.Second), Duration: 60},
			},
			expected: 120,
		},
		{
			name:     "empty logs",
			logs:     []LogTime{},
			expected: 0,
		},
		{
			name: "contained overlap",
			logs: []LogTime{
				{Timestamp: base, Duration: 120},
				{Timestamp: base.Add(-30 * time.Second), Duration: 30},
			},
			expected: 120,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateTrueDuration(tt.logs)

			if got != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  float64
		expected string
	}{
		{
			name:     "minutes only",
			seconds:  300,
			expected: "5m",
		},
		{
			name:     "hours and minutes",
			seconds:  3660,
			expected: "1h 1m",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.seconds)
			if got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestGetYesterdayBounds(t *testing.T) {
	// Use a fixed non-UTC timezone to verify local midnight boundaries
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Skip("timezone not available")
	}

	tests := []struct {
		name          string
		now           time.Time
		expectedStart time.Time
		expectedEnd   time.Time
	}{
		{
			name:          "midday in Kolkata",
			now:           time.Date(2026, 6, 8, 14, 30, 0, 0, loc),
			expectedStart: time.Date(2026, 6, 7, 0, 0, 0, 0, loc),
			expectedEnd:   time.Date(2026, 6, 8, 0, 0, 0, 0, loc),
		},
		{
			name:          "just after midnight in Kolkata",
			now:           time.Date(2026, 6, 8, 0, 5, 0, 0, loc),
			expectedStart: time.Date(2026, 6, 7, 0, 0, 0, 0, loc),
			expectedEnd:   time.Date(2026, 6, 8, 0, 0, 0, 0, loc),
		},
		{
			name:          "log at yesterday start is included",
			now:           time.Date(2026, 6, 8, 12, 0, 0, 0, loc),
			expectedStart: time.Date(2026, 6, 7, 0, 0, 0, 0, loc),
			expectedEnd:   time.Date(2026, 6, 8, 0, 0, 0, 0, loc),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			start, end := getYesterdayBounds(tt.now)
			if !start.Equal(tt.expectedStart) {
				t.Errorf("start: expected %v, got %v", tt.expectedStart, start)
			}
			if !end.Equal(tt.expectedEnd) {
				t.Errorf("end: expected %v, got %v", tt.expectedEnd, end)
			}
		})
	}
}

func TestGetYesterdayBounds_LogFiltering(t *testing.T) {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		t.Skip("timezone not available")
	}

	now := time.Date(2026, 6, 8, 14, 0, 0, 0, loc)
	yesterdayStart, yesterdayEnd := getYesterdayBounds(now)

	logs := []LogTime{
		// Before yesterday start: should be excluded
		{Timestamp: time.Date(2026, 6, 6, 23, 59, 59, 0, loc), Duration: 60},
		// Exactly at yesterday start: should be included
		{Timestamp: time.Date(2026, 6, 7, 0, 0, 0, 0, loc), Duration: 60},
		// Inside yesterday: should be included
		{Timestamp: time.Date(2026, 6, 7, 14, 30, 0, 0, loc), Duration: 60},
		// Exactly at today's midnight (yesterday end): should be excluded
		{Timestamp: time.Date(2026, 6, 8, 0, 0, 0, 0, loc), Duration: 60},
		// After today's midnight: should be excluded
		{Timestamp: time.Date(2026, 6, 8, 12, 0, 0, 0, loc), Duration: 60},
	}

	var filtered []LogTime
	for _, log := range logs {
		if !log.Timestamp.Before(yesterdayStart) && log.Timestamp.Before(yesterdayEnd) {
			filtered = append(filtered, log)
		}
	}

	if len(filtered) != 2 {
		t.Fatalf("expected 2 logs in yesterday bucket, got %d", len(filtered))
	}

	// First included log: exactly at yesterday start
	if !filtered[0].Timestamp.Equal(time.Date(2026, 6, 7, 0, 0, 0, 0, loc)) {
		t.Errorf("first log should be at yesterday start, got %v", filtered[0].Timestamp)
	}

	// Second included log: inside yesterday
	if !filtered[1].Timestamp.Equal(time.Date(2026, 6, 7, 14, 30, 0, 0, loc)) {
		t.Errorf("second log should be inside yesterday, got %v", filtered[1].Timestamp)
	}
}
