package Styles

import (
	"strings"
	"testing"

	"github.com/Rtarun3606k/TakaTime/internal/types"
)

// mockTheme provides a standard set of hex codes for testing
var mockTheme = types.ThemeConfig{
	BackgroundColor:    "#000000",
	TextColor:          "#FFFFFF",
	SubTextColor:       "#AAAAAA",
	BarBackgroundColor: "#222222",
	Color1:             "#FF0000",
	Color2:             "#00FF00",
	Color3:             "#0000FF",
	Color4:             "#FFFF00",
}

func TestInitStyles(t *testing.T) {
	styles := InitStyles(mockTheme)

	tests := []struct {
		name      string
		styleFunc func() string
	}{
		{"Title", func() string { return styles.Title.Render("Test") }},
		{"Text", func() string { return styles.Text.Render("Test") }},
		{"SubText", func() string { return styles.SubText.Render("Test") }},
		{"Box", func() string { return styles.Box.Render("Test") }},
		{"ListLabel", func() string { return styles.ListLabel.Render("Test") }},
		{"ListValue", func() string { return styles.ListValue.Render("Test") }},
		{"ListPercent", func() string { return styles.ListPercent.Render("Test") }},
		{"Color1", func() string { return styles.Color1.Render("Test") }},
		{"Navbar", func() string { return styles.Navbar.Render("Test") }},
		{"Footer", func() string { return styles.Footer.Render("Test") }},
		{"StatCard", func() string { return styles.StatCard.Render("Test") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.styleFunc()
			// Ensure it didn't return an empty string
			if len(got) == 0 {
				t.Errorf("InitStyles() %s generated an empty string", tt.name)
			}
			// Ensure the original text is still present inside the styled string
			if !strings.Contains(got, "Test") {
				t.Errorf("InitStyles() %s lost the original text during rendering", tt.name)
			}
		})
	}
}

func TestBuildStyles(t *testing.T) {
	styles := BuildStyles(mockTheme)

	tests := []struct {
		name      string
		styleFunc func() string
	}{
		{"Color3", func() string { return styles.Color3.Render("Test") }},
		{"Color4", func() string { return styles.Color4.Render("Test") }},
		{"Title", func() string { return styles.Title.Render("Test") }},
		{"Navbar", func() string { return styles.Navbar.Render("Test") }},
		{"Box", func() string { return styles.Box.Render("Test") }},
		{"StatCardValue", func() string { return styles.StatCardValue.Render("Test") }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.styleFunc()
			if len(got) == 0 {
				t.Errorf("BuildStyles() %s generated an empty string", tt.name)
			}
			if !strings.Contains(got, "Test") {
				t.Errorf("BuildStyles() %s lost the original text during rendering", tt.name)
			}
		})
	}
}
