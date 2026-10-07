package scribe

import "io"

// Scriber is the interface for structured, indented CLI output. Output is
// organized as a tree of named sections opened with BeginDescribe/EndDescribe
// (or Child) and leaf lines written with Print, PrintLines, and Error. Each
// nesting level is indented by two spaces. All visual styling is applied by
// the Theme supplied to NewScribe.
type Scriber interface {
	// BeginDescribe emits desc at the current indent and increases indent for
	// subsequent output.
	BeginDescribe(desc string)

	// EndDescribe decrements the indent. For child scribers it also flushes
	// buffered output atomically to the parent writer.
	EndDescribe()

	// Child creates a buffered child scriber at the current indent level.
	// Output is held until EndDescribe, then written atomically to the parent
	// writer — allowing goroutines to produce non-interleaved output.
	Child(desc string) Scriber

	// Print emits a line at the current indent, styled by Theme.Print.
	Print(str string)

	// Printf is the formatted variant of Print.
	Printf(format string, args ...any)

	// PrintLines scans r and emits each line at the current indent, styled by
	// Theme.Print.
	PrintLines(r io.Reader)

	// Error emits a line at the current indent, styled by Theme.Error.
	Error(err error)

	// Errorf is the formatted variant of Error.
	Errorf(format string, args ...any)
}
