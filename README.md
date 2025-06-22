# Comment Remover (commenter)

A performant CLI tool written in Go that safely removes single-line comments from source code files.

> **Note**: This tool was renamed from `cr` to `commenter` to avoid conflicts with the existing Unix `cr` command that converts text files between Unix and DOS line endings.

## Features

- **Safe comment removal**: Only removes single-line comments that are not part of multi-line comments or inside string literals
- **Multiple language support**: TypeScript/JavaScript, Go, SQL, and JSON
- **Performance optimized**: Fast file processing with minimal memory usage
- **Preview mode**: See what would be removed before making changes
- **Detailed logging**: Shows line numbers and content of removed comments
- **Colored output**: Beautiful colored terminal output with emoji icons (can be disabled with `--no-color`)
- **Easy installation**: Available via npm/bun, or build from source

## Supported Languages

| Language              | Extensions                   | Single-line Comment |
| --------------------- | ---------------------------- | ------------------- |
| TypeScript/JavaScript | `.ts`, `.tsx`, `.js`, `.jsx` | `//`                |
| Go                    | `.go`                        | `//`                |
| SQL                   | `.sql`                       | `--`                |
| JSON                  | `.json`                      | `//`                |

## Installation

### Via npm (recommended)

```bash
# Install globally with npm
npm install -g commenter

# Or use with npx (no installation required)
npx commenter <file/path>

# Or use with bun
bun add -g commenter
bunx commenter <file/path>
```

### Via Go (build from source)

```bash
go install github.com/ur-wesley/commentRemover@latest
```

### Manual Installation

1. Download the binary from releases
2. Add to your PATH
3. Or build from source: `make install`

## Usage

```bash
# Preview what comments would be removed (default behavior)
commenter <file/path>

# Process a directory (non-recursive)
commenter src/

# Process a directory recursively
commenter -r project/
commenter --recursive project/        # Long flag

# Actually remove comments and update files
commenter --write <file/path>
commenter -w <file/path>              # Short flag
commenter -w -r src/                  # Write changes recursively

# Disable colored output
commenter --no-color <file/path>
commenter -nc <file/path>             # Short flag

# Show help
commenter --help
commenter -h                          # Short flag

# Combine flags
commenter -w -nc <file/path>          # Write with no colors
commenter -w -r -nc project/          # Recursive write with no colors
```

## Examples

```bash
# Preview comment removal
commenter example.go
commenter src/components/Button.tsx

# Apply changes to file
commenter --write example.go
commenter -w src/components/Button.tsx    # Short flag

# Use with npm runners
npx commenter -w src/utils/helper.ts
bunx commenter -nc build/output.sql       # Short flag for no-color

# Combine flags for efficiency
commenter -w -nc large-file.sql           # Write with no colors
```

## What gets removed

✅ **Removes:**

- Standalone comment lines (e.g., `// This is a comment`)
- Inline comments (e.g., `code(); // comment`)
- Multiple consecutive single-line comments

❌ **Preserves:**

- Multi-line comments (`/* ... */`)
- Comments inside string literals (`"string with // comment"`)
- Single-line comments inside multi-line comment blocks

## Adding Support for New File Types

Want to add support for a new programming language? It's easy! The tool uses a simple language definition system that you can extend.

### Step 1: Update Language Definitions

Edit the `const.go` file and add your new language to the `SupportedLanguages` map:

```go
var SupportedLanguages = map[string]Language{
    // ... existing languages ...

    "python": {
        Name:            "Python",
        Extensions:      []string{".py", ".pyw"},
        SingleLineStart: "#",
        MultiLineStart:  `"""`,
        MultiLineEnd:    `"""`,
    },
    "rust": {
        Name:            "Rust",
        Extensions:      []string{".rs"},
        SingleLineStart: "//",
        MultiLineStart:  "/*",
        MultiLineEnd:    "*/",
    },
    "shell": {
        Name:            "Shell Script",
        Extensions:      []string{".sh", ".bash", ".zsh"},
        SingleLineStart: "#",
        MultiLineStart:  "",  // No multi-line comments
        MultiLineEnd:    "",
    },
}
```

### Step 2: Language Definition Fields

Each language definition requires these fields:

| Field             | Description                  | Example                   | Required |
| ----------------- | ---------------------------- | ------------------------- | -------- |
| `Name`            | Human-readable language name | `"Python"`                | ✅       |
| `Extensions`      | File extensions (with dots)  | `[]string{".py", ".pyw"}` | ✅       |
| `SingleLineStart` | Single-line comment prefix   | `"#"`                     | ✅       |
| `MultiLineStart`  | Multi-line comment start     | `"""`                     | ❌       |
| `MultiLineEnd`    | Multi-line comment end       | `"""`                     | ❌       |

**Note**: If a language doesn't support multi-line comments, leave `MultiLineStart` and `MultiLineEnd` as empty strings (`""`).

### Step 3: Test Your Addition

1. **Create test files** with your new language extension:

   ```bash
   echo "# This is a comment" > test.py
   echo "print('Hello, World!')" >> test.py
   ```

2. **Test comment detection**:

   ```bash
   go run . test.py
   ```

3. **Verify the output** shows detected comments

4. **Run the test suite**:
   ```bash
   go test -v ./...
   ```

### Step 4: Add to Documentation

Update the **Supported Languages** table in this README:

```markdown
| Language | Extensions             | Single-line Comment |
| -------- | ---------------------- | ------------------- |
| Python   | `.py`, `.pyw`          | `#`                 |
| Rust     | `.rs`                  | `//`                |
| Shell    | `.sh`, `.bash`, `.zsh` | `#`                 |
```

### Examples of Common Languages

Here are some popular languages you might want to add:

<details>
<summary>📝 Click to see language definitions</summary>

```go
// C/C++
"c": {
    Name:            "C/C++",
    Extensions:      []string{".c", ".cpp", ".cc", ".cxx", ".h", ".hpp"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},

// Python
"python": {
    Name:            "Python",
    Extensions:      []string{".py", ".pyw"},
    SingleLineStart: "#",
    MultiLineStart:  `"""`,
    MultiLineEnd:    `"""`,
},

// Ruby
"ruby": {
    Name:            "Ruby",
    Extensions:      []string{".rb", ".rbw"},
    SingleLineStart: "#",
    MultiLineStart:  "=begin",
    MultiLineEnd:    "=end",
},

// PHP
"php": {
    Name:            "PHP",
    Extensions:      []string{".php", ".phtml"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},

// Java
"java": {
    Name:            "Java",
    Extensions:      []string{".java"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},

// C#
"csharp": {
    Name:            "C#",
    Extensions:      []string{".cs"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},

// Rust
"rust": {
    Name:            "Rust",
    Extensions:      []string{".rs"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},

// Swift
"swift": {
    Name:            "Swift",
    Extensions:      []string{".swift"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},

// Kotlin
"kotlin": {
    Name:            "Kotlin",
    Extensions:      []string{".kt", ".kts"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},

// Shell Scripts
"shell": {
    Name:            "Shell Script",
    Extensions:      []string{".sh", ".bash", ".zsh", ".fish"},
    SingleLineStart: "#",
    MultiLineStart:  "",
    MultiLineEnd:    "",
},
```

</details>

### Step 5: Contributing Back

Consider contributing your language additions back to the project:

1. **Fork the repository**
2. **Add your language definitions**
3. **Add test cases** in `test_examples/`
4. **Update documentation**
5. **Submit a pull request**

Your contributions help make this tool more useful for everyone! 🚀

## Project Structure

The codebase is organized into focused modules for maintainability:

- `main.go` - CLI parsing and application orchestration
- `const.go` - Language definitions and supported file types
- `processor.go` - Core comment removal logic and file processing
- `discovery.go` - File/directory discovery and batch processing
- `ui.go` - Terminal output, colors, and user interface

## Testing

### Comprehensive Test Suite

Run the full test suite with:

```bash
go test -v ./...                    # Run all unit tests
go test -bench=. -benchmem         # Run performance benchmarks
npm test                           # Test npm package functionality
```

### Test Coverage

- **Unit Tests**: Complete coverage of core functionality
- **Integration Tests**: File and directory processing
- **Benchmark Tests**: Performance testing for large files
- **CLI Tests**: End-to-end command-line interface testing
- **npm Package Tests**: Node.js integration and installation

### Test Categories

- `*_test.go` - Unit tests for each module
- `benchmark_test.go` - Performance benchmarks
- `test/npm-test.js` - npm package integration tests

## Building from Source

1. Clone the repository
2. Build the executable:
   ```bash
   make build
   # or manually:
   go build -o commenter
   ```
3. Install system-wide:
   ```bash
   make install
   ```

## Automated Releases

### GoReleaser Integration

- **Cross-platform builds**: Windows, macOS, Linux (amd64, arm64)
- **Multiple package managers**: Homebrew, Scoop, Winget, npm
- **Automated changelog generation**
- **GitHub Releases with binaries and checksums**

### CI/CD Pipeline

- **GitHub Actions** for automated testing and releases
- **Multi-OS testing** on Ubuntu, Windows, and macOS
- **Go versions 1.20 and 1.21** compatibility
- **Security scanning** with gosec
- **Code quality** with golangci-lint

### Release Process

1. Create and push a git tag: `git tag v1.x.x && git push origin v1.x.x`
2. GitHub Actions automatically:
   - Runs full test suite across platforms
   - Builds binaries for all supported platforms
   - Creates GitHub release with changelog
   - Updates package managers (Homebrew, Scoop, Winget)
   - Publishes to npm registry

### Package Manager Support

- **npm**: `npm install -g commenter` or `npx commenter`
- **Homebrew** (planned): `brew install your-username/tap/commenter`
- **Scoop** (planned): `scoop install commenter`
- **Winget** (planned): `winget install commenter`

## Example Output

**Single File:**

```
📁 File: example.go (Go)
Original lines: 27
Comments removed: 7
Remaining lines: 15

Removed comments:
  Line 5: // This is a standalone comment
  Line 8: fmt.Println("Hello") // Inline comment
  ...

Run with --write to apply changes to the file.
```

**Directory/Batch Processing:**

```
Batch Processing Summary:
Files processed: 15
Total comments removed: 89
Total lines processed: 1,247
Files written successfully: 15

Run with --write to apply changes to all files.
```

## Why "commenter" instead of "cr"?

The original name `cr` conflicts with an existing Unix command that converts text files between Unix and DOS line endings. To avoid this conflict and ensure compatibility across all systems, we renamed the tool to `commenter`.

## CLI Reference

| Short Flag | Long Flag     | Description                                   |
| ---------- | ------------- | --------------------------------------------- |
| `-w`       | `--write`     | Write changes to file (default: preview only) |
| `-r`       | `--recursive` | Process directories recursively               |
| `-nc`      | `--no-color`  | Disable colored output                        |
| `-h`       | `--help`      | Show detailed help message                    |
| `-v`       | `--version`   | Show version information                      |

## Advanced Usage

```bash
# Process entire project recursively
commenter -r .                           # Preview all files
commenter -w -r .                        # Process and save all files

# Process specific directories
commenter -r src/                        # Just the src directory
commenter -w -r src/ tests/              # Multiple directories (requires multiple commands)

# Legacy batch processing (still works)
find . -name "*.ts" -exec commenter -w {} \;

# Use in CI/CD pipelines with no color output
commenter -w -r -nc src/                 # Process directory tree
commenter --write --recursive --no-color src/  # Long flags (explicit)

# Combine with other tools
commenter -w -r src/ && npm run format

# Get help quickly
commenter -h                             # Show help
```
