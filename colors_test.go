package main

import (
	"strings"
	"testing"
)

func TestColorize(t *testing.T) {
	// Ensure color is enabled for test
	SetColorEnabled(true)
	defer SetColorEnabled(false)

	text := "pikachu"
	green := colorGreen(text)
	if !strings.Contains(green, ansiGreen) {
		t.Errorf("expected green ANSI code in %q", green)
	}

	yellow := colorYellow(text)
	if !strings.Contains(yellow, ansiYellow) {
		t.Errorf("expected yellow ANSI code in %q", yellow)
	}

	red := colorRed(text)
	if !strings.Contains(red, ansiRed) {
		t.Errorf("expected red ANSI code in %q", red)
	}

	cyan := colorCyan(text)
	if !strings.Contains(cyan, ansiCyan) {
		t.Errorf("expected cyan ANSI code in %q", cyan)
	}

	// Test disabling colors
	SetColorEnabled(false)
	plain := colorGreen("plain text")
	if plain != "plain text" {
		t.Errorf("expected plain text when color is disabled, got %q", plain)
	}
}

func TestColorType(t *testing.T) {
	SetColorEnabled(true)
	defer SetColorEnabled(false)

	types := []struct {
		typeName     string
		expectedCode string
	}{
		{"normal", ansiWhite},
		{"fire", ansiBrightRed},
		{"water", ansiBrightBlue},
		{"grass", ansiBrightGreen},
		{"electric", ansiBrightYellow},
		{"ice", ansiBrightCyan},
		{"fighting", ansiRed},
		{"poison", ansiMagenta},
		{"ground", ansiYellow},
		{"flying", ansiCyan},
		{"psychic", ansiBrightMagenta},
		{"bug", ansiGreen},
		{"rock", ansiYellow},
		{"ghost", ansiMagenta},
		{"dragon", ansiBlue},
		{"steel", ansiBrightWhite},
		{"dark", ansiBrightBlack},
		{"fairy", ansiBrightMagenta},
		{"unknown-custom", ansiReset},
	}

	for _, tc := range types {
		colored := colorType(tc.typeName)
		if !strings.Contains(colored, tc.expectedCode) {
			t.Errorf("type %q expected code %q, got %q", tc.typeName, tc.expectedCode, colored)
		}
	}
}
