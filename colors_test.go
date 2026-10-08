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
	if !strings.Contains(green, colorGreenCode) && !strings.Contains(green, "\033[32m") {
		t.Errorf("expected green ANSI code in %q", green)
	}

	yellow := colorYellow(text)
	if !strings.Contains(yellow, "\033[33m") {
		t.Errorf("expected yellow ANSI code in %q", yellow)
	}

	red := colorRed(text)
	if !strings.Contains(red, "\033[31m") {
		t.Errorf("expected red ANSI code in %q", red)
	}

	cyan := colorCyan(text)
	if !strings.Contains(cyan, "\033[36m") {
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
		{"fire", "\033[31m"},
		{"water", "\033[34m"},
		{"grass", "\033[32m"},
		{"electric", "\033[33m"},
		{"ice", "\033[36m"},
		{"poison", "\033[35m"},
		{"normal", "\033[0m"},
	}

	for _, tc := range types {
		colored := colorType(tc.typeName)
		if !strings.Contains(colored, tc.expectedCode) {
			t.Errorf("type %q expected code %q, got %q", tc.typeName, tc.expectedCode, colored)
		}
	}
}

const colorGreenCode = "\033[32m"
