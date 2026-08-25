package buildimg

import (
	"os"
	"testing"
	"time"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/fogleman/gg"
)

// --- Helper to get a valid font ---
// Image rendering requires a valid TrueType font.
// This helper tries to load a font from your project. If it can't find one, it skips the test.
func getValidFontBytes(t *testing.T) []byte {
	t.Helper()

	// A list of paths to try, accounting for different go test execution contexts
	possiblePaths := []string{
		"../../cm/report/FiraCodeNerdFontPropo-Retina.ttf", // If run from inside internal/BuildImg
		"../cm/report/FiraCodeNerdFontPropo-Retina.ttf",    // If run from inside internal/
		"cm/report/FiraCodeNerdFontPropo-Retina.ttf",       // If run from the project root
	}

	for _, fontPath := range possiblePaths {
		bytes, err := os.ReadFile(fontPath)
		if err == nil {
			return bytes // Found it! Return the bytes immediately.
		}
	}
	// If it fails all of them, it will print this so you know exactly why it skipped
	t.Skipf("Skipping image tests: Could not find FiraCodeNerdFontPropo-Retina.ttf in any expected location.")
	return nil
}

func TestDrawHeader_InvalidFont(t *testing.T) {
	dc := gg.NewContext(100, 100)
	theme := types.DefaultTheme()
	badFont := []byte("this is not a valid font file")

	err := DrawHeader(dc, "Test Header", badFont, theme, 65.0)
	if err == nil {
		t.Errorf("DrawHeader() expected an error for invalid font data, got nil")
	}

	// Test the unexported drawHeader too
	err = drawHeader(dc, "Test Header", badFont, theme, 40.0)
	if err == nil {
		t.Errorf("drawHeader() expected an error for invalid font data, got nil")
	}
}

func TestDrawListCard(t *testing.T) {
	theme := types.DefaultTheme()
	stats := []types.ListStats{
		{Label: "Go", Value: "10h 0m", Percent: 0.8, Color: "#00ADD8"},
		{Label: "LongProjectNameThatWillBeTruncated", Value: "2h 30m", Percent: 0.2, Color: "#FF0000"},
	}
	updatedAt := time.Now()

	t.Run("Invalid Font Fails Gracefully", func(t *testing.T) {
		badFont := []byte("bad font")
		img, err := DrawListCard("Top Languages", stats, badFont, updatedAt, theme, false)
		if err == nil || img != nil {
			t.Errorf("DrawListCard() expected error and nil image for bad font")
		}
	})

	t.Run("Successful Generation", func(t *testing.T) {
		validFont := getValidFontBytes(t)
		img, err := DrawListCard("Top Languages", stats, validFont, updatedAt, theme, false)
		if err != nil {
			t.Fatalf("DrawListCard() unexpected error: %v", err)
		}
		if img == nil {
			t.Errorf("DrawListCard() returned a nil image on success")
		}
	})
}

func TestDrawTimeCard(t *testing.T) {
	theme := types.DefaultTheme()
	data := types.TimeGridStruct{
		Yestarday: "5h 12m",
		Week:      "30h 0m",
		Month:     "120h 45m",
		AllTime:   "1500h 0m",
	}
	updatedAt := time.Now()

	t.Run("Invalid Font", func(t *testing.T) {
		_, err := DrawTimeCard(data, []byte("bad"), updatedAt, theme, "Tarun")
		if err == nil {
			t.Errorf("DrawTimeCard() expected error for bad font")
		}
	})

	t.Run("Successful Generation With Owner", func(t *testing.T) {
		validFont := getValidFontBytes(t)
		img, err := DrawTimeCard(data, validFont, updatedAt, theme, "Tarun")
		if err != nil {
			t.Fatalf("DrawTimeCard() unexpected error: %v", err)
		}
		if img == nil {
			t.Errorf("DrawTimeCard() returned a nil image")
		}
	})

	t.Run("Successful Generation Without Owner", func(t *testing.T) {
		validFont := getValidFontBytes(t)
		img, err := DrawTimeCard(data, validFont, updatedAt, theme, "")
		if err != nil {
			t.Fatalf("DrawTimeCard() unexpected error: %v", err)
		}
		if img == nil {
			t.Errorf("DrawTimeCard() returned a nil image")
		}
	})
}

func TestDrawTechCard(t *testing.T) {
	theme := types.DefaultTheme()
	editors := []types.ListStats{
		{Label: "Neovim", Percent: 0.9, Color: "#57A143"},
		{Label: "Unknown", Percent: 0.1, Color: "#FFFFFF"}, // Should be filtered out
	}
	osSystems := []types.ListStats{
		{Label: "Linux", Percent: 1.0, Color: "#FCC624"},
	}
	updatedAt := time.Now()

	t.Run("Invalid Font", func(t *testing.T) {
		_, err := DrawTechCard(editors, osSystems, []byte("bad"), updatedAt, theme)
		if err == nil {
			t.Errorf("DrawTechCard() expected error for bad font")
		}
	})

	t.Run("Successful Generation & Filters Unknown", func(t *testing.T) {
		validFont := getValidFontBytes(t)
		img, err := DrawTechCard(editors, osSystems, validFont, updatedAt, theme)
		if err != nil {
			t.Fatalf("DrawTechCard() unexpected error: %v", err)
		}
		if img == nil {
			t.Errorf("DrawTechCard() returned a nil image")
		}
	})
}

func TestLanguageStatsImg(t *testing.T) {
	theme := types.DefaultTheme()
	stats := []types.LanguageStat{
		{Name: "Go", Time: "10h", Percent: 0.5, Color: "#00ADD8"},
	}

	t.Run("Invalid Font", func(t *testing.T) {
		_, err := LanguageStatsImg(stats, []byte("bad"), theme)
		if err == nil {
			t.Errorf("LanguageStatsImg() expected error for bad font")
		}
	})

	t.Run("Successful Generation", func(t *testing.T) {
		validFont := getValidFontBytes(t)
		img, err := LanguageStatsImg(stats, validFont, theme)
		if err != nil {
			t.Fatalf("LanguageStatsImg() unexpected error: %v", err)
		}
		if img == nil {
			t.Errorf("LanguageStatsImg() returned a nil image")
		}
	})
}

func TestHeatmapStatsImg(t *testing.T) {
	theme := types.DefaultTheme()
	history := map[string]float64{
		time.Now().Format("2006-01-02"):                   2.5, // Color 2
		time.Now().AddDate(0, 0, -1).Format("2006-01-02"): 6.0, // Color 4
	}

	t.Run("Invalid Font", func(t *testing.T) {
		_, err := HeatmapStatsImg(history, 800, []byte("bad"), theme, 1.0)
		if err == nil {
			t.Errorf("HeatmapStatsImg() expected error for bad font")
		}
	})

	t.Run("Successful Generation with Scaling", func(t *testing.T) {
		validFont := getValidFontBytes(t)

		// Test standard scale
		img, err := HeatmapStatsImg(history, 800, validFont, theme, 1.0)
		if err != nil {
			t.Fatalf("HeatmapStatsImg() unexpected error at scale 1.0: %v", err)
		}
		if img == nil {
			t.Errorf("HeatmapStatsImg() returned a nil image")
		}

		// Test High-Res scale
		imgHighRes, err := HeatmapStatsImg(history, 800, validFont, theme, 2.0)
		if err != nil {
			t.Fatalf("HeatmapStatsImg() unexpected error at scale 2.0: %v", err)
		}

		// Verify canvas actually scaled up
		if imgHighRes.Bounds().Dx() <= img.Bounds().Dx() {
			t.Errorf("Expected scaled image width to be greater than base image width")
		}
	})
}

func TestSetupContext(t *testing.T) {
	// Tests the unexported setupContext
	theme := types.DefaultTheme()
	dc := setupContext(500, 500, theme)

	if dc.Width() != 500 || dc.Height() != 500 {
		t.Errorf("setupContext() dimensions incorrect. Got %dx%d, want 500x500", dc.Width(), dc.Height())
	}
}
