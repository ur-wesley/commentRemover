package main

import (
	"fmt"
	"strings"

	ts "github.com/tree-sitter/go-tree-sitter"
	treeSitterC "github.com/tree-sitter/tree-sitter-c/bindings/go"
	treeSitterCpp "github.com/tree-sitter/tree-sitter-cpp/bindings/go"
	treeSitterGo "github.com/tree-sitter/tree-sitter-go/bindings/go"
	treeSitterJava "github.com/tree-sitter/tree-sitter-java/bindings/go"
	treeSitterJs "github.com/tree-sitter/tree-sitter-javascript/bindings/go"
	treeSitterPhp "github.com/tree-sitter/tree-sitter-php/bindings/go"
	treeSitterPython "github.com/tree-sitter/tree-sitter-python/bindings/go"
	treeSitterRuby "github.com/tree-sitter/tree-sitter-ruby/bindings/go"
	treeSitterRust "github.com/tree-sitter/tree-sitter-rust/bindings/go"
)

// TreeSitterProcessor handles comment detection using Tree-sitter
type TreeSitterProcessor struct {
	parsers map[string]*ts.Parser
}

// NewTreeSitterProcessor creates a new Tree-sitter processor
func NewTreeSitterProcessor() *TreeSitterProcessor {
	processors := map[string]*ts.Parser{
		"C":                   ts.NewParser(),
		"C++":                 ts.NewParser(),
		"Go":                  ts.NewParser(),
		"Java":                ts.NewParser(),
		"JavaScript":           ts.NewParser(),
		"TypeScript/JavaScript": ts.NewParser(), // alias
		"PHP":                 ts.NewParser(),
		"Python":              ts.NewParser(),
		"Ruby":                ts.NewParser(),
		"Rust":                ts.NewParser(),
		"TypeScript":          ts.NewParser(),
		"SQL":                 ts.NewParser(),
	}

	// Set language parsers - tree-sitter org bindings
	processors["C"].SetLanguage(ts.NewLanguage(treeSitterC.Language()))
	processors["C++"].SetLanguage(ts.NewLanguage(treeSitterCpp.Language()))
	processors["Go"].SetLanguage(ts.NewLanguage(treeSitterGo.Language()))
	processors["Java"].SetLanguage(ts.NewLanguage(treeSitterJava.Language()))
	processors["JavaScript"].SetLanguage(ts.NewLanguage(treeSitterJs.Language()))
	processors["TypeScript/JavaScript"].SetLanguage(ts.NewLanguage(treeSitterJs.Language())) // alias
	processors["PHP"].SetLanguage(ts.NewLanguage(treeSitterPhp.LanguagePHP()))
	processors["Python"].SetLanguage(ts.NewLanguage(treeSitterPython.Language()))
	processors["Ruby"].SetLanguage(ts.NewLanguage(treeSitterRuby.Language()))
	processors["Rust"].SetLanguage(ts.NewLanguage(treeSitterRust.Language()))
	processors["TypeScript"].SetLanguage(ts.NewLanguage(treeSitterJs.Language())) // TS uses JS grammar
	// SQL: no tree-sitter parser available

	return &TreeSitterProcessor{
		parsers: processors,
	}
}

// CommentNode represents a comment found by Tree-sitter
type CommentNode struct {
	StartLine   int
	EndLine     int
	StartColumn int
	EndColumn   int
	Content     string
	Type        string
}

// findCommentNodes finds all comment nodes in the tree
func (tp *TreeSitterProcessor) findCommentNodes(node *ts.Node, content []byte) []CommentNode {
	var comments []CommentNode

	// Check for comment types - tree-sitter uses "comment" kind
	if node.Kind() == "comment" {
		comments = append(comments, CommentNode{
			StartLine:   int(node.StartPosition().Row + 1),
			EndLine:     int(node.EndPosition().Row + 1),
			StartColumn: int(node.StartPosition().Column),
			EndColumn:   int(node.EndPosition().Column),
			Content:     string(content[node.StartByte():node.EndByte()]),
			Type:        node.Kind(),
		})
	}

	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		comments = append(comments, tp.findCommentNodes(child, content)...)
	}

	return comments
}

// ProcessFileWithTreeSitter processes a file using Tree-sitter for accurate comment detection
func (tp *TreeSitterProcessor) ProcessFileWithTreeSitter(content string, lang *Language, ignorePatterns []string) ([]string, int, error) {
	parser, ok := tp.parsers[lang.Name]
	if !ok {
		return nil, 0, fmt.Errorf("unsupported language: %s", lang.Name)
	}

	contentBytes := []byte(content)
	tree := parser.Parse(contentBytes, nil)
	if tree == nil {
		return nil, 0, fmt.Errorf("failed to parse content")
	}
	defer tree.Close()

	rootNode := tree.RootNode()
	comments := tp.findCommentNodes(rootNode, contentBytes)

	lines := strings.Split(content, "\n")
	modifiedLines := make([]string, len(lines))
	copy(modifiedLines, lines)

	commentsRemoved := 0

	for _, comment := range comments {
		if tp.shouldIgnoreComment(comment.Content, ignorePatterns) {
			continue
		}

		if comment.StartLine == comment.EndLine {
			line := lines[comment.StartLine-1]
			if comment.StartColumn == 0 {
				modifiedLines[comment.StartLine-1] = "REMOVE_LINE"
			} else {
				modifiedLines[comment.StartLine-1] = strings.TrimRight(line[:comment.StartColumn], " \t")
			}
		} else {
			for i := comment.StartLine - 1; i < comment.EndLine; i++ {
				if i == comment.StartLine-1 {
					if comment.StartColumn == 0 {
						modifiedLines[i] = "REMOVE_LINE"
					} else {
						modifiedLines[i] = strings.TrimRight(lines[i][:comment.StartColumn], " \t")
					}
				} else if i == comment.EndLine-1 {
					if comment.EndColumn >= len(lines[i]) {
						modifiedLines[i] = "REMOVE_LINE"
					} else {
						modifiedLines[i] = strings.TrimRight(lines[i][comment.EndColumn:], " \t")
					}
				} else {
					modifiedLines[i] = "REMOVE_LINE"
				}
			}
		}
		commentsRemoved++
	}

	finalLines := make([]string, 0, len(modifiedLines))
	for _, line := range modifiedLines {
		if line != "REMOVE_LINE" {
			finalLines = append(finalLines, line)
		}
	}

	return finalLines, commentsRemoved, nil
}

func (tp *TreeSitterProcessor) shouldIgnoreComment(commentContent string, ignorePatterns []string) bool {
	for _, pattern := range ignorePatterns {
		if strings.Contains(commentContent, pattern) {
			return true
		}
	}
	return false
}

func ProcessFileWithRegex(content string, lang *Language, ignorePatterns []string) ([]string, int, error) {
	lines := strings.Split(content, "\n")
	result, err := processFileWithRegex(lines, *lang, false, false, ignorePatterns)
	if err != nil {
		return nil, 0, err
	}
	return result.ModifiedLines, result.CommentsRemoved, nil
}