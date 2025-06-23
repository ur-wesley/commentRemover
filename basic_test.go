package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestGetLanguageByExtension(t *testing.T) {
	tests := []struct {
		filename     string
		expectedLang string
		supported    bool
	}{
		{
			filename:     "main.go",
			expectedLang: "Go",
			supported:    true,
		},
		{
			filename:     "script.js",
			expectedLang: "TypeScript/JavaScript",
			supported:    true,
		},
		{
			filename:     "component.tsx",
			expectedLang: "TypeScript/JavaScript",
			supported:    true,
		},
		{
			filename:     "query.sql",
			expectedLang: "SQL",
			supported:    true,
		},
		{
			filename:     "config.json",
			expectedLang: "JSON",
			supported:    true,
		},
		{
			filename:     "README.md",
			expectedLang: "",
			supported:    false,
		},
		{
			filename:     "no_extension",
			expectedLang: "",
			supported:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			lang, supported := GetLanguageByExtension(tt.filename)

			if supported != tt.supported {
				t.Errorf("Expected supported=%v, got %v", tt.supported, supported)
			}

			if tt.supported {
				if lang == nil {
					t.Fatal("Expected language object, got nil")
				}
				if lang.Name != tt.expectedLang {
					t.Errorf("Expected language %q, got %q", tt.expectedLang, lang.Name)
				}
			} else {
				if lang != nil {
					t.Errorf("Expected nil language for unsupported file, got %v", lang)
				}
			}
		})
	}
}

func TestRemoveSingleLineComment(t *testing.T) {
	lang := Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	tests := []struct {
		name               string
		input              string
		inMultiLineComment bool
		expectedOutput     string
		expectedRemoved    bool
	}{
		{
			name:               "standalone comment",
			input:              "// This is a comment",
			inMultiLineComment: false,
			expectedOutput:     "REMOVE_LINE",
			expectedRemoved:    true,
		},
		{
			name:               "inline comment",
			input:              `fmt.Println("Hello") // This is a comment`,
			inMultiLineComment: false,
			expectedOutput:     `fmt.Println("Hello")`,
			expectedRemoved:    true,
		},
		{
			name:               "no comment",
			input:              `fmt.Println("Hello")`,
			inMultiLineComment: false,
			expectedOutput:     `fmt.Println("Hello")`,
			expectedRemoved:    false,
		},
		{
			name:               "comment in string literal",
			input:              `fmt.Println("Hello // World")`,
			inMultiLineComment: false,
			expectedOutput:     `fmt.Println("Hello // World")`,
			expectedRemoved:    false,
		},
		{
			name:               "inside multi-line comment",
			input:              "// This comment is inside a multi-line block",
			inMultiLineComment: true,
			expectedOutput:     "// This comment is inside a multi-line block",
			expectedRemoved:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, removed := RemoveSingleLineComment(tt.input, lang, tt.inMultiLineComment, false, false)
			if output != tt.expectedOutput {
				t.Errorf("Expected output %q, got %q", tt.expectedOutput, output)
			}
			if removed != tt.expectedRemoved {
				t.Errorf("Expected removed %v, got %v", tt.expectedRemoved, removed)
			}
		})
	}
}

func TestUpdateMultiLineCommentState(t *testing.T) {
	lang := Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	tests := []struct {
		name         string
		input        string
		currentState bool
		expected     bool
	}{
		{
			name:         "start multi-line comment",
			input:        "/* This starts a comment",
			currentState: false,
			expected:     true,
		},
		{
			name:         "end multi-line comment",
			input:        "This ends a comment */",
			currentState: true,
			expected:     false,
		},
		{
			name:         "complete multi-line comment",
			input:        "/* Complete comment */",
			currentState: false,
			expected:     false,
		},
		{
			name:         "no comment markers",
			input:        "Regular code line",
			currentState: false,
			expected:     false,
		},
		{
			name:         "continue in comment",
			input:        "Still in comment",
			currentState: true,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := UpdateMultiLineCommentState(tt.input, lang, tt.currentState)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for input %q with state %v", tt.expected, result, tt.input, tt.currentState)
			}
		})
	}
}

func TestIsInsideStringLiteral(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		pos      int
		expected bool
	}{
		{
			name:     "outside string",
			line:     `fmt.Println("Hello") // comment`,
			pos:      21,
			expected: false,
		},
		{
			name:     "inside double quotes",
			line:     `fmt.Println("Hello // World")`,
			pos:      19,
			expected: true,
		},
		{
			name:     "inside single quotes",
			line:     `char := '//' // comment`,
			pos:      10,
			expected: true,
		},
		{
			name:     "inside backticks",
			line:     "msg := `Hello // World` // comment",
			pos:      14,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsInsideStringLiteral(tt.line, tt.pos)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for line %q at position %d", tt.expected, result, tt.line, tt.pos)
			}
		})
	}
}

func TestProcessFile(t *testing.T) {
	content := `package main

import "fmt"

func main() {
	fmt.Println("Hello")
	/* This is a multi-line comment
	   // This should NOT be removed
	   End of multi-line comment */
	fmt.Println("String with // comment inside")
}
`

	tmpFile, err := os.CreateTemp("", "test_*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	lang := Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	result, err := ProcessFile(tmpFile.Name(), lang, false)
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	if result.OriginalLines != 13 {
		t.Errorf("Expected 13 original lines, got %d", result.OriginalLines)
	}

	if result.CommentsRemoved != 3 {
		t.Errorf("Expected 3 comments removed, got %d", result.CommentsRemoved)
	}

	expectedRemovedLines := []int{5, 7, 8}
	actualRemovedLines := make([]int, len(result.RemovedComments))
	for i, comment := range result.RemovedComments {
		actualRemovedLines[i] = comment.LineNumber
	}

	if !reflect.DeepEqual(expectedRemovedLines, actualRemovedLines) {
		t.Errorf("Expected removed lines %v, got %v", expectedRemovedLines, actualRemovedLines)
	}

	modifiedContent := strings.Join(result.ModifiedLines, "\n")
	if !strings.Contains(modifiedContent, "/* This is a multi-line comment") {
		t.Error("Multi-line comment should be preserved")
	}
	if !strings.Contains(modifiedContent, "// This should NOT be removed") {
		t.Error("Comment inside multi-line comment should be preserved")
	}
	if !strings.Contains(modifiedContent, `"String with // comment inside"`) {
		t.Error("Comment inside string should be preserved")
	}
}

func TestDiscoverFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test_discover_*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := map[string]string{
		"test.go":  "package main\n// comment\nfunc main() {}",
		"test.js":  "// comment\nconsole.log('hello');",
		"test.sql": "-- comment\nSELECT * FROM users;",
		"test.txt": "This should be ignored",
	}

	for filename, content := range testFiles {
		filePath := tempDir + "/" + filename
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	files, err := DiscoverFiles(tempDir, false)
	if err != nil {
		t.Fatalf("DiscoverFiles failed: %v", err)
	}

	expectedCount := 3 // .go, .js, .sql files
	if len(files) != expectedCount {
		t.Errorf("Expected %d files, got %d", expectedCount, len(files))
	}

	// Check that we got the right file types
	extensions := make(map[string]bool)
	for _, file := range files {
		for _, ext := range file.Language.Extensions {
			extensions[ext] = true
		}
	}

	expectedExtensions := []string{".go", ".js", ".sql"}
	for _, ext := range expectedExtensions {
		if !extensions[ext] {
			t.Errorf("Expected to find files with extension %s", ext)
		}
	}
}
