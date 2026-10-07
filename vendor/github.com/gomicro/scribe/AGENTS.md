# Agent Instructions

## Overview

`scribe` is a Go library for structured, indented CLI output. It formats output as a tree of
labeled describe blocks and leaf print/error lines, with each nesting level indented by two
spaces. All visual styling is delegated to a caller-supplied `Theme`, so the same structured
output can be rendered plain or with ANSI color without changing call sites.

## Tech Stack

- Go 1.23.3
- Test assertions: `github.com/alecthomas/assert`
- Mock writer for tests: `github.com/gomicro/penname`

## Build & Validate

```sh
go fmt ./...
go vet ./...
go build ./...
go test ./...
```

Run these four in order after any Go change. There is no separate linter config.

## Repository Layout

| Path | Purpose |
|------|---------|
| `scribe.go` | `Scribe` struct; implements `Scriber` |
| `scriber.go` | `Scriber` interface definition |
| `theme.go` | `Theme` struct, `DefaultTheme()`, sentinel errors, `ValidateTheme` |
| `scribe_test.go` | Tests: describe nesting, full output tree, themed output, child scribers |
| `color/` | ANSI escape color helpers (`FgXxx` constants + bold `XxxFg` wrapping fns) |
| `vendor/` | Vendored dependencies — never edit directly; use `go mod vendor` |

## Architecture Notes

- `Scriber` is the public interface; callers only hold `Scriber`, never `*Scribe`.
- `NewScribe(writer io.Writer, theme *Theme) (Scriber, error)` is the sole constructor; it
  validates the theme and returns an error if any decorator is nil.
- `Child(desc string) Scriber` creates a buffered child scriber. All output is held in memory
  until `EndDescribe()` is called, at which point the entire block is written atomically to the
  parent's writer. This allows concurrent goroutines to each own a child and produce non-interleaved
  output.
- `color/` is a standalone utility package — it has no import relationship with the root package.
  Callers compose it with `Theme` decorator functions as they see fit.

### `Scriber` interface

```go
BeginDescribe(desc string)          // open a named section; increments indent
EndDescribe()                       // close section; decrements indent (flushes if child)
Child(desc string) Scriber          // buffered child for concurrent output
Print(str string)                   // one indented line via Theme.Print
Printf(format string, args ...any)  // formatted Print
PrintLines(r io.Reader)             // scan reader line-by-line; emit each via Theme.Print
Error(err error)                    // one indented line via Theme.Error
Errorf(format string, args ...any)  // formatted Error
```

### `Theme`

```go
type Theme struct {
    Describe func(string) string
    Print    func(string) string
    Error    func(error) string
}
```

All three fields are required — `ValidateTheme` returns a sentinel error for any nil field.
`DefaultTheme()` returns a theme with no-op decorators (strings pass through unchanged).

## Key Conventions

- Wrap errors: `fmt.Errorf("methodName: operation: %w", err)` — lowercase, no period.
- Never use the blank identifier `_` to discard errors or return values.
- `t.Parallel()` at the top of every test function and subtest.
- Tests use `penname.New()` as the mock `io.Writer`; assert with `github.com/alecthomas/assert`.
- `vendor/` is read-only — only `go mod vendor` may modify it.
- `color/` helper functions apply **bold + color** (`\x1b[1;XXm`); use the raw `FgXxx` constants
  when non-bold output is needed.
