package ui

import (
	"os"
)

// ANSI escape code constants for terminal styling.
const (
	ansiReset         = "\033[0m"
	ansiBold          = "\033[1m"
	ansiRed           = "\033[31m"
	ansiGreen         = "\033[32m"
	ansiYellow        = "\033[33m"
	ansiBlue          = "\033[34m"
	ansiMagenta       = "\033[35m"
	ansiCyan          = "\033[36m"
	ansiWhite         = "\033[37m"
	ansiBrightBlack   = "\033[90m" // Dark Gray
	ansiBrightRed     = "\033[91m"
	ansiBrightGreen   = "\033[92m"
	ansiBrightYellow  = "\033[93m"
	ansiBrightBlue    = "\033[94m"
	ansiBrightMagenta = "\033[95m"
	ansiBrightCyan    = "\033[96m"
	ansiBrightWhite   = "\033[97m"
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

// ColorGreen returns green-styled text for success messages.
func ColorGreen(text string) string {
	return colorize(ansiGreen, text)
}

// colorRed returns red-styled text for errors and warning messages.
func colorRed(text string) string {
	return colorize(ansiRed, text)
}

// ColorYellow returns yellow-styled text for escapes and notifications.
func ColorYellow(text string) string {
	return colorize(ansiYellow, text)
}

// ColorCyan returns cyan-styled text for prompts and titles.
func ColorCyan(text string) string {
	return colorize(ansiCyan, text)
}

// ColorType returns the elemental type formatted with an authentic, type-specific ANSI color
// for all 18 official Pokémon elemental types.
func ColorType(typeName string) string {
	switch typeName {
	case "normal":
		return colorize(ansiWhite, typeName)
	case "fire":
		return colorize(ansiBrightRed, typeName)
	case "water":
		return colorize(ansiBrightBlue, typeName)
	case "grass":
		return colorize(ansiBrightGreen, typeName)
	case "electric":
		return colorize(ansiBrightYellow, typeName)
	case "ice":
		return colorize(ansiBrightCyan, typeName)
	case "fighting":
		return colorize(ansiRed, typeName)
	case "poison":
		return colorize(ansiMagenta, typeName)
	case "ground":
		return colorize(ansiYellow, typeName)
	case "flying":
		return colorize(ansiCyan, typeName)
	case "psychic":
		return colorize(ansiBrightMagenta, typeName)
	case "bug":
		return colorize(ansiGreen, typeName)
	case "rock":
		return colorize(ansiYellow, typeName)
	case "ghost":
		return colorize(ansiMagenta, typeName)
	case "dragon":
		return colorize(ansiBlue, typeName)
	case "steel":
		return colorize(ansiBrightWhite, typeName)
	case "dark":
		return colorize(ansiBrightBlack, typeName)
	case "fairy":
		return colorize(ansiBrightMagenta, typeName)
	default:
		return colorize(ansiReset, typeName)
	}
}
