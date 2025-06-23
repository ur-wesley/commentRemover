# Extending Comment Remover

This guide explains how to add support for new programming languages to the Comment Remover tool.

## Overview

The tool uses a simple language definition system in `const.go`. Each language is defined with:

- File extensions to recognize
- Single-line comment syntax
- Multi-line comment syntax (if supported)

## Adding a New Language

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

Update the **Supported Languages** table in `README.md`:

```markdown
| Language | Extensions             | Single-line Comment |
| -------- | ---------------------- | ------------------- |
| Python   | `.py`, `.pyw`          | `#`                 |
| Rust     | `.rs`                  | `//`                |
| Shell    | `.sh`, `.bash`, `.zsh` | `#`                 |
```

## Common Language Examples

Here are some popular languages you might want to add:

### C/C++

```go
"c": {
    Name:            "C/C++",
    Extensions:      []string{".c", ".cpp", ".cc", ".cxx", ".h", ".hpp"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},
```

### Python

```go
"python": {
    Name:            "Python",
    Extensions:      []string{".py", ".pyw"},
    SingleLineStart: "#",
    MultiLineStart:  `"""`,
    MultiLineEnd:    `"""`,
},
```

### Rust

```go
"rust": {
    Name:            "Rust",
    Extensions:      []string{".rs"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},
```

### Shell Scripts

```go
"shell": {
    Name:            "Shell Script",
    Extensions:      []string{".sh", ".bash", ".zsh", ".fish"},
    SingleLineStart: "#",
    MultiLineStart:  "",  // No multi-line comments
    MultiLineEnd:    "",
},
```

### PHP

```go
"php": {
    Name:            "PHP",
    Extensions:      []string{".php", ".phtml"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},
```

### Java

```go
"java": {
    Name:            "Java",
    Extensions:      []string{".java"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},
```

### C#

```go
"csharp": {
    Name:            "C#",
    Extensions:      []string{".cs"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},
```

### Swift

```go
"swift": {
    Name:            "Swift",
    Extensions:      []string{".swift"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},
```

### Kotlin

```go
"kotlin": {
    Name:            "Kotlin",
    Extensions:      []string{".kt", ".kts"},
    SingleLineStart: "//",
    MultiLineStart:  "/*",
    MultiLineEnd:    "*/",
},
```

## Testing Your Changes

### Create Comprehensive Tests

Add test cases to `basic_test.go`:

```go
func TestPythonCommentRemoval(t *testing.T) {
    content := `# This is a comment
print("Hello, World!")  # Inline comment
"""
Multi-line comment
# This should NOT be removed
End of multi-line comment
"""
x = 42  # Another comment`

    // ... test implementation
}
```

### Test Edge Cases

Consider these scenarios:

- Comments inside strings
- Nested multi-line comments
- Comments at the end of files
- Empty lines with comments
- Special characters in comments

## Contributing Back

If you add support for a new language, consider contributing it back to the project:

1. **Fork the repository**
2. **Add your language definitions**
3. **Add test cases** in `basic_test.go`
4. **Update documentation** in `README.md`
5. **Submit a pull request**

Your contributions help make this tool more useful for everyone! 🚀

## Troubleshooting

### Common Issues

1. **Comments not being detected**: Check that `SingleLineStart` matches exactly
2. **Multi-line comments being removed**: Verify `MultiLineStart` and `MultiLineEnd` are correct
3. **Comments inside strings being removed**: The tool should handle this automatically
4. **File extensions not recognized**: Ensure extensions are in the `Extensions` slice

### Debug Tips

- Use `go run . --help` to see supported file types
- Test with a simple file first
- Check the test output for detailed information
- Use `--no-color` flag for cleaner output in CI environments
