package gogist

import (
	"os"
	"testing"
	"time"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/fogleman/gg"
)

// getValidFontBytes safely tries to load the font file from a few possible relative paths.
func getValidFontBytes(t *testing.T) []byte {
	t.Helper()

	possiblePaths := []string{
		"../../cm/report/FiraCodeNerdFontPropo-Retina.ttf", // If run from internal/GoGist
		"../cm/report/FiraCodeNerdFontPropo-Retina.ttf",    // If run from internal/
		"cm/report/FiraCodeNerdFontPropo-Retina.ttf",       // If run from project root
	}

	for _, fontPath := range possiblePaths {
		bytes, err := os.ReadFile(fontPath)
		if err == nil {
			return bytes
		}
	}

	t.Skipf("Skipping font-dependent tests: Could not find FiraCodeNerdFontPropo-Retina.ttf")
	return nil
}

func TestSetupContext(t *testing.T) {
	theme := types.ThemeConfig{
		BackgroundColor: "#111111",
	}

	w, h := 800, 600
	dc := SetupContext(w, h, theme)

	if dc == nil {
		t.Fatal("SetupContext() returned a nil context")
	}

	if dc.Width() != w {
		t.Errorf("SetupContext() width = %d, want %d", dc.Width(), w)
	}

	if dc.Height() != h {
		t.Errorf("SetupContext() height = %d, want %d", dc.Height(), h)
	}
}

func TestLoadFontFace(t *testing.T) {
	t.Run("Invalid Font Data", func(t *testing.T) {
		badData := []byte("this is not a valid font")
		face, err := LoadFontFace(badData, 12.0)
		if err == nil {
			t.Errorf("LoadFontFace() expected error for invalid font data, got nil")
		}
		if face != nil {
			t.Errorf("LoadFontFace() expected nil face for invalid font data, got %v", face)
		}
	})

	t.Run("Valid Font Data", func(t *testing.T) {
		validData := getValidFontBytes(t)
		face, err := LoadFontFace(validData, 12.0)
		if err != nil {
			t.Fatalf("LoadFontFace() unexpected error: %v", err)
		}
		if face == nil {
			t.Errorf("LoadFontFace() returned nil face on success")
		}
	})
}

func TestLoadFontFace_Unexported(t *testing.T) {
	// Testing the unexported version of the function
	t.Run("Invalid Font Data", func(t *testing.T) {
		badData := []byte("this is not a valid font")
		face, err := loadFontFace(badData, 12.0)
		if err == nil {
			t.Errorf("loadFontFace() expected error for invalid font data, got nil")
		}
		if face != nil {
			t.Errorf("loadFontFace() expected nil face for invalid font data, got %v", face)
		}
	})

	t.Run("Valid Font Data", func(t *testing.T) {
		validData := getValidFontBytes(t)
		face, err := loadFontFace(validData, 12.0)
		if err != nil {
			t.Fatalf("loadFontFace() unexpected error: %v", err)
		}
		if face == nil {
			t.Errorf("loadFontFace() returned nil face on success")
		}
	})
}

func TestDrawFooter(t *testing.T) {
	theme := types.ThemeConfig{
		SubTextColor: "#888888",
	}
	updatedAt := time.Now()
	dc := gg.NewContext(500, 500)

	t.Run("Invalid Font Data", func(t *testing.T) {
		badData := []byte("bad font bytes")
		err := DrawFooter(dc, badData, theme, updatedAt)
		if err == nil {
			t.Errorf("DrawFooter() expected an error for bad font data, got nil")
		}
	})

	t.Run("Valid Font Data", func(t *testing.T) {
		validData := getValidFontBytes(t)
		err := DrawFooter(dc, validData, theme, updatedAt)
		if err != nil {
			t.Fatalf("DrawFooter() unexpected error: %v", err)
		}
	})
}
