package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProcessingOptions struct {
	Write       bool
	NoColor     bool
	Recursive   bool
	Consecutive bool
	NoWarnLarge bool
	Extensions  []string
}

type ProcessingStats struct {
	FilesProcessed   int
	FilesSkipped     int
	TotalComments    int
	TotalLines       int
	SuccessfulWrites int
	FailedWrites     int
	Errors           []string
}

type FileInfo struct {
	Path     string
	Language Language
}

func DiscoverGlobFiles(pattern string) ([]FileInfo, error) {
	var files []FileInfo

	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid glob pattern '%s': %v", pattern, err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("no files match pattern: %s", pattern)
	}

	for _, match := range matches {
		if stat, err := os.Stat(match); err == nil && stat.IsDir() {
			continue
		}

		lang, supported := GetLanguageByExtension(match)
		if supported {
			files = append(files, FileInfo{
				Path:     match,
				Language: *lang,
			})
		}
	}

	return files, nil
}

func DiscoverFiles(inputPath string, recursive bool) ([]FileInfo, error) {
	var files []FileInfo

	if strings.Contains(inputPath, "*") || strings.Contains(inputPath, "?") || strings.Contains(inputPath, "[") {
		return DiscoverGlobFiles(inputPath)
	}

	stat, err := os.Stat(inputPath)
	if err != nil {
		return nil, fmt.Errorf("path does not exist: %s", inputPath)
	}

	if stat.IsDir() {
		err = processDirectory(inputPath, recursive, &files)
		if err != nil {
			return nil, err
		}
	} else {
		lang, supported := GetLanguageByExtension(inputPath)
		if !supported {
			return nil, fmt.Errorf("unsupported file type: %s", filepath.Ext(inputPath))
		}
		files = append(files, FileInfo{
			Path:     inputPath,
			Language: *lang,
		})
	}

	return files, nil
}

func processDirectory(dirPath string, recursive bool, files *[]FileInfo) error {
	if recursive {
		return filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if info.IsDir() {
				return nil
			}

			lang, supported := GetLanguageByExtension(path)
			if supported {
				*files = append(*files, FileInfo{
					Path:     path,
					Language: *lang,
				})
			}

			return nil
		})
	} else {
		entries, err := os.ReadDir(dirPath)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			fullPath := filepath.Join(dirPath, entry.Name())
			lang, supported := GetLanguageByExtension(fullPath)
			if supported {
				*files = append(*files, FileInfo{
					Path:     fullPath,
					Language: *lang,
				})
			}
		}
	}

	return nil
}

func FilterFilesByExtensions(files []FileInfo, extensions []string) []FileInfo {
	if len(extensions) == 0 {
		return files
	}

	var filtered []FileInfo
	for _, file := range files {
		fileExt := strings.ToLower(filepath.Ext(file.Path))
		for _, ext := range extensions {
			if strings.ToLower(ext) == fileExt {
				filtered = append(filtered, file)
				break
			}
		}
	}

	return filtered
}

func ProcessMultipleFiles(files []FileInfo, options ProcessingOptions, totalDuration time.Duration) *ProcessingStats {
	stats := &ProcessingStats{}
	useColor := !options.NoColor

	for _, file := range files {
		result, err := ProcessFile(file.Path, file.Language, options.Consecutive)
		if err != nil {
			stats.FailedWrites++
			stats.Errors = append(stats.Errors, fmt.Sprintf("%s: %v", file.Path, err))
			continue
		}

		if !options.NoWarnLarge && result.OriginalLines > 500 && len(files) > 1 {
			printWarning(useColor, "Large file: %s (%d lines)", file.Path, result.OriginalLines)
		}

		stats.FilesProcessed++
		stats.TotalComments += result.CommentsRemoved
		stats.TotalLines += result.OriginalLines

		if options.Write {
			if err := WriteFile(file.Path, result.ModifiedLines); err != nil {
				stats.FailedWrites++
				stats.Errors = append(stats.Errors, fmt.Sprintf("Failed to write %s: %v", file.Path, err))
			} else {
				stats.SuccessfulWrites++
			}
		}

		if len(files) == 1 {
			printFileResult(file.Path, file.Language, result, !options.NoColor, totalDuration, !options.NoWarnLarge)
		}
	}

	return stats
}

func printFileResult(filePath string, lang Language, result *CommentRemovalResult, useColor bool, duration time.Duration, showLargeWarning bool) {
	printInfo(useColor, "File: %s (%s)", filePath, lang.Name)

	if showLargeWarning && result.OriginalLines > 500 {
		printWarning(useColor, "Large file detected: %d lines (>500 LOC)", result.OriginalLines)
	}

	printStat(useColor, "Original lines", result.OriginalLines)
	printStat(useColor, "Comments removed", result.CommentsRemoved)
	printStat(useColor, "Remaining lines", result.RemainingLines)

	if len(result.RemovedComments) > 0 {
		fmt.Printf("\n%sRemoved comments:%s\n", colorize(useColor, ColorYellow+ColorBold), colorize(useColor, ColorReset))
		for _, comment := range result.RemovedComments {
			fmt.Printf("  %sLine %d:%s %s%s%s\n",
				colorize(useColor, ColorBlue),
				comment.LineNumber,
				colorize(useColor, ColorReset),
				colorize(useColor, ColorDim),
				strings.TrimSpace(comment.Content),
				colorize(useColor, ColorReset))
		}
	}

	printExecutionTime(useColor, duration)
}

func printBatchStats(stats *ProcessingStats, options ProcessingOptions) {
	useColor := !options.NoColor

	fmt.Printf("\n%sBatch Processing Summary:%s\n", colorize(useColor, ColorBold+ColorCyan), colorize(useColor, ColorReset))
	printStat(useColor, "Files processed", stats.FilesProcessed)
	printStat(useColor, "Total comments removed", stats.TotalComments)
	printStat(useColor, "Total lines processed", stats.TotalLines)

	if options.Write {
		printStat(useColor, "Files written successfully", stats.SuccessfulWrites)
		if stats.FailedWrites > 0 {
			fmt.Printf("%sFailed writes: %d%s\n", colorize(useColor, ColorRed), stats.FailedWrites, colorize(useColor, ColorReset))
		}
	}

	if len(stats.Errors) > 0 {
		fmt.Printf("\n%sErrors:%s\n", colorize(useColor, ColorRed+ColorBold), colorize(useColor, ColorReset))
		for _, err := range stats.Errors {
			fmt.Printf("  %s%s%s\n", colorize(useColor, ColorRed), err, colorize(useColor, ColorReset))
		}
	}
}
