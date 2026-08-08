package persnalization

import (
	"strings"
	"testing"

	"github.com/Rtarun3606k/TakaTime/internal/Styles"
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/charmbracelet/lipgloss"
)

func TestGetCoderPersona(t *testing.T) {
	tests := []struct {
		name string
		dist types.ActivityDistribution
		want string
	}{
		{
			name: "Morning Max",
			dist: types.ActivityDistribution{Morning: 10, Afternoon: 5, Evening: 2, Night: 1},
			want: "🌅  Early Bird",
		},
		{
			name: "Afternoon Max",
			dist: types.ActivityDistribution{Morning: 2, Afternoon: 12, Evening: 5, Night: 0},
			want: "☀️  Afternoon Architect",
		},
		{
			name: "Evening Max",
			dist: types.ActivityDistribution{Morning: 0, Afternoon: 4, Evening: 8, Night: 3},
			want: "🌆  Evening Engineer",
		},
		{
			name: "Night Max",
			dist: types.ActivityDistribution{Morning: 1, Afternoon: 1, Evening: 2, Night: 15},
			want: "🦉  Midnight Vampire",
		},
		{
			name: "All Zeros",
			dist: types.ActivityDistribution{Morning: 0, Afternoon: 0, Evening: 0, Night: 0},
			want: "💤 Resting Developer",
		},
		{
			name: "Ties Default to Earliest (Morning tie)",
			dist: types.ActivityDistribution{Morning: 5, Afternoon: 5, Evening: 0, Night: 0},
			want: "🌅  Early Bird",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCoderPersona(tt.dist); got != tt.want {
				t.Errorf("GetCoderPersona() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildStreakBox(t *testing.T) {
	// Mock basic styles for testing (colors don't matter for string containment)
	mockStyles := Styles.AppStyles{
		Color1:  lipgloss.NewStyle(),
		Color2:  lipgloss.NewStyle(),
		Text:    lipgloss.NewStyle(),
		SubText: lipgloss.NewStyle(),
		Box:     lipgloss.NewStyle(),
	}

	tests := []struct {
		name          string
		streak        int
		todayHours    float64
		avgHours      float64
		maxHours      float64
		maxDate       string
		width         int
		expectedTexts []string
	}{
		{
			name:       "Active Streak with Valid Max Date",
			streak:     5,
			todayHours: 4.5,
			avgHours:   3.0,
			maxHours:   8.2,
			maxDate:    "2026-04-03",
			width:      50,
			expectedTexts: []string{
				"🔥 5 Day Streak",
				"Today:  4.5h / Avg:  3.0h",
				"🏆 Record:  8.2h on Apr 03, 2026",
				"█", // Should have filled bar characters
			},
		},
		{
			name:       "No Streak and Fallback Date",
			streak:     0,
			todayHours: 1.0,
			avgHours:   2.5,
			maxHours:   0.0,
			maxDate:    "",
			width:      40,
			expectedTexts: []string{
				"❄️  No Active Streak",
				"Today:  1.0h / Avg:  2.5h",
				"🏆 Record: N/A",
				"░", // Should have empty bar characters
			},
		},
		{
			name:       "Average Zero Fallback (Prevents Div by Zero)",
			streak:     1,
			todayHours: 0.5,
			avgHours:   0.0, // Should trigger the <= 0.1 fallback to 1.0
			maxHours:   2.0,
			maxDate:    "2026-01-01",
			width:      40,
			expectedTexts: []string{
				"Today:  0.5h / Avg:  1.0h",
			},
		},
		{
			name:       "Over 100% Goal (Percent Cap)",
			streak:     10,
			todayHours: 10.0,
			avgHours:   2.0,
			maxHours:   10.0,
			maxDate:    "2026-08-01",
			width:      40,
			expectedTexts: []string{
				"Today: 10.0h / Avg:  2.0h",
			},
		},
		{
			name:       "Invalid Date Parsing Fallback",
			streak:     2,
			todayHours: 2.0,
			avgHours:   2.0,
			maxHours:   5.0,
			maxDate:    "invalid-date-format",
			width:      40,
			expectedTexts: []string{
				"🏆 Record:  5.0h on invalid-date-format", // Should fallback to raw string
			},
		},
		{
			name:       "Extreme Width Clamping (Too Small)",
			streak:     1,
			todayHours: 1.0,
			avgHours:   1.0,
			maxHours:   0,
			maxDate:    "",
			width:      30, // Triggers barWidth < 10 fallback
			expectedTexts: []string{
				"█",
			},
		},
		{
			name:       "Extreme Width Clamping (Too Large)",
			streak:     1,
			todayHours: 1.0,
			avgHours:   1.0,
			maxHours:   0,
			maxDate:    "",
			width:      1000, // Triggers barWidth > 25 fallback
			expectedTexts: []string{
				"█",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BuildStreakBox(tt.streak, tt.todayHours, tt.avgHours, tt.maxHours, tt.maxDate, mockStyles, tt.width)

			for _, expected := range tt.expectedTexts {
				if !strings.Contains(got, expected) {
					t.Errorf("BuildStreakBox() output missing expected text: %q\nGot Output:\n%s", expected, got)
				}
			}

			// Ensure the overall header is always present
			if !strings.Contains(got, "━ Daily Target ━") {
				t.Errorf("BuildStreakBox() missing standard header")
			}
		})
	}
}

