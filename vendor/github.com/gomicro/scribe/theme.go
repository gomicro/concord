package scribe

type Theme struct {
	Describe func(string) string
	Print    func(string) string
}

var DefaultTheme = &Theme{
	Describe: NoopDecorator,
	Print:    NoopDecorator,
}

func NoopDecorator(s string) string {
	return s
}
