package dbqueryv2

import (
	"strings"

	utils "github.com/Rtarun3606k/TakaTime/internal/Utils"
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/go-enry/go-enry/v2"
)

// CleanTelemetryLanguage standardizes incoming chaotic IDE strings into clean categories.
func CleanTelemetryLanguage(rawInput string) string {
	cleaned := strings.ToLower(strings.TrimSpace(rawInput))
	// check absolut map
	if exactMatch, exists := types.ExtensionOverrides[cleaned]; exists {
		return exactMatch
	}

	if lang, _ := enry.GetLanguageByExtension("dummy." + cleaned); lang != "" {

		// check enry  guess against the map to prevent hallucinations!
		enryCleaned := strings.ToLower(lang)
		if perfectMatch, exists := types.ExtensionOverrides[enryCleaned]; exists {
			return perfectMatch
		}

		// Return enry's perfectly cased guess
		return lang
	}

	// If even enry doesn't know, just format the string nicely
	return utils.ToTitleCase(cleaned)
}
