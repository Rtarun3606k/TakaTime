package buildimg

import (
	_ "embed"
	"fmt"
	"image"
	"log"
	"time"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// Updated to accept 'theme'
func LanguageStatsImg(stats []types.LanguageStat, fontData []byte, theme types.ThemeConfig) (image.Image, error) {
	currettime := time.Now()

	// 1. Setup Canvas (Use theme background)
	dc := setupContext(types.CanvasWidth, types.CanvasHeight, theme)

	// 2. Draw Header
	// FIXED: Now passing 'theme' and size '40.0'
	if err := drawHeader(dc, "TakaTime Language Stats", fontData, theme, 40.0); err != nil {
		log.Printf("error: failed to draw header: %v", err)
		return nil, err
	}

	// 3. Draw the List of Languages
	// FIXED: Now passing 'theme'
	if err := drawStatsList(dc, stats, fontData, theme); err != nil {
		log.Printf("error: failed to draw stats list: %v", err)
		return nil, err
	}

	// 4. Draw Footer
	footerText := fmt.Sprintf("Last Updated: %s", currettime.Format("2006-01-02 15:04:05"))
	// FIXED: Now passing 'theme'
	if err := drawFooter(dc, footerText, fontData, theme); err != nil {
		log.Printf("error: failed to draw footer: %v", err)
		return nil, err
	}

	return dc.Image(), nil
}

// --- Helper Functions ---

// Updated to accept theme
func setupContext(w, h int, theme types.ThemeConfig) *gg.Context {
	dc := gg.NewContext(w, h)
	dc.SetHexColor(theme.BackgroundColor) // Use theme!
	dc.Clear()
	return dc
}

func drawHeader(dc *gg.Context, text string, fontData []byte, theme types.ThemeConfig, size float64) error {
	headerFace, err := loadFontFace(fontData, size)
	if err != nil {
		log.Printf("error: failed to load header font: %v", err)
		return err
	}
	dc.SetFontFace(headerFace)
	dc.SetHexColor(theme.TextColor)

	textWidth, _ := dc.MeasureString(text)
	dc.DrawString(text, (float64(dc.Width())-textWidth)/2, size+40)
	return nil
}

// Updated to accept theme
func drawStatsList(dc *gg.Context, stats []types.LanguageStat, fontData []byte, theme types.ThemeConfig) error {
	listFace, err := loadFontFace(fontData, 20)
	if err != nil {
		log.Printf("error: failed to load list font: %v", err)
		return fmt.Errorf("failed to load list font: %w", err)
	}
	dc.SetFontFace(listFace)

	startY := 150.0
	gap := 60.0

	for i, stat := range stats {
		y := startY + (float64(i) * gap)
		drawRow(dc, stat, y, theme) // Pass theme down
	}
	return nil
}

// Updated to accept theme
func drawRow(dc *gg.Context, stat types.LanguageStat, y float64, theme types.ThemeConfig) {
	// A. Language Name
	dc.SetHexColor(theme.TextColor)
	dc.DrawString(stat.Name, 50, y)

	// B. Time Text
	dc.SetHexColor(theme.SubTextColor)
	dc.DrawString(stat.Time, 200, y)

	// C. Progress Bar Background
	barX := 350.0
	barW := 350.0
	barH := 15.0
	dc.SetHexColor(theme.BarBackgroundColor)
	dc.DrawRoundedRectangle(barX, y-10, barW, barH, 4)
	dc.Fill()

	// D. Progress Bar Fill
	dc.SetHexColor(stat.Color)
	fillWidth := barW * stat.Percent
	if fillWidth < 10 {
		fillWidth = 10
	}
	dc.DrawRoundedRectangle(barX, y-10, fillWidth, barH, 4)
	dc.Fill()

	// E. Percentage Text
	dc.SetHexColor(theme.TextColor)
	percentStr := fmt.Sprintf("%.1f%%", stat.Percent*100)
	dc.DrawString(percentStr, barX+barW+20, y)
}

func loadFontFace(fontBytes []byte, points float64) (font.Face, error) {
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Printf("error: failed to parse font: %v", err)
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size:    points,
		Hinting: font.HintingFull,
	})
	return face, nil
}

// Updated to accept theme
func drawFooter(dc *gg.Context, text string, fontData []byte, theme types.ThemeConfig) error {
	footerFace, err := loadFontFace(fontData, 12)
	if err != nil {
		log.Printf("error: failed to load footer font: %v", err)
		return err
	}
	dc.SetFontFace(footerFace)

	dc.SetHexColor(theme.SubTextColor) // Use theme subtext
	w, h := float64(dc.Width()), float64(dc.Height())
	dc.DrawStringAnchored(text, w-10, h-10, 1.0, 0.0)
	return nil
}

func HeatmapStatsImg(history map[string]float64, baseWidth int, fontData []byte, theme types.ThemeConfig, scale float64) (image.Image, error) {
	// 1. Image Layout Constants
	baseHeight := 210.0
	cellSize := 10.0 * scale
	cellSpacing := 3.0 * scale
	cols := 53.0 // 52 weeks + overlapping days

	// Apply scale to the overall canvas dimensions
	actualWidth := int(float64(baseWidth) * scale)
	actualHeight := int(baseHeight * scale)

	// Calculate grid width to center it
	gridWidth := cols * (cellSize + cellSpacing)
	startX := (float64(actualWidth) - gridWidth) / 2.0

	// Minimum padding for Y-Axis labels (scaled)
	if startX < (40.0 * scale) {
		startX = 40.0 * scale
	}

	startY := 75.0 * scale

	// 2. Initialize Context with High-Res Dimensions
	dc := gg.NewContext(actualWidth, actualHeight)

	// Load Font from byte slice
	f, err := truetype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %v", err)
	}

	// Scale the fonts
	normalFace := truetype.NewFace(f, &truetype.Options{Size: 12.0 * scale})
	titleFace := truetype.NewFace(f, &truetype.Options{Size: 16.0 * scale})

	// Background
	dc.SetHexColor(theme.BackgroundColor)
	dc.Clear()

	// 3. Draw Title
	dc.SetFontFace(titleFace)
	dc.SetHexColor(theme.TextColor)
	dc.DrawStringAnchored("━ 365-Day Contribution Graph ━", float64(actualWidth)/2.0, 30.0*scale, 0.5, 0.5)

	dc.SetFontFace(normalFace)

	// 4. Time Setup
	today := time.Now()
	start := today.AddDate(0, 0, -364)
	offset := int(start.Weekday())
	curr := start.AddDate(0, 0, -offset)

	var lastMonth time.Month = 0

	// 5. Draw the Grid & Header Labels
	col := 0
	for !curr.After(today) {
		// Draw Month Header
		if curr.Month() != lastMonth {
			if col > 0 {
				dc.SetHexColor(theme.SubTextColor)
				xPos := startX + float64(col)*(cellSize+cellSpacing)
				dc.DrawStringAnchored(curr.Format("Jan"), xPos, startY-(15.0*scale), 0, 0.5)
			}
			lastMonth = curr.Month()
		}

		// Draw 7 Days
		for row := 0; row < 7; row++ {
			if !curr.Before(start) && !curr.After(today) {
				dateStr := curr.Format("2006-01-02")
				val := history[dateStr]

				colorHex := theme.BarBackgroundColor
				// FIX: Added the 4-tier color threshold logic
				if val > 0 && val <= 1.0 {
					colorHex = theme.Color1
				} else if val > 1.0 && val <= 3.0 {
					colorHex = theme.Color2
				} else if val > 3.0 && val <= 5.0 {
					colorHex = theme.Color3
				} else if val > 5.0 {
					colorHex = theme.Color4 // The intense color for heavy coding days!
				}

				x := startX + float64(col)*(cellSize+cellSpacing)
				y := startY + float64(row)*(cellSize+cellSpacing)

				dc.SetHexColor(colorHex)
				dc.DrawRoundedRectangle(x, y, cellSize, cellSize, 2.0*scale)
				dc.Fill()
			}
			curr = curr.AddDate(0, 0, 1)
		}
		col++
	}

	// 6. Draw Y-Axis Labels (Mon, Wed, Fri)
	dc.SetHexColor(theme.SubTextColor)
	textOffsetX := startX - (12.0 * scale)

	dc.DrawStringAnchored("Mon", textOffsetX, startY+float64(1)*(cellSize+cellSpacing)+(cellSize/2.0)+cellSpacing, 1, 0.5)
	dc.DrawStringAnchored("Wed", textOffsetX, startY+float64(3)*(cellSize+cellSpacing)+(cellSize/2.0)+cellSpacing, 1, 0.5)
	dc.DrawStringAnchored("Fri", textOffsetX, startY+float64(5)*(cellSize+cellSpacing)+(cellSize/2.0)+cellSpacing, 1, 0.5)

	// 7. Draw Legend at the bottom right
	legendY := startY + float64(8)*(cellSize+cellSpacing)

	// FIX: Shifted the legend X-start from 150.0 to 165.0 to accommodate the extra 5th box
	legendX := startX + gridWidth - (165.0 * scale)

	dc.SetHexColor(theme.SubTextColor)
	dc.DrawStringAnchored("Less", legendX-(20.0*scale), legendY+(cellSize/2.0), 1, 0.5)

	// FIX: Added theme.Color4 to the slice
	colors := []string{theme.SubTextColor, theme.Color1, theme.Color2, theme.Color3, theme.Color4}
	for i, c := range colors {
		x := legendX + float64(i)*(cellSize+cellSpacing)
		dc.SetHexColor(c)
		dc.DrawRoundedRectangle(x, legendY, cellSize, cellSize, 2.0*scale)
		dc.Fill()
	}

	dc.SetHexColor(theme.SubTextColor)
	dc.DrawStringAnchored("More", legendX+float64(len(colors))*(cellSize+cellSpacing)+(5.0*scale), legendY+(cellSize/2.0), 0, 0.5)

	return dc.Image(), nil
}
