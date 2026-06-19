package gogist

import (
	"fmt"
	"time"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func SetupContext(w, h int, theme types.ThemeConfig) *gg.Context {

	dc := gg.NewContext(w, h)

	dc.SetHexColor(theme.BackgroundColor)
	dc.Clear()
	return dc

}

func LoadFontFace(fontBytes []byte, points float64) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {

		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size:    points,
		Hinting: font.HintingFull,
	})
	return face, nil
}

func loadFontFace(fontBytes []byte, points float64) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size:    points,
		Hinting: font.HintingFull,
	})
	return face, nil
}

func DrawFooter(dc *gg.Context, fontData []byte, theme types.ThemeConfig, updatedAt time.Time) error {
	footerFace, err := loadFontFace(fontData, 12)
	if err != nil {
		return err
	}
	dc.SetFontFace(footerFace)
	dc.SetHexColor(theme.SubTextColor)

	text := fmt.Sprintf("Last Updated: %s", updatedAt.Format("2006-01-02 15:04:05"))
	w, h := float64(dc.Width()), float64(dc.Height())

	dc.DrawStringAnchored(text, w-10, h-10, 1.0, 0.0)
	return nil
}
