package utils

import (
	"os"

	"github.com/go-enry/go-enry/v2"
)

// DetectLanguage Smart Two-Layer Architecture
func DetectLanguage(filename string) string {
	// ---------------------------------------------------------
	// LAYER 1: Zero Disk I/O
	// ---------------------------------------------------------
	lang, safe := enry.GetLanguageByExtension(filename)

	if safe && lang != "" {
		// log.Printf("[Fast Path] 0 bytes read for %s", filename) // Uncomment for debugging
		return lang
	}

	// ---------------------------------------------------------
	// LAYER 2: Read 4KB
	// ---------------------------------------------------------
	filePtr, err := os.Open(filename)
	if err == nil {
		defer filePtr.Close()

		buffer := make([]byte, 4096) // 4KB chunk limit
		n, _ := filePtr.Read(buffer)

		// Pass the 4KB chunk into enry's deep analyzer
		lang = enry.GetLanguage(filename, buffer[:n])

	} else {
		lang = enry.GetLanguage(filename, nil)
	}

	// ---------------------------------------------------------
	// FALLBACK
	// ---------------------------------------------------------
	if lang == "" {
		return "text"
	}
	return lang
}
