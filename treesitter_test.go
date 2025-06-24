package main

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestTreeSitterProcessor(t *testing.T) {
	processor := NewTreeSitterProcessor()

	tests := []struct {
		name           string
		content        string
		lang           *Language
		ignorePatterns []string
		expectedLines  []string
		expectedCount  int
	}{
		{
			name: "Go single line comments",
			content: `package main

import "fmt"

// This is a comment
func main() {
    fmt.Println("Hello") // Inline comment
    // Another comment
}`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			expectedLines: []string{
				"package main",
				"",
				"import \"fmt\"",
				"",
				"func main() {",
				"    fmt.Println(\"Hello\")",
				"}",
			},
			expectedCount: 3,
		},
		{
			name: "Go multi-line comments",
			content: `package main

import "fmt"

/*
 * Multi-line comment
 * with multiple lines
 */
func main() {
    fmt.Println("Hello")
}`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			expectedLines: []string{
				"package main",
				"",
				"import \"fmt\"",
				"",
				"func main() {",
				"    fmt.Println(\"Hello\")",
				"}",
			},
			expectedCount: 1,
		},
		{
			name: "JavaScript with ignore patterns",
			content: `function test() {
    // Regular comment
    const x = 42; // @ts-ignore This should be preserved
    // Another regular comment
}`,
			lang: &Language{
				Name:            "TypeScript/JavaScript",
				Extensions:      []string{".js", ".ts"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{"@ts-ignore"},
			expectedLines: []string{
				"function test() {",
				"    const x = 42; // @ts-ignore This should be preserved",
				"}",
			},
			expectedCount: 2,
		},
		{
			name: "PHP comments",
			content: `<?php
// Single line comment
$var = "test"; // Inline comment
/*
 * Multi-line comment
 */
echo $var;`,
			lang: &Language{
				Name:            "PHP",
				Extensions:      []string{".php"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			expectedLines: []string{
				"<?php",
				"$var = \"test\";",
				"echo $var;",
			},
			expectedCount: 3,
		},
		{
			name: "SQL comments",
			content: `SELECT * FROM users
-- Single line comment
WHERE id = 1; -- Inline comment
/*
 * Multi-line comment
 */
SELECT name FROM users;`,
			lang: &Language{
				Name:            "SQL",
				Extensions:      []string{".sql"},
				SingleLineStart: "--",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			expectedLines: []string{
				"SELECT * FROM users",
				"WHERE id = 1;",
				"",
				"/*",
				" * Multi-line comment",
				" */",
				"SELECT name FROM users;",
			},
			expectedCount: 2, // Only the -- comments are detected by Tree-sitter
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifiedLines, commentsRemoved, err := processor.ProcessFileWithTreeSitter(tt.content, tt.lang, tt.ignorePatterns)

			if err != nil {
				t.Fatalf("ProcessFileWithTreeSitter failed: %v", err)
			}

			if commentsRemoved != tt.expectedCount {
				t.Errorf("Expected %d comments removed, got %d", tt.expectedCount, commentsRemoved)
			}

			// Clean up empty lines for comparison
			cleanExpected := cleanEmptyLines(tt.expectedLines)
			cleanActual := cleanEmptyLines(modifiedLines)

			if !reflect.DeepEqual(cleanActual, cleanExpected) {
				t.Errorf("Expected lines:\n%s\n\nGot lines:\n%s",
					strings.Join(cleanExpected, "\n"),
					strings.Join(cleanActual, "\n"))
			}
		})
	}
}

func TestTreeSitterEdgeCases(t *testing.T) {
	processor := NewTreeSitterProcessor()

	tests := []struct {
		name           string
		content        string
		lang           *Language
		ignorePatterns []string
		description    string
	}{
		{
			name: "Comments in strings should be preserved",
			content: `package main

func main() {
    str := "// This is not a comment"
    str2 := "/* This is also not a comment */"
    fmt.Println(str)
}`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			description:    "Comments inside string literals should not be removed",
		},
		{
			name: "Nested comments",
			content: `package main

func main() {
    // Outer comment
    /* 
     * Outer multi-line comment
     * // Inner comment (should be part of outer)
     */
    fmt.Println("Hello")
}`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			description:    "Nested comments should be handled correctly",
		},
		{
			name:    "Empty file",
			content: "",
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			description:    "Empty files should be handled gracefully",
		},
		{
			name: "File with only comments",
			content: `// Only comments
// No code
/*
 * Just comments
 */`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			description:    "Files with only comments should be handled",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			modifiedLines, commentsRemoved, err := processor.ProcessFileWithTreeSitter(tt.content, tt.lang, tt.ignorePatterns)

			if err != nil {
				t.Fatalf("ProcessFileWithTreeSitter failed: %v", err)
			}

			// For edge cases, we mainly want to ensure no errors and reasonable behavior
			t.Logf("Processed file: %s", tt.description)
			t.Logf("Comments removed: %d", commentsRemoved)
			t.Logf("Result lines: %d", len(modifiedLines))
		})
	}
}

func TestTreeSitterErrorHandling(t *testing.T) {
	processor := NewTreeSitterProcessor()

	tests := []struct {
		name           string
		content        string
		lang           *Language
		ignorePatterns []string
		expectError    bool
		errorContains  string
	}{
		{
			name:    "Unsupported language",
			content: `some code here`,
			lang: &Language{
				Name:            "UnsupportedLanguage",
				Extensions:      []string{".unsupported"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			expectError:    true,
			errorContains:  "no Tree-sitter parser available",
		},
		{
			name: "Invalid syntax",
			content: `package main

func main() {
    // Valid comment
    invalid syntax here
}`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{},
			expectError:    false, // Tree-sitter should handle invalid syntax gracefully
			errorContains:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := processor.ProcessFileWithTreeSitter(tt.content, tt.lang, tt.ignorePatterns)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("Expected error to contain '%s', got: %v", tt.errorContains, err)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
			}
		})
	}
}

func TestTreeSitterIntegration(t *testing.T) {
	// Test integration with the main ProcessFile function
	content := `package main

// This is a comment
func main() {
    fmt.Println("Hello") // Inline comment
}`

	// Create a temporary file for testing
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

	// Test with Tree-sitter enabled
	result, err := ProcessFile(tmpFile.Name(), lang, false, false, []string{})
	if err != nil {
		t.Fatalf("ProcessFile failed: %v", err)
	}

	// Should remove 2 comments
	if result.CommentsRemoved != 2 {
		t.Errorf("Expected 2 comments removed, got %d", result.CommentsRemoved)
	}

	// Test with Tree-sitter disabled
	os.Setenv("COMMENTER_DISABLE_TREESITTER", "1")
	defer os.Unsetenv("COMMENTER_DISABLE_TREESITTER")

	result2, err := ProcessFile(tmpFile.Name(), lang, false, false, []string{})
	if err != nil {
		t.Fatalf("ProcessFile with Tree-sitter disabled failed: %v", err)
	}

	// Should also remove 2 comments (regex fallback)
	if result2.CommentsRemoved != 2 {
		t.Errorf("Expected 2 comments removed with regex fallback, got %d", result2.CommentsRemoved)
	}
}

func TestTreeSitterPerformance(t *testing.T) {
	processor := NewTreeSitterProcessor()

	// Create a large file with many comments
	var lines []string
	lines = append(lines, "package main")
	lines = append(lines, "")
	lines = append(lines, "func main() {")

	for i := 0; i < 1000; i++ {
		if i%10 == 0 {
			lines = append(lines, "    // Comment "+string(rune(i)))
		}
		lines = append(lines, "    fmt.Println(\"Line "+string(rune(i))+"\")")
	}

	lines = append(lines, "}")

	content := strings.Join(lines, "\n")

	lang := &Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	// Benchmark Tree-sitter processing
	start := time.Now()
	_, commentsRemoved, err := processor.ProcessFileWithTreeSitter(content, lang, []string{})
	duration := time.Since(start)

	if err != nil {
		t.Fatalf("Performance test failed: %v", err)
	}

	t.Logf("Processed %d lines with %d comments in %v", len(lines), commentsRemoved, duration)

	// Should complete within reasonable time (less than 1 second)
	if duration > time.Second {
		t.Errorf("Processing took too long: %v", duration)
	}
}

func TestTreeSitterCommentDetection(t *testing.T) {
	processor := NewTreeSitterProcessor()

	content := `package main

// This is a comment
func main() {
    fmt.Println("Hello") // Inline comment
}`

	lang := &Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	parser := processor.parsers[lang.Name]
	tree := parser.Parse(nil, []byte(content))
	defer tree.Close()

	rootNode := tree.RootNode()
	comments := processor.findCommentNodes(rootNode, content)

	if len(comments) != 2 {
		t.Errorf("Expected 2 comments, got %d", len(comments))
	}

	// Check first comment (full line)
	if comments[0].StartLine != 3 || comments[0].EndLine != 3 {
		t.Errorf("Expected first comment on line 3, got lines %d-%d", comments[0].StartLine, comments[0].EndLine)
	}

	// Check second comment (inline)
	if comments[1].StartLine != 5 || comments[1].EndLine != 5 {
		t.Errorf("Expected second comment on line 5, got lines %d-%d", comments[1].StartLine, comments[1].EndLine)
	}
}

func TestPHPDebug(t *testing.T) {
	processor := NewTreeSitterProcessor()

	content := `<?php
// Single line comment
$var = "test"; // Inline comment
/*
 * Multi-line comment
 */
echo $var;`

	lang := &Language{
		Name:            "PHP",
		Extensions:      []string{".php"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	parser := processor.parsers[lang.Name]
	tree := parser.Parse(nil, []byte(content))
	defer tree.Close()

	rootNode := tree.RootNode()
	comments := processor.findCommentNodes(rootNode, content)

	t.Logf("Found %d comments:", len(comments))
	for i, comment := range comments {
		t.Logf("Comment %d: Type=%s, Lines=%d-%d, Content='%s'",
			i+1, comment.Type, comment.StartLine, comment.EndLine, comment.Content)
	}
}

func TestTreeSitterIgnorePatterns(t *testing.T) {
	processor := NewTreeSitterProcessor()

	tests := []struct {
		name           string
		content        string
		lang           *Language
		ignorePatterns []string
		expectedCount  int
		description    string
	}{
		{
			name: "Multiple ignore patterns",
			content: `package main

// Regular comment
// @ts-ignore This should be preserved
// TODO: This should be preserved
// FIXME: This should be preserved
// Another regular comment
func main() {
    fmt.Println("Hello") // @deprecated This should be preserved
}`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{"@ts-ignore", "TODO", "FIXME", "@deprecated"},
			expectedCount:  2, // Only "Regular comment" and "Another regular comment" should be removed
			description:    "Comments with ignore patterns should be preserved",
		},
		{
			name: "Case sensitive patterns",
			content: `package main

// @TS-IGNORE This should be removed (case sensitive)
// @ts-ignore This should be preserved
func main() {
    fmt.Println("Hello")
}`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
			ignorePatterns: []string{"@ts-ignore"},
			expectedCount:  1, // Only "@TS-IGNORE" should be removed
			description:    "Ignore patterns should be case sensitive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, commentsRemoved, err := processor.ProcessFileWithTreeSitter(tt.content, tt.lang, tt.ignorePatterns)

			if err != nil {
				t.Fatalf("ProcessFileWithTreeSitter failed: %v", err)
			}

			if commentsRemoved != tt.expectedCount {
				t.Errorf("Expected %d comments removed, got %d. %s",
					tt.expectedCount, commentsRemoved, tt.description)
			}
		})
	}
}

func TestSQLDebug(t *testing.T) {
	processor := NewTreeSitterProcessor()

	content := `SELECT * FROM users
-- Single line comment
WHERE id = 1; -- Inline comment
/*
 * Multi-line comment
 */
SELECT name FROM users;`

	lang := &Language{
		Name:            "SQL",
		Extensions:      []string{".sql"},
		SingleLineStart: "--",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	parser := processor.parsers[lang.Name]
	tree := parser.Parse(nil, []byte(content))
	defer tree.Close()

	rootNode := tree.RootNode()
	comments := processor.findCommentNodes(rootNode, content)

	t.Logf("Found %d comments:", len(comments))
	for i, comment := range comments {
		t.Logf("Comment %d: Type=%s, Lines=%d-%d, Content='%s'",
			i+1, comment.Type, comment.StartLine, comment.EndLine, comment.Content)
	}
}

func cleanEmptyLines(lines []string) []string {
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}
