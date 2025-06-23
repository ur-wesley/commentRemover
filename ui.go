package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorWhite  = "\033[37m"
	ColorBold   = "\033[1m"
	ColorDim    = "\033[2m"
)

func isTerminal() bool {
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
		return true
	}
	return false
}

func colorize(useColor bool, color string) string {
	if useColor {
		return color
	}
	return ""
}

func printError(useColor bool, format string, args ...interface{}) {
	prefix := fmt.Sprintf("%sError:%s ", colorize(useColor, ColorRed+ColorBold), colorize(useColor, ColorReset))
	fmt.Fprintf(os.Stderr, prefix+format+"\n", args...)
}

func printSuccess(useColor bool, format string, args ...interface{}) {
	prefix := fmt.Sprintf("%s✓%s ", colorize(useColor, ColorGreen+ColorBold), colorize(useColor, ColorReset))
	fmt.Printf(prefix+format+"\n", args...)
}

func printInfo(useColor bool, format string, args ...interface{}) {
	prefix := fmt.Sprintf("%s📁%s ", colorize(useColor, ColorBlue), colorize(useColor, ColorReset))
	fmt.Printf(prefix+format+"\n", args...)
}

func printWarning(useColor bool, format string, args ...interface{}) {
	prefix := fmt.Sprintf("%s⚠%s ", colorize(useColor, ColorYellow+ColorBold), colorize(useColor, ColorReset))
	fmt.Printf(prefix+format+"\n", args...)
}

func printStat(useColor bool, label string, value int) {
	fmt.Printf("%s%s:%s %s%d%s\n",
		colorize(useColor, ColorCyan),
		label,
		colorize(useColor, ColorReset),
		colorize(useColor, ColorBold),
		value,
		colorize(useColor, ColorReset))
}

func showHelpMessage(useColor bool) {
	programName := filepath.Base(os.Args[0])

	fmt.Printf("%s%s%s - Comment Remover\n", colorize(useColor, ColorBold+ColorBlue), programName, colorize(useColor, ColorReset))
	fmt.Printf("A performant CLI tool that safely removes single-line comments from source code files.\n\n")

	fmt.Printf("%sUSAGE:%s\n", colorize(useColor, ColorBold+ColorYellow), colorize(useColor, ColorReset))
	fmt.Printf("  %s [OPTIONS] [<file/path/pattern>]\n", programName)
	fmt.Printf("  %s                              # Process current directory recursively%s\n", programName, colorize(useColor, ColorDim))
	fmt.Printf("  %s src/                         # Process src directory recursively%s\n", programName, colorize(useColor, ColorDim))
	fmt.Printf("  %s \"*.go\"                       # Process all .go files in current directory%s\n", programName, colorize(useColor, ColorDim))
	fmt.Printf("  %s \"./src/**/*.ts\"               # Process all .ts files recursively in src%s\n\n", programName, colorize(useColor, ColorReset))

	fmt.Printf("%sOPTIONS:%s\n", colorize(useColor, ColorBold+ColorYellow), colorize(useColor, ColorReset))
	fmt.Printf("  %s-w, --write%s      Write changes to file instead of just logging\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-r, --recursive%s  Process directories recursively (default: true)\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-c, --consecutive%s Remove consecutive single-line comments (default: false)\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-e, --exclude%s    Comma-separated glob patterns to exclude (e.g., '*test.go,*.min.js')\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-i, --ignore-pattern%s Comma-separated patterns to ignore in comments (e.g., '@ts-ignore,@deprecated')\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-nc, --no-color%s  Disable colored output\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-nwl, --no-warn-large%s Disable warnings for large files (>500 LOC)\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-h, --help%s       Show this help message\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-v, --version%s    Show version information\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s-m, --remove-single-multiline%s Remove single-line comments using multi-line patterns (e.g., /* comment */)\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("\n")

	fmt.Printf("%sSUPPORTED FILE TYPES:%s\n", colorize(useColor, ColorBold+ColorYellow), colorize(useColor, ColorReset))
	for _, lang := range SupportedLanguages {
		fmt.Printf("  %s%s%s: %s (comment: %s%s%s)\n",
			colorize(useColor, ColorCyan),
			lang.Name,
			colorize(useColor, ColorReset),
			strings.Join(lang.Extensions, ", "),
			colorize(useColor, ColorDim),
			lang.SingleLineStart,
			colorize(useColor, ColorReset))
	}

	fmt.Printf("\n%sEXAMPLES:%s\n", colorize(useColor, ColorBold+ColorYellow), colorize(useColor, ColorReset))
	fmt.Printf("  %s%s                             # Process current directory recursively (default)\n", programName, "")
	fmt.Printf("  %s example.go%s                  # Preview comment removal from single file\n", programName, "")
	fmt.Printf("  %s %s-w%s example.ts%s               # Remove comments and save file\n", programName, colorize(useColor, ColorGreen), colorize(useColor, ColorReset), "")
	fmt.Printf("  %s src/%s                        # Process src directory recursively\n", programName, "")
	fmt.Printf("  %s %s-w -r%s project/%s              # Recursively process and save all files\n", programName, colorize(useColor, ColorGreen), colorize(useColor, ColorReset), "")
	fmt.Printf("  %s \"*.go\"%s                      # Process all .go files in current directory\n", programName, "")
	fmt.Printf("  %s \"./src/**/*.ts\"%s             # Process all .ts files recursively in src\n", programName, "")
	fmt.Printf("  %s %s-w%s \"src/**/*.{ts,js}\"%s       # Process and save .ts/.js files in src\n", programName, colorize(useColor, ColorGreen), colorize(useColor, ColorReset), "")
	fmt.Printf("  %s %s-e%s \"*test.go,*.min.js\"%s       # Exclude test files and minified files\n", programName, colorize(useColor, ColorGreen), colorize(useColor, ColorReset), "")
	fmt.Printf("  %s %s-i%s \"@ts-ignore,@deprecated\"%s   # Ignore comments with specific patterns\n", programName, colorize(useColor, ColorGreen), colorize(useColor, ColorReset), "")
	fmt.Printf("  %s %s-c%s file.ts%s                   # Remove consecutive comments too\n", programName, colorize(useColor, ColorGreen), colorize(useColor, ColorReset), "")
	fmt.Printf("  %s %s-w -nc%s src/utils.js%s          # Save with no colors\n", programName, colorize(useColor, ColorGreen), colorize(useColor, ColorReset), "")
	fmt.Printf("  %s %s-m%s file.js%s                   # Remove single-line multi-line comments\n", programName, colorize(useColor, ColorGreen), colorize(useColor, ColorReset), "")
	fmt.Printf("\n")

	fmt.Printf("%sWHAT GETS REMOVED:%s\n", colorize(useColor, ColorBold+ColorYellow), colorize(useColor, ColorReset))
	fmt.Printf("  %s✓%s Standalone comment lines (e.g., %s// This is a comment%s)\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset), colorize(useColor, ColorDim), colorize(useColor, ColorReset))
	fmt.Printf("  %s✓%s Inline comments (e.g., %scode(); // comment%s)\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset), colorize(useColor, ColorDim), colorize(useColor, ColorReset))
	fmt.Printf("  %s✓%s Multiple consecutive single-line comments\n\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset))
	fmt.Printf("  %s✓%s Single-line multi-line comments (e.g., %s/* comment */%s) [with -m]\n\n", colorize(useColor, ColorGreen), colorize(useColor, ColorReset), colorize(useColor, ColorDim), colorize(useColor, ColorReset))

	fmt.Printf("%sWHAT GETS PRESERVED:%s\n", colorize(useColor, ColorBold+ColorYellow), colorize(useColor, ColorReset))
	fmt.Printf("  %s×%s Multi-line comments (%s/* ... */%s)\n", colorize(useColor, ColorRed), colorize(useColor, ColorReset), colorize(useColor, ColorDim), colorize(useColor, ColorReset))
	fmt.Printf("  %s×%s Comments inside string literals (%s\"string with // comment\"%s)\n", colorize(useColor, ColorRed), colorize(useColor, ColorReset), colorize(useColor, ColorDim), colorize(useColor, ColorReset))
	fmt.Printf("  %s×%s Single-line comments inside multi-line comment blocks\n", colorize(useColor, ColorRed), colorize(useColor, ColorReset))
}

func printExecutionTime(useColor bool, duration time.Duration) {
	var timeStr string
	var unit string

	if duration < time.Microsecond {
		timeStr = fmt.Sprintf("%.0f", float64(duration.Nanoseconds()))
		unit = "ns"
	} else if duration < time.Millisecond {
		timeStr = fmt.Sprintf("%.0f", float64(duration.Nanoseconds())/1000)
		unit = "µs"
	} else if duration < time.Second {
		timeStr = fmt.Sprintf("%.2f", float64(duration.Nanoseconds())/1000000)
		unit = "ms"
	} else {
		timeStr = fmt.Sprintf("%.3f", duration.Seconds())
		unit = "s"
	}

	fmt.Printf("\n%sExecution time:%s %s%s%s%s\n",
		colorize(useColor, ColorDim),
		colorize(useColor, ColorReset),
		colorize(useColor, ColorBold),
		timeStr,
		unit,
		colorize(useColor, ColorReset))
}
