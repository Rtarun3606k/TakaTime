package utils

import (
	"reflect"
	"testing"

	"github.com/Rtarun3606k/TakaTime/internal/types"
)

func TestThemeSwitcher(t *testing.T) {
	tests := []struct {
		name      string
		themeFlag string
		want      types.ThemeConfig
	}{
		{
			name:      "Dracula Theme",
			themeFlag: "dracula",
			want:      types.DraculaTheme,
		},
		{
			name:      "Catppuccin Theme",
			themeFlag: "catppuccin",
			want:      types.CatppuccinTheme,
		},
		{
			name:      "Rose Pine Theme",
			themeFlag: "rosepine",
			want:      types.RosepineTheme,
		},
		{
			name:      "Tokyo Night Theme",
			themeFlag: "tokyonight",
			want:      types.TokyoNightTheme,
		},
		{
			name:      "Light Theme",
			themeFlag: "light",
			want:      types.LightTheme,
		},
		{
			name:      "Explicit Dark (Defaults)",
			themeFlag: "dark",
			want:      types.DefaultTheme(),
		},
		{
			name:      "Empty String Fallback",
			themeFlag: "",
			want:      types.DefaultTheme(),
		},
		{
			name:      "Unknown Theme Fallback",
			themeFlag: "non-existent-theme-123",
			want:      types.DefaultTheme(),
		},
		{
			name:      "Case Sensitivity Fallback (Expects all lowercase)",
			themeFlag: "Dracula", // Note the capital D
			want:      types.DefaultTheme(), // Fails the strict "dracula" switch case
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ThemeSwitcher(tt.themeFlag)
			
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ThemeSwitcher(%q) returned unexpected theme struct", tt.themeFlag)
			}
		})
	}
}
