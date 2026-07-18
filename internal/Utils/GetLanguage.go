package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/go-enry/go-enry/v2"
)

// DetectLanguage Smart Three-Layer Architecture
func DetectLanguage(filename string) string {
	baseName := strings.ToLower(filepath.Base(filename))
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(baseName), "."))

	// ---------------------------------------------------------
	// LAYER 0: The Absolute Override
	// ---------------------------------------------------------
	if lang, exists := types.ExtensionOverrides[baseName]; exists {
		return lang
	}
	if lang, exists := types.ExtensionOverrides[ext]; exists {
		return lang
	}

	// ---------------------------------------------------------
	// LAYER 1: Enry Extension Guess (Zero Disk I/O)
	// ---------------------------------------------------------
	lang, safe := enry.GetLanguageByExtension(filename)

	if safe && lang != "" {
		return formatOutput(lang, ext)
	}

	// ---------------------------------------------------------
	// LAYER 2: Read 4KB for Deep Analysis
	// ---------------------------------------------------------
	filePtr, err := os.Open(filename)
	if err == nil {
		defer filePtr.Close()
		buffer := make([]byte, 4096)
		n, _ := filePtr.Read(buffer)
		lang = enry.GetLanguage(filename, buffer[:n])
	} else {
		lang = enry.GetLanguage(filename, nil)
	}

	// Pass both the Enry result AND the raw extension
	return formatOutput(lang, ext)
}

// formatOutput catches hallucinations and dynamically capitalizes unknowns
func formatOutput(rawLang string, fallbackExt string) string {
	if rawLang != "" {
		// Catch Enry hallucinations (e.g., Ecmarkup -> HTML)
		cleaned := strings.ToLower(strings.TrimSpace(rawLang))
		if perfectMatch, exists := types.ExtensionOverrides[cleaned]; exists {
			return perfectMatch
		}

		// Return rawLang directly to keep enry's exact casing (e.g., "Objective-C")
		return rawLang
	}

	// 2. Enry completely failed. Use the file extension as the language!
	if fallbackExt != "" {
		// e.g., ".svelte" -> "Svelte"
		return ToTitleCase(fallbackExt)
	}

	// 3. Ultimate Fallback (No extension, Enry failed)
	return "Plain Text"
}

// toTitleCase replaces the deprecated strings.Title() function
// It safely capitalizes the first letter of any string (e.g., "svelte" -> "Svelte")
func ToTitleCase(s string) string {
	if len(s) == 0 {
		return ""
	}
	return strings.ToUpper(string(s[0])) + s[1:]
}
