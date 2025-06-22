# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- GoReleaser configuration for automated releases
- GitHub Actions CI/CD pipeline with multi-platform testing
- Comprehensive test suite with unit tests, benchmarks, and integration tests
- npm package support with automated npm publishing
- Winget package support for Windows package manager
- Homebrew and Scoop package manager support
- golangci-lint configuration for code quality
- Comprehensive documentation and examples

### Changed

- Modular architecture split into focused files:
  - `main.go` - CLI parsing and orchestration
  - `const.go` - Language definitions
  - `processor.go` - Core comment removal logic
  - `discovery.go` - File/directory discovery
  - `ui.go` - Terminal output and colors
- Enhanced CLI with short flags (`-w`, `-r`, `-nc`, `-h`)
- Improved directory processing with recursive support
- Better error handling and user feedback
- Cross-platform binary naming and installation

### Fixed

- String literal detection for complex escape sequences
- Multi-line comment handling
- File extension case sensitivity
- Windows executable naming in build processes

## [1.0.0] - Initial Release

### Added

- Basic CLI tool for removing single-line comments
- Support for Go, TypeScript/JavaScript, SQL, and JSON files
- Preview mode (default) and write mode (`--write`)
- Colored terminal output with emoji indicators
- Safe comment removal that preserves:
  - Multi-line comments
  - Comments inside string literals
  - Comments within multi-line comment blocks
- Cross-platform support (Windows, macOS, Linux)
- npm package distribution
- Professional help system and documentation

### Core Features

- **File Processing**: Single file comment removal
- **Language Support**: Go (//), TypeScript/JavaScript (//), SQL (--), JSON (//)
- **Safety**: Smart detection to avoid removing comments in strings or multi-line blocks
- **User Experience**: Beautiful colored output, clear statistics, helpful error messages
- **Installation**: Multiple installation methods (npm, Go install, manual)

### Technical Implementation

- Written in Go for performance and cross-platform compatibility
- Regex-based comment detection with context awareness
- Terminal color detection for CI/CD compatibility
- Comprehensive error handling and validation
