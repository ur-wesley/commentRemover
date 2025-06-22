package main

import (
	"slices"
	"strings"
)

type Language struct {
	Name            string
	Extensions      []string
	SingleLineStart string
	MultiLineStart  string
	MultiLineEnd    string
}

var SupportedLanguages = map[string]Language{
	"typescript": {
		Name:            "TypeScript/JavaScript",
		Extensions:      []string{".ts", ".tsx", ".js", ".jsx"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	},
	"go": {
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	},
	"sql": {
		Name:            "SQL",
		Extensions:      []string{".sql"},
		SingleLineStart: "--",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	},
	"json": {
		Name:            "JSON",
		Extensions:      []string{".json"},
		SingleLineStart: "//",
		MultiLineStart:  "",
		MultiLineEnd:    "",
	},
}

func GetLanguageByExtension(filename string) (*Language, bool) {
	ext := strings.ToLower(filename)
	dotIndex := strings.LastIndex(ext, ".")
	if dotIndex == -1 {
		return nil, false
	}

	extension := ext[dotIndex:]

	for _, lang := range SupportedLanguages {
		if slices.Contains(lang.Extensions, extension) {
			return &lang, true
		}
	}

	return nil, false
}
