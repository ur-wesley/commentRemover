package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type PackageJSON struct {
	Version string `json:"version"`
}

func getVersionFromPackageJSON() string {
	execPath, err := os.Executable()
	if err != nil {
		return version
	}

	packageJSONPath := filepath.Join(filepath.Dir(execPath), "package.json")

	if _, err := os.Stat(packageJSONPath); os.IsNotExist(err) {
		packageJSONPath = "package.json"
	}

	file, err := os.Open(packageJSONPath)
	if err != nil {
		return version
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return version
	}

	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return version
	}

	if pkg.Version != "" {
		return pkg.Version
	}

	return version
}

func main() {
	startTime := time.Now()

	var write bool
	var noColor bool
	var recursive bool
	var showHelp bool
	var showVersion bool
	var consecutive bool
	var noWarnLarge bool

	flag.BoolVar(&write, "write", false, "Write changes to file instead of just logging")
	flag.BoolVar(&write, "w", false, "Write changes to file (shorthand)")
	flag.BoolVar(&recursive, "recursive", true, "Process directories recursively (default: true)")
	flag.BoolVar(&recursive, "r", true, "Process directories recursively (shorthand, default: true)")
	flag.BoolVar(&consecutive, "consecutive", false, "Remove consecutive single-line comments (default: false)")
	flag.BoolVar(&consecutive, "c", false, "Remove consecutive single-line comments (shorthand)")
	flag.BoolVar(&noColor, "no-color", false, "Disable colored output")
	flag.BoolVar(&noColor, "nc", false, "Disable colored output (shorthand)")
	flag.BoolVar(&noWarnLarge, "no-warn-large", false, "Disable warnings for large files (>500 LOC)")
	flag.BoolVar(&noWarnLarge, "nwl", false, "Disable warnings for large files (shorthand)")
	flag.BoolVar(&showHelp, "help", false, "Show help message")
	flag.BoolVar(&showHelp, "h", false, "Show help message (shorthand)")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&showVersion, "v", false, "Show version information (shorthand)")
	flag.Parse()

	options := ProcessingOptions{
		Write:       write,
		NoColor:     noColor,
		Recursive:   recursive,
		Consecutive: consecutive,
		NoWarnLarge: noWarnLarge,
	}

	useColor := !noColor && isTerminal()

	if showVersion {
		fmt.Printf("%s version %s\n", filepath.Base(os.Args[0]), getVersionFromPackageJSON())
		fmt.Printf("Built: %s\n", date)
		fmt.Printf("Commit: %s\n", commit)
		os.Exit(0)
	}

	if showHelp {
		showHelpMessage(useColor)
		os.Exit(0)
	}

	var inputPath string
	if flag.NArg() < 1 {
		inputPath = "."
	} else {
		inputPath = flag.Arg(0)
	}

	files, err := DiscoverFiles(inputPath, recursive)
	if err != nil {
		printError(useColor, "%v", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		printError(useColor, "No supported files found in '%s'", inputPath)
		os.Exit(1)
	}

	duration := time.Since(startTime)

	stats := ProcessMultipleFiles(files, options, duration)

	if len(files) == 1 {
		if write {
			printSuccess(useColor, "File updated successfully!")
		} else {
			fmt.Printf("\n%sRun with --write to apply changes to the file.%s\n",
				colorize(useColor, ColorCyan),
				colorize(useColor, ColorReset))
		}
	} else {
		printBatchStats(stats, options)
		if !write {
			fmt.Printf("\n%sRun with --write to apply changes to all files.%s\n",
				colorize(useColor, ColorCyan),
				colorize(useColor, ColorReset))
		}
		printExecutionTime(useColor, duration)
	}

	if len(stats.Errors) > 0 {
		os.Exit(1)
	}
}
