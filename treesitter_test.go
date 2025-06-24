package main

import (
	"reflect"
	"strings"
	"testing"
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

func cleanEmptyLines(lines []string) []string {
	var result []string
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	return result
}
