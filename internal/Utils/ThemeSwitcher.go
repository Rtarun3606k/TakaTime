package utils

import "github.com/Rtarun3606k/TakaTime/internal/types"

func ThemeSwitcher(themeFlag string) types.ThemeConfig {

	var theme types.ThemeConfig

	theme = types.DefaultTheme()
	switch themeFlag {
	case "light":
		theme = types.LightTheme
	case "dracula":
		theme = types.DraculaTheme
	case "nord":
		theme = types.NordTheme
	case "gruvbox":
		theme = types.GruvboxTheme
	case "monokai":
		theme = types.MonokaiTheme
	case "cyberpunk":
		theme = types.CyberpunkTheme
	case "tokyonight":
		theme = types.TokyoNightTheme
	case "everforest":
		theme = types.EverforestTheme
	case "iceberg":
		theme = types.IcebergTheme
	case "sunset":
		theme = types.SunsetTheme
	case "deepocean":
		theme = types.DeepOceanTheme
	case "midnight":
		theme = types.MidnightPurpleTheme
	case "catppuccin":
		theme = types.CatppuccinTheme
	case "solarized":
		theme = types.SolarizedDarkTheme
	case "onedark":
		theme = types.OneDarkProTheme
	case "material":
		theme = types.MaterialDarkTheme
	case "synthwave":
		theme = types.SynthwaveTheme
	case "rosepine":
		theme = types.RosepineTheme
	default:
		// Default to Dark if unknown or explicitly "dark"
		theme = types.DefaultTheme()
	}
	return theme
}
