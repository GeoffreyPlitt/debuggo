package main

import (
	"os"
	"testing"

	"github.com/GeoffreyPlitt/debuggo"
)

// TestBasicExample ensures all example code is executed for coverage
func TestBasicExample(t *testing.T) {
	// Redirect stderr to /dev/null to avoid test output pollution
	origStderr := os.Stderr
	origStdout := os.Stdout
	devNull, _ := os.Open(os.DevNull)
	os.Stderr = devNull
	os.Stdout = devNull
	defer func() {
		os.Stderr = origStderr
		os.Stdout = origStdout
	}()

	// Save original DEBUG value
	originalDebug := os.Getenv("DEBUG")
	defer os.Setenv("DEBUG", originalDebug)

	// Test with all debug enabled
	os.Setenv("DEBUG", "*")
	debuggo.ReloadDebugSettings()

	// Call main directly - this will execute the complete example
	main()

	// Test with only DB debug enabled
	os.Setenv("DEBUG", "db")
	debuggo.ReloadDebugSettings()
	debugApp("This should not appear")
	debugDb("This should appear")

	// Test with no debug enabled
	os.Setenv("DEBUG", "")
	debuggo.ReloadDebugSettings()
	getDetailedInfo() // Call directly for coverage
}
