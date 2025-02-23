package scribe

type Theme struct {
	Describe func(string) string
	Done     func(string) string
}

var DefaultTheme = &Theme{
	Describe: NoopDecorator,
	Done:     NoopDecorator,
}

func NoopDecorator(s string) string {
	return s
}
