package utils

import (
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var titleCaser = cases.Title(language.English)

func ToTitleCase(s string) string {
	return titleCaser.String(strings.ToLower(s))
}
