// Package scribe provides structured, indented output for CLI tools. Output is
// organized as a tree of named sections and leaf lines, with each nesting level
// indented by two spaces. All visual styling is delegated to a caller-supplied
// Theme, so the same structured output can be rendered plain or with ANSI color
// without changing call sites.
package scribe

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"
)

type Scribe struct {
	writer io.Writer
	mu     *sync.Mutex
	level  int
	theme  *Theme
	parent io.Writer     // non-nil for buffered children
	buf    *bytes.Buffer // non-nil for buffered children
}

// NewScribe creates a new Scriber that writes to writer using the given theme.
// It returns an error if the theme fails validation (any decorator is nil).
func NewScribe(writer io.Writer, theme *Theme) (Scriber, error) {
	err := ValidateTheme(theme)
	if err != nil {
		return nil, fmt.Errorf("new scribe: %w", err)
	}

	return &Scribe{
		writer: writer,
		mu:     &sync.Mutex{},
		theme:  theme,
	}, nil
}

func (s *Scribe) BeginDescribe(desc string) {
	s.println()
	s.printt(s.theme.Describe(desc))
	s.level++
}

func (s *Scribe) EndDescribe() {
	s.level--
	if s.parent != nil {
		s.mu.Lock()
		fmt.Fprintf(s.parent, "%s", s.buf)
		s.mu.Unlock()
		s.buf.Reset()
		s.parent = nil
	}
}

// Child creates a new buffered Scriber labeled with desc at the current indent
// level. All output written to the child is held in memory until EndDescribe is
// called, at which point the entire block is written atomically to the parent's
// writer. This allows multiple goroutines to each own a child and produce
// grouped, non-interleaved output.
func (s *Scribe) Child(desc string) Scriber {
	buf := &bytes.Buffer{}
	child := &Scribe{
		writer: buf,
		mu:     s.mu,
		level:  s.level,
		theme:  s.theme,
		parent: s.writer,
		buf:    buf,
	}
	child.println()
	child.printt(child.theme.Describe(desc))
	child.level++
	return child
}

func (s *Scribe) Print(str string) {
	s.level++
	s.printt(s.theme.Print(str))
	s.level--
}

func (s *Scribe) Printf(format string, args ...any) {
	s.Print(fmt.Sprintf(format, args...))
}

func (s *Scribe) PrintLines(r io.Reader) {
	scanner := bufio.NewScanner(r)

	s.level++

	for scanner.Scan() {
		s.printt(s.theme.Print(scanner.Text()))
	}

	s.println()
	s.level--
}

func (s *Scribe) Error(err error) {
	s.level++
	s.printt(s.theme.Error(err))
	s.level--
}

func (s *Scribe) Errorf(format string, args ...any) {
	s.Error(fmt.Errorf(format, args...))
}

func (s *Scribe) print(str string) {
	fmt.Fprintf(s.writer, "%v\n", str)
}

func (s *Scribe) printt(str string) {
	s.print(fmt.Sprintf("%v%v", s.space(), str))
}

func (s *Scribe) println() {
	s.print("")
}

func (s *Scribe) space() string {
	return strings.Repeat(" ", s.level*2)
}
