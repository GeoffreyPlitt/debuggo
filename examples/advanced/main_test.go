package main

import (
	"os"
	"testing"

	"github.com/GeoffreyPlitt/debuggo"
)

// TestAdvancedExample ensures all functions are executed for coverage
func TestAdvancedExample(t *testing.T) {
	// Redirect stderr and stdout to /dev/null
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

	// Test scenarios that are guaranteed to hit all code paths
	debugScenarios := []string{
		"app:*",                   // All app components
		"app:server:*",            // Only server components
		"*,!app:server:websocket", // Everything except websocket
	}

	for _, setting := range debugScenarios {
		// Set DEBUG and run the main function
		os.Setenv("DEBUG", setting)
		debuggo.ReloadDebugSettings()

		// Call main which will execute everything
		main()

		// Ensure these specific functions are called for coverage
		reconfigureDebugSettings()
		simulateHttpRequest("/api/direct-test")
		simulateWebSocketMessage("direct-test")
	}
}
