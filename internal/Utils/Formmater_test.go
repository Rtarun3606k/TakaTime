package utils

import (
	"testing"
)

func TestToTitleCase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "All Lowercase",
			input: "hello world",
			want:  "Hello World",
		},
		{
			name:  "All Uppercase",
			input: "HELLO WORLD",
			want:  "Hello World", // Should lowercase first, then capitalize
		},
		{
			name:  "Mixed Case",
			input: "gOlAnG pRoGrAmMiNg",
			want:  "Golang Programming",
		},
		{
			name:  "Single Word",
			input: "takatime",
			want:  "Takatime",
		},
		{
			name:  "With Punctuation",
			input: "hello, world!",
			want:  "Hello, World!",
		},
		{
			name:  "With Numbers",
			input: "123 testing string",
			want:  "123 Testing String",
		},
		{
			name:  "Empty String",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToTitleCase(tt.input)
			if got != tt.want {
				t.Errorf("ToTitleCase(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
