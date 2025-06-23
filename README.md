# Comment Remover (commenter)

A performant CLI tool that safely removes single-line comments from source code files.

## Features

- **Safe comment removal**: Only removes single-line comments that are not part of multi-line comments or inside string literals
- **Multiple language support**: TypeScript/JavaScript, Go, SQL, and JSON
- **Performance optimized**: Fast file processing with minimal memory usage
- **Preview mode**: See what would be removed before making changes
- **Smart file filtering**: Respects `.gitignore` and `.commenterignore` files
- **Flexible exclusion**: Use `--exclude` flag for runtime pattern exclusion

## Supported Languages

| Language              | Extensions                   | Single-line Comment |
| --------------------- | ---------------------------- | ------------------- |
| TypeScript/JavaScript | `.ts`, `.tsx`, `.js`, `.jsx` | `//`                |
| Go                    | `.go`                        | `//`                |
| SQL                   | `.sql`                       | `--`                |
| JSON                  | `.json`                      | `//`                |
| PHP                   | `.php`, `.phtml`             | `//`                |
| C#                    | `.cs`                        | `//`                |

## Installation

### Via npm/bun (recommended)

The tool automatically downloads the latest pre-built binary from GitHub releases for your platform.

```bash
# Install globally with npm
npm install -g @ur-wesley/commenter

# Or use with npx (no installation required - auto-downloads binary)
npx @ur-wesley/commenter <file/path>

# Or use with bun
bun add -g @ur-wesley/commenter
bunx @ur-wesley/commenter <file/path>  # Auto-downloads binary

# Or install locally and use with bun run
bun add @ur-wesley/commenter
bun run commenter <file/path>
```

### Via Go (build from source)

```bash
go install github.com/ur-wesley/commentRemover@latest
```

### Manual Installation

1. Download the binary from [GitHub releases](https://github.com/ur-wesley/commentRemover/releases)
2. Add to your PATH
3. Or build from source: `go build -o commenter .`

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

# Remove single-line multi-line comments (e.g., /* comment */)
commenter --remove-single-multiline <file/path>
commenter -m <file/path>              # Short flag

# Exclude files with patterns
commenter -e "*test.go,*.min.js" src/  # Exclude test and minified files
commenter --exclude "*.spec.js" .      # Exclude spec files

# Disable colored output
commenter --no-color <file/path>
commenter -nc <file/path>             # Short flag

# Show help
commenter --help
commenter -h                          # Short flag
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
npx @ur-wesley/commenter -w src/utils/helper.ts
bunx @ur-wesley/commenter -nc build/output.sql       # Short flag for no-color

# Use with bun run (when installed locally)
bun run commenter -w src/utils/helper.ts
bun run commenter -nc build/output.sql

# Exclude test files and minified files
commenter -e "*test.go,*.min.js" src/
commenter -w -e "*.spec.js" project/

# Combine flags for efficiency
commenter -w -nc large-file.sql           # Write with no colors
commenter -m src/file.js                  # Remove single-line multi-line comments
```

## What gets removed

✅ **Removes:**

- Standalone comment lines (e.g., `// This is a comment`)
- Inline comments (e.g., `code(); // comment`)
- Multiple consecutive single-line comments
- Single-line multi-line comments (e.g., `/* comment */`)

❌ **Preserves:**

- Multi-line comments (`/* ... */`)
- Comments inside string literals (`"string with // comment"`)
- Single-line comments inside multi-line comment blocks

## File Filtering

The tool respects ignore files and patterns:

- **`.gitignore`**: Standard Git ignore patterns
- **`.commenterignore`**: Tool-specific ignore patterns (same syntax as `.gitignore`)
- **`--exclude` flag**: Runtime glob patterns (e.g., `--exclude "*test.go,*.min.js"`)

If both `.gitignore` and `.commenterignore` exist, both are respected (union of rules).

## Adding Support for New File Types

See [EXTENDING.md](EXTENDING.md) for detailed instructions on adding support for new programming languages.

## License

MIT License - see [LICENSE](LICENSE) file for details.
