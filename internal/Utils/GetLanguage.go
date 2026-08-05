package utils

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Rtarun3606k/TakaTime/internal/types"
)

// // DetectLanguage Smart Three-Layer Architecture
// func DetectLanguage(filename string) string {
// 	baseName := strings.ToLower(filepath.Base(filename))
// 	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(baseName), "."))
//
// 	// ---------------------------------------------------------
// 	// LAYER 0: The Absolute Override
// 	// ---------------------------------------------------------
// 	if lang, exists := types.ExtensionOverrides[baseName]; exists {
// 		return lang
// 	}
// 	if lang, exists := types.ExtensionOverrides[ext]; exists {
// 		return lang
// 	}
//
// 	// ---------------------------------------------------------
// 	// LAYER 1: Enry Extension Guess (Zero Disk I/O)
// 	// ---------------------------------------------------------
// 	lang, safe := enry.GetLanguageByExtension(filename)
//
// 	if safe && lang != "" {
// 		return formatOutput(lang, ext)
// 	}
//
// 	// ---------------------------------------------------------
// 	// LAYER 2: Read 4KB for Deep Analysis
// 	// ---------------------------------------------------------
// 	filePtr, err := os.Open(filename)
// 	if err == nil {
// 		defer filePtr.Close()
// 		buffer := make([]byte, 4096)
// 		n, _ := filePtr.Read(buffer)
// 		lang = enry.GetLanguage(filename, buffer[:n])
// 	} else {
// 		lang = enry.GetLanguage(filename, nil)
// 	}
//
// 	// Pass both the Enry result AND the raw extension
// 	return formatOutput(lang, ext)
// }
//
// // formatOutput catches hallucinations and dynamically capitalizes unknowns
// func formatOutput(rawLang string, fallbackExt string) string {
// 	if rawLang != "" {
// 		// Catch Enry hallucinations (e.g., Ecmarkup -> HTML)
// 		cleaned := strings.ToLower(strings.TrimSpace(rawLang))
// 		if perfectMatch, exists := types.ExtensionOverrides[cleaned]; exists {
// 			return perfectMatch
// 		}
//
// 		// Return rawLang directly to keep enry's exact casing (e.g., "Objective-C")
// 		return rawLang
// 	}
//
// 	// 2. Enry completely failed. Use the file extension as the language!
// 	if fallbackExt != "" {
// 		// e.g., ".svelte" -> "Svelte"
// 		return ToTitleCase(fallbackExt)
// 	}
//
// 	// 3. Ultimate Fallback (No extension, Enry failed)
// 	return "Plain Text"
// }
//
// // toTitleCase replaces the deprecated strings.Title() function
// // It safely capitalizes the first letter of any string (e.g., "svelte" -> "Svelte")
// func ToTitleCase(s string) string {
// 	if len(s) == 0 {
// 		return ""
// 	}
// 	return strings.ToUpper(string(s[0])) + s[1:]
// }

func DetectFromHeuristics(ext string, content []byte) ([]string, bool) {

	log.Printf("[Heuristic] Entered DetectFromHeuristics")
	rules, ok := types.HeuristicMap[ext]
	if !ok {
		return nil, false
	}

	candidates := make([]string, 0, len(rules))
	for _, rule := range rules {
		candidates = append(candidates, rule.Language)
	}

	bestScore := -1.0
	bestLanguage := ""

	for _, rule := range rules {
		matches := 0
		total := len(rule.Patterns)

		//positive patterns
		for _, pattern := range rule.Patterns {
			ok := pattern.Match(content)

			log.Printf("[Heuristic] %s: pattern=%q match=%v",
				rule.Language, pattern, ok)

			if ok {
				matches++
			}
		}

		//negative patterns
		negativeMatched := false
		for _, pattern := range rule.Negative {

			if pattern.Match(content) {
				log.Printf("[Heuristic] %s: negative pattern=%q matched",
					rule.Language, pattern)
				negativeMatched = true
				break
			}
		}

		if negativeMatched {
			continue
		}

		var score float64
		if total > 0 {
			score = float64(matches) / float64(total)
		}

		log.Printf("[Heuristic] %s score=%.2f (%d/%d)",
			rule.Language, score, matches, total)

		if score > bestScore {
			bestScore = score
			bestLanguage = rule.Language
		}

		//return if the score is 100% or 1.0
		if score == 1.0 {
			result := []string{rule.Language}

			for _, lang := range candidates {
				if lang != rule.Language {
					result = append(result, lang)
				}
			}

			return result, true
		}
	}

	if bestScore > 0 {
		result := []string{bestLanguage}

		for _, lang := range candidates {
			if lang != bestLanguage {
				result = append(result, lang)
			}
		}

		log.Printf("[Heuristic] Selected language: %s (score=%.2f)",
			bestLanguage, bestScore)

		return result, true
	}

	log.Printf("[Heuristic] No rule matched; returning candidates: %v", candidates)
	return candidates, false
}

func DetectLanguage(path string) (bool, string, []string, error) {
	filename := filepath.Base(path)

	log.Printf("[Detect] File: %s", filename)

	// Layer 1: special filenames.
	if lang, ok := types.ExtensionMap[filename]; ok {
		log.Printf("[Detect] Matched filename -> %v", lang)
		return true, lang, nil, nil
	}

	// Layer 2: extension lookup.
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")

	// Unique extension.
	if lang, ok := types.ExtensionMap[ext]; ok {
		log.Printf("[Detect] Unique extension match -> %s", lang)
		return true, lang, nil, nil
	}

	// Ambiguous extension.
	candidates, ok := types.AmbiguousExtensionMap[ext]
	if !ok {
		log.Printf("[Detect] Unknown extension %q", ext)
		return false, "", nil, nil
	}

	log.Printf("[Detect] Ambiguous extension -> %v", candidates)

	// No heuristics available.
	if _, ok := types.HeuristicMap[ext]; !ok {
		log.Printf("[Detect] No heuristics for extension %q", ext)
		return false, "", candidates, nil
	}

	// log.Printf("[Detect] Running heuristics for %q", ext)
	content, err := ReadFile(path)
	if err != nil {
		log.Printf("[Detect] ReadFile error: %v", err)
		return false, "", nil, err
	}

	langs, matched := DetectFromHeuristics(ext, content)

	if matched {
		return true, langs[0], candidates, nil
	}

	return false, "", candidates, nil
}

func ReadFile(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := make([]byte, 4096)

	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return nil, err
	}

	content := buf[:n]

	return content, nil
}
