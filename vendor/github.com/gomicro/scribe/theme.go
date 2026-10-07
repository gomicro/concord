package scribe

import "fmt"

// Theme controls how each category of output line is formatted before it is
// written. All three fields are required; NewScribe returns an error if any
// is nil. Use NoopDecorator and NoopErrDecorator for fields that should pass
// the string through unchanged.
type Theme struct {
	// Describe formats the label emitted by BeginDescribe and Child.
	Describe func(string) string
	// Print formats lines emitted by Print, Printf, and PrintLines.
	Print func(string) string
	// Error formats lines emitted by Error and Errorf.
	Error func(error) string
}

// DefaultTheme returns a Theme with no-op decorators.
func DefaultTheme() *Theme {
	return &Theme{
		Describe: NoopDecorator,
		Print:    NoopDecorator,
		Error:    NoopErrDecorator,
	}
}

// Sentinel errors returned by ValidateTheme for missing Theme decorator fields.
var (
	ErrThemeDescribeMissing = fmt.Errorf("theme missing describe decorator")
	ErrThemePrintMissing    = fmt.Errorf("theme missing print decorator")
	ErrThemeErrorMissing    = fmt.Errorf("theme missing error decorator")
)

// NoopDecorator is a no-op for Theme.Describe and Theme.Print fields.
func NoopDecorator(s string) string {
	return s
}

// NoopErrDecorator is a no-op for the Theme.Error field.
func NoopErrDecorator(err error) string {
	return err.Error()
}

// ValidateTheme checks that all three decorator fields of theme are non-nil.
// It returns a wrapped sentinel error for the first nil field found.
func ValidateTheme(theme *Theme) error {
	if theme.Describe == nil {
		return fmt.Errorf("theme validation: %w", ErrThemeDescribeMissing)
	}

	if theme.Print == nil {
		return fmt.Errorf("theme validation: %w", ErrThemePrintMissing)
	}

	if theme.Error == nil {
		return fmt.Errorf("theme validation: %w", ErrThemeErrorMissing)
	}

	return nil
}
