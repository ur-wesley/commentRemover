package main

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/php"
	"github.com/smacker/go-tree-sitter/sql"
)

// TreeSitterProcessor handles comment detection using Tree-sitter
type TreeSitterProcessor struct {
	parsers map[string]*sitter.Parser
}

// NewTreeSitterProcessor creates a new Tree-sitter processor
func NewTreeSitterProcessor() *TreeSitterProcessor {
	processors := map[string]*sitter.Parser{
		"Go":                    sitter.NewParser(),
		"TypeScript/JavaScript": sitter.NewParser(),
		"PHP":                   sitter.NewParser(),
		"SQL":                   sitter.NewParser(),
	}

	// Set language parsers
	processors["Go"].SetLanguage(golang.GetLanguage())
	processors["TypeScript/JavaScript"].SetLanguage(javascript.GetLanguage())
	processors["PHP"].SetLanguage(php.GetLanguage())
	processors["SQL"].SetLanguage(sql.GetLanguage())

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
	Type        string // "single_line", "multi_line", "documentation"
}

// ProcessFileWithTreeSitter processes a file using Tree-sitter for accurate comment detection
func (tp *TreeSitterProcessor) ProcessFileWithTreeSitter(content string, lang *Language, ignorePatterns []string) ([]string, int, error) {
	parser, exists := tp.parsers[lang.Name]
	if !exists {
		// Return error to trigger fallback to regex-based processing for unsupported languages
		return nil, 0, fmt.Errorf("no Tree-sitter parser available for language: %s", lang.Name)
	}

	// Parse the content
	tree := parser.Parse(nil, []byte(content))
	defer tree.Close()

	// Get the root node
	rootNode := tree.RootNode()

	// Find all comment nodes
	comments := tp.findCommentNodes(rootNode, content)

	// Process lines with comment information
	lines := strings.Split(content, "\n")
	modifiedLines := make([]string, len(lines))
	copy(modifiedLines, lines)

	commentsRemoved := 0

	for _, comment := range comments {
		// Check if comment should be ignored
		if tp.shouldIgnoreComment(comment.Content, ignorePatterns) {
			continue
		}

		// Remove the comment
		if comment.StartLine == comment.EndLine {
			// Single line comment
			line := lines[comment.StartLine-1]
			if comment.StartColumn == 0 {
				// Full line comment
				modifiedLines[comment.StartLine-1] = "REMOVE_LINE"
			} else {
				// Inline comment - preserve indentation
				modifiedLines[comment.StartLine-1] = strings.TrimRight(line[:comment.StartColumn], " \t")
			}
		} else {
			// Multi-line comment
			for i := comment.StartLine - 1; i < comment.EndLine; i++ {
				if i == comment.StartLine-1 {
					// First line of multi-line comment
					if comment.StartColumn == 0 {
						modifiedLines[i] = "REMOVE_LINE"
					} else {
						// Preserve indentation
						modifiedLines[i] = strings.TrimRight(lines[i][:comment.StartColumn], " \t")
					}
				} else if i == comment.EndLine-1 {
					// Last line of multi-line comment
					if comment.EndColumn >= len(lines[i]) {
						modifiedLines[i] = "REMOVE_LINE"
					} else {
						// Preserve indentation
						modifiedLines[i] = strings.TrimRight(lines[i][comment.EndColumn:], " \t")
					}
				} else {
					// Middle lines of multi-line comment
					modifiedLines[i] = "REMOVE_LINE"
				}
			}
		}
		commentsRemoved++
	}

	// Clean up REMOVE_LINE markers
	finalLines := make([]string, 0, len(modifiedLines))
	for _, line := range modifiedLines {
		if line != "REMOVE_LINE" {
			finalLines = append(finalLines, line)
		}
	}

	return finalLines, commentsRemoved, nil
}

// findCommentNodes recursively finds all comment nodes in the syntax tree
func (tp *TreeSitterProcessor) findCommentNodes(node *sitter.Node, content string) []CommentNode {
	var comments []CommentNode
	processedNodes := make(map[*sitter.Node]bool)

	var findComments func(*sitter.Node)
	findComments = func(n *sitter.Node) {
		// Check if current node is a comment and hasn't been processed
		if tp.isCommentNode(n) && !processedNodes[n] {
			comment := tp.nodeToComment(n, content)
			comments = append(comments, comment)
			processedNodes[n] = true
		}

		// Recursively check children
		for i := 0; i < int(n.ChildCount()); i++ {
			child := n.Child(i)
			findComments(child)
		}
	}

	findComments(node)
	return comments
}

// isCommentNode checks if a node represents a comment
func (tp *TreeSitterProcessor) isCommentNode(node *sitter.Node) bool {
	nodeType := node.Type()
	return strings.Contains(nodeType, "comment") ||
		nodeType == "comment" ||
		nodeType == "line_comment" ||
		nodeType == "block_comment" ||
		nodeType == "documentation_comment" ||
		nodeType == "comment_block" ||
		nodeType == "line_comment_block"
}

// nodeToComment converts a Tree-sitter node to a CommentNode
func (tp *TreeSitterProcessor) nodeToComment(node *sitter.Node, content string) CommentNode {
	startPoint := node.StartPoint()
	endPoint := node.EndPoint()

	// Convert Tree-sitter points to line/column numbers
	startLine := int(startPoint.Row) + 1
	endLine := int(endPoint.Row) + 1
	startColumn := int(startPoint.Column)
	endColumn := int(endPoint.Column)

	// Extract comment content
	lines := strings.Split(content, "\n")
	var commentContent string

	if startLine == endLine {
		// Single line comment
		line := lines[startLine-1]
		commentContent = line[startColumn:endColumn]
	} else {
		// Multi-line comment
		var parts []string
		for i := startLine - 1; i < endLine; i++ {
			if i == startLine-1 {
				parts = append(parts, lines[i][startColumn:])
			} else if i == endLine-1 {
				parts = append(parts, lines[i][:endColumn])
			} else {
				parts = append(parts, lines[i])
			}
		}
		commentContent = strings.Join(parts, "\n")
	}

	// Determine comment type
	commentType := "multi_line"
	if startLine == endLine {
		commentType = "single_line"
	}
	if strings.Contains(node.Type(), "documentation") {
		commentType = "documentation"
	}

	return CommentNode{
		StartLine:   startLine,
		EndLine:     endLine,
		StartColumn: startColumn,
		EndColumn:   endColumn,
		Content:     commentContent,
		Type:        commentType,
	}
}

// shouldIgnoreComment checks if a comment should be ignored based on patterns
func (tp *TreeSitterProcessor) shouldIgnoreComment(commentContent string, ignorePatterns []string) bool {
	for _, pattern := range ignorePatterns {
		if strings.Contains(commentContent, pattern) {
			return true
		}
	}
	return false
}

// ProcessFileWithRegex is a fallback function for unsupported languages
func ProcessFileWithRegex(content string, lang *Language, ignorePatterns []string) ([]string, int, error) {
	lines := strings.Split(content, "\n")

	// Use the existing regex-based processing
	result, err := processFileWithRegex(lines, *lang, false, false, ignorePatterns)
	if err != nil {
		return nil, 0, err
	}

	return result.ModifiedLines, result.CommentsRemoved, nil
}
