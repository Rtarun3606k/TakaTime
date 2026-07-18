package utils

import "testing"

func TestDetectLanguage(t *testing.T) {

	tests := []struct {
		name     string
		filename string
		expected string
	}{
		// 1. Map Overrides (Layer 0)
		{"Full exact filename", "/home/user/project/.env", "Environment File"},
		{"Full exact filename Go", "go.mod", "Go"},
		{"Extension exact override (TSX)", "src/app/page.tsx", "TypeScript React"},
		{"Extension exact override (Rust)", "main.rs", "Rust"},

		// 2. Enry Hallucination Catches
		{"Catches HTML hallucination", "index.html", "HTML"},
		{"Catches YAML hallucination", "config.yml", "YAML"},

		// 3. Enry Casing Preservation (Layer 1)
		{"Preserves C++ casing", "src/main.cpp", "C++"},
		{"Preserves Objective-C casing", "app/main.m", "Objective-C"},

		// 4. Ultimate Fallbacks
		{"Unknown extension Title Cased", "components/button.svelte", "Svelte"},
		{"No extension at all", "unknownfile", "Plain Text"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := DetectLanguage(tc.filename)
			if got != tc.expected {
				t.Errorf("DetectLanguage(%q) = %q; want %q", tc.filename, got, tc.expected)
			}
		})
	}
}
