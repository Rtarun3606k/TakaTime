package dbqueryv2

import "testing"

func TestCleanTelemetryLanguage(t *testing.T) {
	tests := []struct {
		name     string
		rawInput string
		expected string
	}{
		// 1. Map Overrides (Layer 0)
		{"Empty string", "", "Plain Text"},
		{"Whitespace handling", "   ", "Plain Text"},
		{"VS Code React alias", "javascriptreact", "JavaScript React"},
		{"Messy casing override", "  tYpEsCrIpTrEaCt  ", "TypeScript React"},
		{"Config mapping", "dotenv", "Environment File"},

		// 2. Enry Hallucination Catches
		{"Enry YAML hallucination", "miniyaml", "YAML"},
		{"Enry Rust hallucination", "renderscript", "Rust"},
		{"Enry Solidity hallucination", "gerber image", "Solidity"},

		// 3. Enry Casing Preservation (Layer 1)
		{"Preserves C++ casing", "cpp", "C++"},
		{"Preserves Objective-C casing", "m", "Objective-C"},

		// 4. Ultimate Fallback (Layer 2)
		{"Unknown framework Title Cased", "qwik", "Qwik"},
		{"Unknown framework messy casing", "SVELTE", "Svelte"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := CleanTelemetryLanguage(tc.rawInput)
			if got != tc.expected {
				t.Errorf("CleanTelemetryLanguage(%q) = %q; want %q", tc.rawInput, got, tc.expected)
			}
		})
	}
}
