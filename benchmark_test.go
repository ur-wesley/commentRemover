package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func BenchmarkRemoveSingleLineComment(b *testing.B) {
	lang := Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	testLines := []string{
		"// This is a standalone comment",
		`fmt.Println("Hello World") // This is an inline comment`,
		`fmt.Println("No comment here")`,
		`fmt.Println("String with // comment inside")`,
		"var x = 42 // Another inline comment",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, line := range testLines {
			RemoveSingleLineComment(line, lang, false, false, false)
		}
	}
}

func BenchmarkUpdateMultiLineCommentState(b *testing.B) {
	lang := Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	testLines := []string{
		"/* Start of comment",
		"Inside comment",
		"Still inside",
		"End of comment */",
		"Regular code line",
		"/* Complete comment */",
	}

	b.ResetTimer()
	state := false
	for i := 0; i < b.N; i++ {
		for _, line := range testLines {
			state = UpdateMultiLineCommentState(line, lang, state)
		}
	}
}

func BenchmarkIsInsideStringLiteral(b *testing.B) {
	testLine := `fmt.Println("This is a string with // comment inside") // Real comment`
	positions := []int{10, 25, 35, 45, 55}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, pos := range positions {
			IsInsideStringLiteral(testLine, pos)
		}
	}
}

func BenchmarkProcessFile(b *testing.B) {
	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello World")
	
	/* This is a multi-line comment
	   // This should NOT be removed
	   End of multi-line comment */
	   
	fmt.Println("String with // comment inside")
	
	var x = 42
	
	// Multiple
	// Sequential
	// Comments
	
	fmt.Printf("Value: %d\n", x)
}

`

	tmpFile, err := os.CreateTemp("", "benchmark_*.go")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		b.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	lang := Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ProcessFile(tmpFile.Name(), lang, false, false, []string{})
		if err != nil {
			b.Fatalf("ProcessFile failed: %v", err)
		}
	}
}

func BenchmarkGetLanguageByExtension(b *testing.B) {
	filenames := []string{
		"main.go",
		"script.js",
		"component.tsx",
		"query.sql",
		"config.json",
		"README.md",
		"style.css",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, filename := range filenames {
			GetLanguageByExtension(filename)
		}
	}
}

func BenchmarkLargeFileProcessing(b *testing.B) {
	var lines []string
	for i := 0; i < 1000; i++ {
		lines = append(lines,
			"package main",
			"",
			"import \"fmt\"",
			"",
			"// This is a comment to remove",
			"func test() {",
			"    fmt.Println(\"Hello\") // Inline comment",
			"    var x = 42",
			"    /* Multi-line comment",
			"       // Should not be removed",
			"       End comment */",
			"}",
		)
	}
	content := strings.Join(lines, "\n")

	tmpFile, err := os.CreateTemp("", "large_benchmark_*.go")
	if err != nil {
		b.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		b.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	lang := Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ProcessFile(tmpFile.Name(), lang, false, false, []string{})
		if err != nil {
			b.Fatalf("ProcessFile failed: %v", err)
		}
	}
}

func BenchmarkDiscoverFilesRecursive(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "benchmark_discover_*")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	for i := 0; i < 5; i++ {
		for j := 0; j < 5; j++ {
			subDir := fmt.Sprintf("%s/dir%d/subdir%d", tempDir, i, j)
			if err := os.MkdirAll(subDir, 0755); err != nil {
				b.Fatalf("Failed to create subdir: %v", err)
			}

			testFile := fmt.Sprintf("%s/test%d.go", subDir, j)
			content := "package main\n// Test comment\nfunc main() {}\n"
			if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
				b.Fatalf("Failed to create test file: %v", err)
			}
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DiscoverFiles(tempDir, true, []string{})
		if err != nil {
			b.Fatalf("DiscoverFiles failed: %v", err)
		}
	}
}
