package main

import (
	"slices"
	"strings"
)

type Language struct {
	Name                        string
	Extensions                  []string
	SingleLineStart             string
	MultiLineStart              string
	MultiLineEnd                string
	AdditionalMultiLinePatterns []MultiLinePattern
}

type MultiLinePattern struct {
	Start string
	End   string
}

var SupportedLanguages = map[string]Language{
	"typescript": {
		Name:            "TypeScript/JavaScript",
		Extensions:      []string{".ts", ".tsx", ".js", ".jsx"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
		AdditionalMultiLinePatterns: []MultiLinePattern{
			{Start: "{/*", End: "*/}"},
		},
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
	"php": {
		Name:            "PHP",
		Extensions:      []string{".php", ".phtml"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	},
	"csharp": {
		Name:            "C#",
		Extensions:      []string{".cs"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
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
