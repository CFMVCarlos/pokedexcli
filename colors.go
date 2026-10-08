package main

import (
	"os"
)

// ANSI escape code constants for terminal styling.
const (
	ansiReset   = "\033[0m"
	ansiBold    = "\033[1m"
	ansiRed     = "\033[31m"
	ansiGreen   = "\033[32m"
	ansiYellow  = "\033[33m"
	ansiBlue    = "\033[34m"
	ansiMagenta = "\033[35m"
	ansiCyan    = "\033[36m"
)

// colorEnabled controls whether ANSI escape codes are emitted.
// Defaults to true, but can be disabled via NO_COLOR environment variable.
var colorEnabled = os.Getenv("NO_COLOR") == ""

// SetColorEnabled toggles ANSI terminal color formatting.
func SetColorEnabled(enabled bool) {
	colorEnabled = enabled
}

// colorize wraps text with the given ANSI color escape sequence if colors are enabled.
func colorize(colorCode, text string) string {
	if !colorEnabled {
		return text
	}
	return colorCode + text + ansiReset
}

// colorGreen returns green-styled text for success messages.
func colorGreen(text string) string {
	return colorize(ansiGreen, text)
}

// colorRed returns red-styled text for errors and warning messages.
func colorRed(text string) string {
	return colorize(ansiRed, text)
}

// colorYellow returns yellow-styled text for escapes and notifications.
func colorYellow(text string) string {
	return colorize(ansiYellow, text)
}

// colorCyan returns cyan-styled text for prompts and titles.
func colorCyan(text string) string {
	return colorize(ansiCyan, text)
}

// colorType returns the elemental type formatted with a type-specific ANSI color.
func colorType(typeName string) string {
	switch typeName {
	case "fire":
		return colorize(ansiRed, typeName)
	case "water":
		return colorize(ansiBlue, typeName)
	case "grass", "bug":
		return colorize(ansiGreen, typeName)
	case "electric":
		return colorize(ansiYellow, typeName)
	case "ice":
		return colorize(ansiCyan, typeName)
	case "poison", "psychic", "ghost":
		return colorize(ansiMagenta, typeName)
	default:
		return colorize(ansiReset, typeName)
	}
}
