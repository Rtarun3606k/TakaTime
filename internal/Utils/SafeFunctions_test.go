package utils

import (
	"testing"
)

func TestSafeTruncateString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "String Longer Than Max (Truncates and adds ...)",
			input:  "TakaTimeDashboard",
			maxLen: 8,
			want:   "TakaTime...", // 8 chars + 3 dots
		},
		{
			name:   "String Shorter Than Max (Pads with spaces)",
			input:  "Go",
			maxLen: 5,
			want:   "Go      ", // "Go" (2) + 6 spaces = 8 total (5 + 3)
		},
		{
			name:   "String Exactly Max Length (Pads with 3 spaces)",
			input:  "Hello",
			maxLen: 5,
			want:   "Hello   ", // "Hello" (5) + 3 spaces = 8 total
		},
		{
			name:   "Empty String (Pads to maxLen + 3)",
			input:  "",
			maxLen: 4,
			want:   "       ", // 0 chars + 7 spaces
		},
		{
			name:   "Unicode Characters (Truncates properly without breaking bytes)",
			input:  "こんにちは世界", // "Hello World" in Japanese (7 runes)
			maxLen: 3,
			want:   "こんに...", // 3 runes + 3 dots
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafeTruncateString(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("SafeTruncateString(%q, %d) = %q (len %d), want %q (len %d)",
					tt.input, tt.maxLen, got, len([]rune(got)), tt.want, len([]rune(tt.want)))
			}
		})
	}
}

func TestSafePadText(t *testing.T) {
	tests := []struct {
		name       string
		text       string
		width      int
		alignRight bool
		want       string
	}{
		{
			name:       "Align Left (Pads on the right)",
			text:       "Rust",
			width:      10,
			alignRight: false,
			want:       "Rust      ",
		},
		{
			name:       "Align Right (Pads on the left)",
			text:       "12h 45m",
			width:      10,
			alignRight: true,
			want:       "   12h 45m",
		},
		{
			name:       "Exact Width Align Left",
			text:       "Python",
			width:      6,
			alignRight: false,
			want:       "Python",
		},
		{
			name:       "Exact Width Align Right",
			text:       "Golang",
			width:      6,
			alignRight: true,
			want:       "Golang",
		},
		{
			name:       "Width Smaller Than Text (Does not truncate)",
			text:       "TakaTime",
			width:      4,
			alignRight: false,
			want:       "TakaTime",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SafePadText(tt.text, tt.width, tt.alignRight)
			if got != tt.want {
				t.Errorf("SafePadText(%q, %d, %v) = %q, want %q",
					tt.text, tt.width, tt.alignRight, got, tt.want)
			}
		})
	}
}
