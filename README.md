# Comment Remover (commenter)

A performant CLI tool written in Go that safely removes single-line comments from source code files.

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

### Via npm/bun (recommended)

The tool automatically downloads the latest pre-built binary from GitHub releases for your platform.

```bash
# Install globally with npm
npm install -g @ur-wesley/commenter

# Or use with npx (no installation required)
npx @ur-wesley/commenter <file/path>

# Or use with bun
bun add -g @ur-wesley/commenter
bunx @ur-wesley/commenter <file/path>

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
npx @ur-wesley/commenter -w src/utils/helper.ts
bunx @ur-wesley/commenter -nc build/output.sql       # Short flag for no-color

# Use with bun run (when installed locally)
bun run commenter -w src/utils/helper.ts
bun run commenter -nc build/output.sql

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
    MultiLineStart:  `
```
