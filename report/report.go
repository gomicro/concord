package report

import "fmt"

const (
	colorRed    = "\033[1;31m"
	colorGreen  = "\033[1;32m"
	colorYellow = "\033[1;33m"
	colorBlue   = "\033[1;34m"
	colorPurple = "\033[1;35m"
	colorCyan   = "\033[1;36m"
	colorWhite  = "\033[1;37m"
	colorReset  = "\033[0m"
)

func PrintHeader(text string) {
	fmt.Printf("%s%s%s", colorBlue, text, colorReset)
}

func Println() {
	fmt.Println()
}

func PrintWarn(text string) {
	fmt.Printf("  %s%s%s", colorYellow, text, colorReset)
}

func PrintError(text string) {
	fmt.Printf("  %s%s%s", colorRed, text, colorReset)
}
