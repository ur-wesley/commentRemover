package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type CommentRemovalResult struct {
	OriginalLines   int
	CommentsRemoved int
	RemainingLines  int
	ModifiedLines   []string
	RemovedComments []RemovedComment
}

type RemovedComment struct {
	LineNumber int
	Content    string
}

func ProcessFile(filePath string, lang Language, consecutive bool, removeSingleLineMultiline bool, ignorePatterns []string) (*CommentRemovalResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	var removedComments []RemovedComment
	scanner := bufio.NewScanner(file)

	const maxCapacity = 10 * 1024 * 1024
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, maxCapacity)

	lineNumber := 0
	inMultiLineComment := false

	var allLines []string
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		if strings.Contains(err.Error(), "token too long") {
			return nil, fmt.Errorf("file contains lines longer than %d MB (likely a minified/bundled file): %v", maxCapacity/(1024*1024), err)
		}
		return nil, err
	}

	for i, line := range allLines {
		lineNumber++
		originalLine := line

		if lang.MultiLineStart != "" && lang.MultiLineEnd != "" {
			inMultiLineComment = UpdateMultiLineCommentState(line, lang, inMultiLineComment)
		}

		isConsecutive := isPartOfConsecutiveComments(allLines, i, lang)

		processedLine, removed := RemoveSingleLineComment(line, lang, inMultiLineComment, consecutive, isConsecutive)

		if removed && len(ignorePatterns) > 0 && (strings.HasPrefix(strings.TrimSpace(originalLine), lang.SingleLineStart) || strings.Contains(originalLine, lang.SingleLineStart)) {
			if shouldIgnoreComment(originalLine, ignorePatterns) {
				removed = false
				processedLine = originalLine
			}
		}

		if !removed && removeSingleLineMultiline && lang.MultiLineStart != "" && lang.MultiLineEnd != "" {
			if singleLine, content := RemoveSingleLineMultilineComment(line, lang); singleLine {
				if len(ignorePatterns) > 0 && shouldIgnoreComment(content, ignorePatterns) {
					removed = false
					processedLine = content
				} else {
					removed = true
					processedLine = "REMOVE_LINE"
					originalLine = content
				}
			}
		}

		if removed {
			removedComments = append(removedComments, RemovedComment{
				LineNumber: lineNumber,
				Content:    originalLine,
			})
		}

		if processedLine != "REMOVE_LINE" {
			lines = append(lines, processedLine)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return &CommentRemovalResult{
		OriginalLines:   lineNumber,
		CommentsRemoved: len(removedComments),
		RemainingLines:  len(lines),
		ModifiedLines:   lines,
		RemovedComments: removedComments,
	}, nil
}

func UpdateMultiLineCommentState(line string, lang Language, currentState bool) bool {
	if lang.MultiLineStart == "" || lang.MultiLineEnd == "" {
		return false
	}

	inComment := currentState
	i := 0

	for i < len(line) {
		if !inComment && strings.HasPrefix(line[i:], lang.MultiLineStart) {
			inComment = true
			i += len(lang.MultiLineStart)
		} else if inComment && strings.HasPrefix(line[i:], lang.MultiLineEnd) {
			inComment = false
			i += len(lang.MultiLineEnd)
		} else {
			patternFound := false
			for _, pattern := range lang.AdditionalMultiLinePatterns {
				if !inComment && strings.HasPrefix(line[i:], pattern.Start) {
					inComment = true
					i += len(pattern.Start)
					patternFound = true
					break
				} else if inComment && strings.HasPrefix(line[i:], pattern.End) {
					inComment = false
					i += len(pattern.End)
					patternFound = true
					break
				}
			}
			if !patternFound {
				i++
			}
		}
	}

	return inComment
}

func isPartOfConsecutiveComments(lines []string, currentIndex int, lang Language) bool {
	if currentIndex < 0 || currentIndex >= len(lines) {
		return false
	}

	currentLine := strings.TrimSpace(lines[currentIndex])
	if !strings.HasPrefix(currentLine, lang.SingleLineStart) {
		return false
	}

	hasPreviousComment := false
	hasNextComment := false

	if currentIndex > 0 {
		prevLine := strings.TrimSpace(lines[currentIndex-1])
		hasPreviousComment = strings.HasPrefix(prevLine, lang.SingleLineStart)
	}

	if currentIndex < len(lines)-1 {
		nextLine := strings.TrimSpace(lines[currentIndex+1])
		hasNextComment = strings.HasPrefix(nextLine, lang.SingleLineStart)
	}

	return hasPreviousComment || hasNextComment
}

func RemoveSingleLineComment(line string, lang Language, inMultiLineComment bool, consecutive bool, isConsecutive bool) (string, bool) {
	if inMultiLineComment {
		return line, false
	}

	commentIndex := -1
	for i := 0; i <= len(line)-len(lang.SingleLineStart); i++ {
		if strings.HasPrefix(line[i:], lang.SingleLineStart) {
			if !IsInsideStringLiteral(line, i) {
				commentIndex = i
				break
			}
		}
	}

	if commentIndex == -1 {
		return line, false
	}

	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, lang.SingleLineStart) {
		if isConsecutive && !consecutive {
			return line, false
		}
		return "REMOVE_LINE", true
	}

	beforeComment := strings.TrimRightFunc(line[:commentIndex], func(r rune) bool {
		return r == ' ' || r == '\t'
	})

	if beforeComment == "" {
		if isConsecutive && !consecutive {
			return line, false
		}
		return "REMOVE_LINE", true
	}

	return beforeComment, true
}

func IsInsideStringLiteral(line string, pos int) bool {
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false

	for i := 0; i < pos && i < len(line); i++ {
		char := line[i]

		if i > 0 && line[i-1] == '\\' {
			continue
		}

		switch char {
		case '\'':
			if !inDoubleQuote && !inBacktick {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote && !inBacktick {
				inDoubleQuote = !inDoubleQuote
			}
		case '`':
			if !inSingleQuote && !inDoubleQuote {
				inBacktick = !inBacktick
			}
		}
	}

	return inSingleQuote || inDoubleQuote || inBacktick
}

func WriteFile(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}

	return nil
}

func RemoveSingleLineMultilineComment(line string, lang Language) (bool, string) {
	trimmed := strings.TrimSpace(line)
	if strings.HasPrefix(trimmed, lang.MultiLineStart) && strings.HasSuffix(trimmed, lang.MultiLineEnd) &&
		strings.Count(trimmed, lang.MultiLineStart) == 1 && strings.Count(trimmed, lang.MultiLineEnd) == 1 {
		inner := strings.TrimSpace(trimmed[len(lang.MultiLineStart) : len(trimmed)-len(lang.MultiLineEnd)])
		if inner != "" {
			return true, line
		}
	}
	return false, ""
}

func shouldIgnoreComment(commentLine string, ignorePatterns []string) bool {
	commentContent := strings.TrimSpace(commentLine)

	if strings.Contains(commentContent, "//") {
		lastCommentIndex := strings.LastIndex(commentContent, "//")
		if lastCommentIndex != -1 {
			commentContent = strings.TrimSpace(commentContent[lastCommentIndex+2:])
		}
	} else if strings.Contains(commentContent, "--") {
		lastCommentIndex := strings.LastIndex(commentContent, "--")
		if lastCommentIndex != -1 {
			commentContent = strings.TrimSpace(commentContent[lastCommentIndex+2:])
		}
	} else if strings.HasPrefix(commentContent, "/*") && strings.HasSuffix(commentContent, "*/") {
		commentContent = strings.TrimSpace(commentContent[2 : len(commentContent)-2])
	}

	for _, pattern := range ignorePatterns {
		if strings.Contains(commentContent, pattern) {
			return true
		}
	}

	return false
}
