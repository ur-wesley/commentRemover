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
			filename:     "script.php",
			expectedLang: "PHP",
			supported:    true,
		},
		{
			filename:     "template.phtml",
			expectedLang: "PHP",
			supported:    true,
		},
		{
			filename:     "Program.cs",
			expectedLang: "C#",
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

func TestJSXCommentDetection(t *testing.T) {
	lang := Language{
		Name:            "TypeScript/JavaScript",
		Extensions:      []string{".ts", ".tsx", ".js", ".jsx"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
		AdditionalMultiLinePatterns: []MultiLinePattern{
			{Start: "{/*", End: "*/}"},
		},
	}

	tests := []struct {
		name         string
		input        string
		currentState bool
		expected     bool
	}{
		{
			name:         "start JSX comment",
			input:        "  {/* This is a JSX comment",
			currentState: false,
			expected:     true,
		},
		{
			name:         "end JSX comment",
			input:        "This ends a JSX comment */}",
			currentState: true,
			expected:     false,
		},
		{
			name:         "complete JSX comment",
			input:        "  {/* Complete JSX comment */}",
			currentState: false,
			expected:     false,
		},
		{
			name:         "JSX comment with special chars",
			input:        "  {/* Comment with -- // *** */}",
			currentState: false,
			expected:     false,
		},
		{
			name:         "continue in JSX comment",
			input:        "Still in JSX comment",
			currentState: true,
			expected:     true,
		},
		{
			name:         "mixed JSX and regular comments",
			input:        "  {/* JSX comment */} /* Regular comment",
			currentState: false,
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

	result, err := ProcessFile(tmpFile.Name(), lang, false, false, []string{})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	if result.OriginalLines != 11 {
		t.Errorf("Expected 11 original lines, got %d", result.OriginalLines)
	}

	if result.CommentsRemoved != 0 {
		t.Errorf("Expected 0 comments removed, got %d", result.CommentsRemoved)
	}

	expectedRemovedLines := []int{}
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

	files, err := DiscoverFiles(tempDir, false, []string{})
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

func TestRemoveSingleLineMultilineComment(t *testing.T) {
	lang := Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	// Single-line multi-line comment
	line := "   /* This is a single-line multi-line comment */   "
	removed, content := RemoveSingleLineMultilineComment(line, lang)
	if !removed {
		t.Errorf("Expected single-line multi-line comment to be removed")
	}
	if content != line {
		t.Errorf("Expected content to match original line")
	}

	// Not a single-line multi-line comment (just the markers)
	line2 := "/* */"
	removed2, _ := RemoveSingleLineMultilineComment(line2, lang)
	if removed2 {
		t.Errorf("Did not expect just the markers to be removed as single-line comment")
	}
}

func TestProcessFile_RemoveSingleLineMultiline(t *testing.T) {
	content := `package main

func main() {
	fmt.Println("Hello")
	/* This is a multi-line comment
	   // This should NOT be removed
	   End of multi-line comment */
	fmt.Println("String with // comment inside")
	/* Single-line multi-line comment */
}`

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

	result, err := ProcessFile(tmpFile.Name(), lang, false, true, []string{})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	if result.CommentsRemoved != 1 {
		t.Errorf("Expected 1 comment removed, got %d", result.CommentsRemoved)
	}

	for _, comment := range result.RemovedComments {
		if !strings.Contains(comment.Content, "Single-line multi-line comment") {
			t.Errorf("Expected removed comment to be the single-line multi-line comment, got: %s", comment.Content)
		}
	}

	modifiedContent := strings.Join(result.ModifiedLines, "\n")
	if strings.Contains(modifiedContent, "/* Single-line multi-line comment */") {
		t.Error("Single-line multi-line comment should be removed")
	}
}

func TestPHPCommentRemoval(t *testing.T) {
	content := `<?php
// This is a single-line comment
/* This is a multi-line comment
   that spans multiple lines */

class Example {
    // Class property comment
    private $name = "test";
    
    /* Single-line multi-line comment */
    
    public function __construct() {
        // Constructor comment
        echo "Hello World"; // Inline comment
    }
}`

	tmpFile, err := os.CreateTemp("", "test_*.php")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	lang := Language{
		Name:            "PHP",
		Extensions:      []string{".php", ".phtml"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	result, err := ProcessFile(tmpFile.Name(), lang, false, false, []string{})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	if result.CommentsRemoved != 4 {
		t.Errorf("Expected 4 comments removed, got %d", result.CommentsRemoved)
	}

	modifiedContent := strings.Join(result.ModifiedLines, "\n")
	if !strings.Contains(modifiedContent, "/* This is a multi-line comment") {
		t.Error("Multi-line comment should be preserved")
	}
	if !strings.Contains(modifiedContent, "/* Single-line multi-line comment */") {
		t.Error("Single-line multi-line comment should be preserved when flag is false")
	}
}

func TestCSharpCommentRemoval(t *testing.T) {
	content := `using System;

// This is a single-line comment
/* This is a multi-line comment
   that spans multiple lines */

namespace Example
{
    // Class comment
    public class Program
    {
        // Property comment
        public string Name { get; set; }
        
        /* Single-line multi-line comment */
        
        public void Test()
        {
            string message = "String with // comment inside";
            Console.WriteLine(message);
        }
    }
}`

	tmpFile, err := os.CreateTemp("", "test_*.cs")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	lang := Language{
		Name:            "C#",
		Extensions:      []string{".cs"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	result, err := ProcessFile(tmpFile.Name(), lang, false, true, []string{})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	if result.CommentsRemoved != 4 {
		t.Errorf("Expected 4 comments removed, got %d", result.CommentsRemoved)
	}

	modifiedContent := strings.Join(result.ModifiedLines, "\n")
	if !strings.Contains(modifiedContent, "/* This is a multi-line comment") {
		t.Error("Multi-line comment should be preserved")
	}
	if strings.Contains(modifiedContent, "/* Single-line multi-line comment */") {
		t.Error("Single-line multi-line comment should be removed when flag is true")
	}
}

func TestIgnorePatterns(t *testing.T) {
	content := `// This is a regular comment
// @ts-ignore This should be preserved
// @deprecated This should be preserved
// TODO: This should be preserved
// FIXME: This should be preserved

function test() {
    // Regular inline comment
    const x = 42; // @ts-ignore inline ignore comment
    const y = 10; // Regular inline comment
    
    /* Regular multi-line comment */
    
    /* @ts-ignore multi-line ignore comment */
    
    return x + y;
}

// @ts-expect-error This should be preserved too`

	tmpFile, err := os.CreateTemp("", "test_*.ts")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	lang := Language{
		Name:            "TypeScript/JavaScript",
		Extensions:      []string{".ts", ".tsx", ".js", ".jsx"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	ignorePatterns := []string{"@ts-ignore", "@deprecated", "TODO", "FIXME", "@ts-expect-error"}

	result, err := ProcessFile(tmpFile.Name(), lang, false, true, ignorePatterns)
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	// Should only remove 1 comment: regular inline comment
	if result.CommentsRemoved != 1 {
		t.Errorf("Expected 1 comment removed, got %d", result.CommentsRemoved)
	}

	modifiedContent := strings.Join(result.ModifiedLines, "\n")

	// Check that ignore pattern comments are preserved
	preservedComments := []string{
		"@ts-ignore This should be preserved",
		"@deprecated This should be preserved",
		"TODO: This should be preserved",
		"FIXME: This should be preserved",
		"@ts-ignore inline ignore comment",
		"@ts-ignore multi-line ignore comment",
		"@ts-expect-error This should be preserved too",
	}

	for _, comment := range preservedComments {
		if !strings.Contains(modifiedContent, comment) {
			t.Errorf("Expected preserved comment containing: %s", comment)
		}
	}

	// Check that regular comments are removed
	removedComments := []string{
		"Regular inline comment",
	}

	for _, comment := range removedComments {
		if strings.Contains(modifiedContent, comment) {
			t.Errorf("Expected removed comment: %s", comment)
		}
	}
}

func TestShouldIgnoreComment(t *testing.T) {
	ignorePatterns := []string{"@ts-ignore", "TODO", "FIXME"}

	tests := []struct {
		name     string
		comment  string
		expected bool
	}{
		{
			name:     "standalone ignore comment",
			comment:  "// @ts-ignore This should be ignored",
			expected: true,
		},
		{
			name:     "inline ignore comment",
			comment:  "const x = 42; // @ts-ignore inline comment",
			expected: true,
		},
		{
			name:     "multi-line ignore comment",
			comment:  "/* @ts-ignore multi-line comment */",
			expected: true,
		},
		{
			name:     "TODO comment",
			comment:  "// TODO: Fix this later",
			expected: true,
		},
		{
			name:     "FIXME comment",
			comment:  "// FIXME: This needs attention",
			expected: true,
		},
		{
			name:     "regular comment",
			comment:  "// This is a regular comment",
			expected: false,
		},
		{
			name:     "SQL comment with ignore pattern",
			comment:  "-- @ts-ignore SQL comment",
			expected: true,
		},
		{
			name:     "SQL comment without ignore pattern",
			comment:  "-- Regular SQL comment",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldIgnoreComment(tt.comment, ignorePatterns)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for comment: %s", tt.expected, result, tt.comment)
			}
		})
	}
}
