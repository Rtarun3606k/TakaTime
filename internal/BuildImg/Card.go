package buildimg

import (
	"fmt"
	"image"
	"log"
	"strings"
	"time"

	gogist "github.com/Rtarun3606k/TakaTime/internal/GoGist"
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/fogleman/gg"
)

func DrawHeader(dc *gg.Context, text string, fontData []byte, theme types.ThemeConfig, size float64) error {
	// Load font with the specific size requested (e.g., 65.0 or 30.0)
	headerFace, err := loadFontFace(fontData, size)
	if err != nil {
		log.Printf("error: failed to load header font: %v", err)
		return err
	}
	dc.SetFontFace(headerFace)
	dc.SetHexColor(theme.TextColor)

	// Calculate center position
	textWidth, _ := dc.MeasureString(text)
	x := (float64(dc.Width()) - textWidth) / 2
	y := size + 40 // Push down by size + padding so it never clips

	dc.DrawString(text, x, y)
	return nil
}

// 1. GENERIC LIST CARD (Languages / Projects)
func DrawListCard(title string, stats []types.ListStats, fontData []byte, updatedAt time.Time, theme types.ThemeConfig, isProject bool) (image.Image, error) {
	W, H := 1300, 600
	dc := gogist.SetupContext(W, H, theme)

	// Header
	if err := DrawHeader(dc, title, fontData, theme, 60); err != nil {
		log.Printf("error: failed to draw header: %v", err)
		return nil, err
	}

	listFace, _ := loadFontFace(fontData, 48)
	dc.SetFontFace(listFace)

	startY := 200.0
	gap := 85.0

	// Bar Config (Kept shorter as requested)
	barX := 640.0
	barW := 450.0
	barH := 30.0

	// Helper to truncate text > 8 chars with ".."
	truncate := func(s string) string {
		runes := []rune(s) // Use runes to handle emojis/special chars correctly
		if len(runes) > 8 {
			return string(runes[:8]) + ".."
		}
		return s
	}

	for i, stat := range stats {
		y := startY + (float64(i) * gap)

		// A. Label (Truncated)
		dc.SetHexColor(theme.TextColor)
		// We truncate here so "taka-time.nvim" becomes "taka-tim.."
		dc.DrawString(truncate(stat.Label), 60, y)

		// B. Value
		dc.SetHexColor(theme.SubTextColor)
		if isProject {

			dc.DrawString(stat.Value, 400, y)
		} else {

			dc.DrawString(stat.Value, 320, y)
		}

		// C. Bar BG
		dc.SetHexColor(theme.BarBackgroundColor)
		dc.DrawRoundedRectangle(barX, y-20, barW, barH, 10)
		dc.Fill()

		// D. Bar Fill
		dc.SetHexColor(stat.Color)
		fillW := barW * stat.Percent
		if fillW < 20 {
			fillW = 20
		}
		dc.DrawRoundedRectangle(barX, y-20, fillW, barH, 10)
		dc.Fill()

		// E. Percent Text
		dc.SetHexColor(theme.TextColor)
		dc.DrawString(fmt.Sprintf("%.1f%%", stat.Percent*100), barX+barW+30, y)
	}

	gogist.DrawFooter(dc, fontData, theme, updatedAt)
	return dc.Image(), nil
}

// 2. TIME GRID CARD (2x2)
func DrawTimeCard(data types.TimeGridStruct, fontData []byte, updatedAt time.Time, theme types.ThemeConfig, ownerName string) (image.Image, error) {
	W, H := 1200, 230
	dc := gogist.SetupContext(W, H, theme)

	// FIXED: Manually Draw Title at Top-Left (20, 45) instead of Center
	// This prevents it from crashing into the columns
	headerFace, _ := loadFontFace(fontData, 35)
	dc.SetFontFace(headerFace)
	dc.SetHexColor(theme.TextColor)
	title := ""
	if len(ownerName) < 1 {

		title = "Coding Activity"
	} else {

		title = ownerName + "'s Coding Activity"
	}
	// dc.DrawString("Coding Activity", 20, 45)
	textWidth, _ := dc.MeasureString(title)
	x := (float64(W) - textWidth) / 2
	y := 60.0 // Keep it high up

	dc.DrawString(title, x, y)
	// Fonts
	labelFace, _ := loadFontFace(fontData, 20)
	valFace, _ := loadFontFace(fontData, 42)

	drawColumn := func(x float64, label, val, colorHex string) {
		// Label
		dc.SetFontFace(labelFace)
		dc.SetHexColor(theme.SubTextColor)
		dc.DrawString(label, x, 130)

		// Value
		dc.SetFontFace(valFace)
		dc.SetHexColor(colorHex)
		dc.DrawString(val, x, 180)
	}

	colW := 300.0
	// Shift content slightly right to balance with the left-aligned title
	padding := 30.0

	drawColumn(0*colW+padding, "Yesterday", data.Yestarday, theme.Color3)
	drawColumn(1*colW+padding, "Last 7 Days", data.Week, theme.Color2)
	drawColumn(2*colW+padding, "Last 30 Days", data.Month, theme.Color1)
	drawColumn(3*colW+padding, "All Time", data.AllTime, theme.Color4)

	gogist.DrawFooter(dc, fontData, theme, updatedAt)
	return dc.Image(), nil
}

// 3. TECH STACK CARD -> BIGGER FONTS (High Visibility)
func DrawTechCard(editors []types.ListStats, osSystems []types.ListStats, fontData []byte, updatedAt time.Time, theme types.ThemeConfig) (image.Image, error) {
	W, H := 1600, 500
	dc := gogist.SetupContext(W, H, theme)

	if err := DrawHeader(dc, "Environment Stats", fontData, theme, 42); err != nil {
		log.Printf("error: failed to draw header: %v", err)
		return nil, err
	}

	// Filter Unknown
	filterUnknown := func(items []types.ListStats) []types.ListStats {
		var clean []types.ListStats
		for _, item := range items {
			if strings.ToLower(item.Label) != "unknown" && item.Label != "" {
				clean = append(clean, item)
			}
		}
		return clean
	}
	editors = filterUnknown(editors)
	osSystems = filterUnknown(osSystems)

	// Sub-headers
	subHeaderFace, _ := loadFontFace(fontData, 36)
	dc.SetFontFace(subHeaderFace)

	// Shifted Right Column Start: 650 -> 700 to create a safety gap
	dc.SetHexColor(theme.Color1)
	dc.DrawString("Editors", 60, 190)

	dc.SetHexColor(theme.Color2)
	dc.DrawString("Operating Systems", 900, 190)

	drawMini := func(list []types.ListStats, xOffset float64) {
		smallFace, _ := loadFontFace(fontData, 30)
		dc.SetFontFace(smallFace)

		startY := 260.0
		rowH := 75.0

		for i, item := range list {
			y := startY + float64(i)*rowH

			// Label
			dc.SetHexColor(theme.TextColor)
			dc.DrawString(item.Label, xOffset, y)

			// FIXED MATH FOR NO OVERLAP:
			// Left Col (Starts 50): BarX = 270. Width = 300. Ends = 570.
			// Right Col (Starts 700): Safe gap of 130px.

			barW := 300.0
			barH := 24.0
			barStartX := xOffset + 220.0

			// BG
			dc.SetHexColor(theme.BarBackgroundColor)
			dc.DrawRoundedRectangle(barStartX, y-18, barW, barH, 8)
			dc.Fill()

			// Fill
			dc.SetHexColor(item.Color)
			fillW := barW * item.Percent
			if fillW < 20 {
				fillW = 20
			}
			dc.DrawRoundedRectangle(barStartX, y-18, fillW, barH, 8)
			dc.Fill()

			// Percent
			dc.SetHexColor(theme.TextColor)
			dc.DrawString(fmt.Sprintf("%.0f%%", item.Percent*100), barStartX+barW+25, y)
		}
	}

	// Draw Left (Editors)
	drawMini(editors, 50)

	// Draw Right (OS) - Start further right at 700
	drawMini(osSystems, 900)

	gogist.DrawFooter(dc, fontData, theme, updatedAt)
	return dc.Image(), nil
}
