package utils

import (
	"strings"
	"testing"
)

func TestGenerateOutput(t *testing.T) {
	got := GenerateOutput()

	expectedElements := []string{
		"<h2 align=\"center\">TakaTime Weekly Report</h2>",
		"<img src=\"./public/taka-time.png\" width=\"100%\" alt=\"Time Stats\" />",
		"<img src=\"./public/taka-languages30.png\" width=\"400\" alt=\"Languages\" />",
		"<img src=\"./public/taka-projects30.png\" width=\"400\" alt=\"Projects\" />",
		"<img src=\"./public/taka-languages.png\" width=\"400\" alt=\"Languages\" />",
		"<img src=\"./public/taka-projects.png\" width=\"400\" alt=\"Projects\" />",
		"<img src=\"./public/taka-heatmap.png\" width=\"100%\" alt=\"Heatmap\" />",
		"<img src=\"./public/taka-tech.png\" width=\"100%\" alt=\"Tech Stack\" />",
		"Generated automatically by <a href=\"https://github.com/Rtarun3606k/TakaTime\">TakaTime</a>",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(got, expected) {
			t.Errorf("GenerateOutput() missing expected element:\n%s\n\nGot Output:\n%s", expected, got)
		}
	}

	if !strings.HasPrefix(got, "<h2") {
		t.Errorf("GenerateOutput() should start with an <h2> tag, got: %s", got[:10])
	}

	if !strings.HasSuffix(got, "TakaTime</a></em></p>") {
		t.Errorf("GenerateOutput() should end with the footer closing tags")
	}
}
