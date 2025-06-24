package main

import (
	"os"
	"strings"
	"testing"
)

func BenchmarkTreeSitterVsRegex(b *testing.B) {
	// Create test content with various comment types
	content := `package main

import "fmt"

// Single line comment
func main() {
    // Another comment
    fmt.Println("Hello") // Inline comment
    
    /* 
     * Multi-line comment
     * with multiple lines
     */
    
    str := "// This is not a comment"
    str2 := "/* This is also not a comment */"
    
    // More comments
    for i := 0; i < 10; i++ {
        // Loop comment
        fmt.Println(i) // Print comment
    }
}`

	lang := &Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	processor := NewTreeSitterProcessor()

	b.Run("TreeSitter", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := processor.ProcessFileWithTreeSitter(content, lang, []string{})
			if err != nil {
				b.Fatalf("Tree-sitter processing failed: %v", err)
			}
		}
	})

	b.Run("Regex", func(b *testing.B) {
		// Disable Tree-sitter for regex benchmark
		os.Setenv("COMMENTER_DISABLE_TREESITTER", "1")
		defer os.Unsetenv("COMMENTER_DISABLE_TREESITTER")

		lines := strings.Split(content, "\n")
		for i := 0; i < b.N; i++ {
			_, err := processFileWithRegex(lines, *lang, false, false, []string{})
			if err != nil {
				b.Fatalf("Regex processing failed: %v", err)
			}
		}
	})
}

func BenchmarkTreeSitterLargeFile(b *testing.B) {
	// Create a large file with many comments
	var lines []string
	lines = append(lines, "package main")
	lines = append(lines, "")
	lines = append(lines, "import (")
	lines = append(lines, "    \"fmt\"")
	lines = append(lines, "    \"strings\"")
	lines = append(lines, ")")
	lines = append(lines, "")
	lines = append(lines, "func main() {")

	for i := 0; i < 1000; i++ {
		if i%10 == 0 {
			lines = append(lines, "    // Comment "+string(rune(i)))
		}
		if i%50 == 0 {
			lines = append(lines, "    /*")
			lines = append(lines, "     * Multi-line comment "+string(rune(i)))
			lines = append(lines, "     */")
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

	processor := NewTreeSitterProcessor()

	b.Run("TreeSitterLargeFile", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := processor.ProcessFileWithTreeSitter(content, lang, []string{})
			if err != nil {
				b.Fatalf("Tree-sitter processing failed: %v", err)
			}
		}
	})
}

func BenchmarkTreeSitterDifferentLanguages(b *testing.B) {
	processor := NewTreeSitterProcessor()

	testCases := []struct {
		name    string
		content string
		lang    *Language
	}{
		{
			name: "Go",
			content: `package main

// Comment
func main() {
    fmt.Println("Hello") // Inline
}`,
			lang: &Language{
				Name:            "Go",
				Extensions:      []string{".go"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
		},
		{
			name: "JavaScript",
			content: `function test() {
    // Comment
    const x = 42; // Inline
}`,
			lang: &Language{
				Name:            "TypeScript/JavaScript",
				Extensions:      []string{".js"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
		},
		{
			name: "PHP",
			content: `<?php
// Comment
$var = "test"; // Inline
echo $var;`,
			lang: &Language{
				Name:            "PHP",
				Extensions:      []string{".php"},
				SingleLineStart: "//",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
		},
		{
			name: "SQL",
			content: `SELECT * FROM users
-- Comment
WHERE id = 1; -- Inline
SELECT name FROM users;`,
			lang: &Language{
				Name:            "SQL",
				Extensions:      []string{".sql"},
				SingleLineStart: "--",
				MultiLineStart:  "/*",
				MultiLineEnd:    "*/",
			},
		},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _, err := processor.ProcessFileWithTreeSitter(tc.content, tc.lang, []string{})
				if err != nil {
					b.Fatalf("Tree-sitter processing failed for %s: %v", tc.name, err)
				}
			}
		})
	}
}

func BenchmarkTreeSitterWithIgnorePatterns(b *testing.B) {
	processor := NewTreeSitterProcessor()

	content := `package main

// Regular comment
// @ts-ignore This should be preserved
// TODO: This should be preserved
// FIXME: This should be preserved
// Another regular comment
func main() {
    fmt.Println("Hello") // @deprecated This should be preserved
    // More regular comments
}`

	lang := &Language{
		Name:            "Go",
		Extensions:      []string{".go"},
		SingleLineStart: "//",
		MultiLineStart:  "/*",
		MultiLineEnd:    "*/",
	}

	ignorePatterns := []string{"@ts-ignore", "TODO", "FIXME", "@deprecated"}

	b.Run("WithIgnorePatterns", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := processor.ProcessFileWithTreeSitter(content, lang, ignorePatterns)
			if err != nil {
				b.Fatalf("Tree-sitter processing failed: %v", err)
			}
		}
	})

	b.Run("WithoutIgnorePatterns", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _, err := processor.ProcessFileWithTreeSitter(content, lang, []string{})
			if err != nil {
				b.Fatalf("Tree-sitter processing failed: %v", err)
			}
		}
	})
}
