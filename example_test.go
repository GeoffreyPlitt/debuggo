package debuggo_test

import (
	"os"

	"github.com/GeoffreyPlitt/debuggo"
)

// This example demonstrates basic usage of the debuggo package.
func Example() {
	// Set DEBUG environment variable to enable debugging
	os.Setenv("DEBUG", "app:*")
	debuggo.ReloadDebugSettings()

	// Create debug loggers for different components
	debugApp := debuggo.Debug("app:main")
	debugDB := debuggo.Debug("app:database")
	debugAPI := debuggo.Debug("app:api")

	// Log debug messages
	debugApp("Application starting")
	debugDB("Connected to database")

	// Conditionally execute expensive debug operations
	if debuggo.IsEnabled("app:api") {
		debugAPI("API server listening on %s", "localhost:8080")
	}

	// Output is sent to stderr by default
}

// This example demonstrates how to selectively enable debug components
func Example_selective() {
	// Only enable database debugging
	os.Setenv("DEBUG", "app:database")
	debuggo.ReloadDebugSettings()

	debug := debuggo.Debug("app:database")
	debug("This debug message will appear")

	// This won't be output
	other := debuggo.Debug("app:other")
	other("This message won't appear")

	// Output is sent to stderr
}

// This example demonstrates wildcard and negation features
func Example_wildcardAndNegation() {
	// Enable all debugging except for the metrics component
	os.Setenv("DEBUG", "*,!app:metrics")
	debuggo.ReloadDebugSettings()

	debugApp := debuggo.Debug("app:main")
	debugAPI := debuggo.Debug("app:api")
	debugMetrics := debuggo.Debug("app:metrics")

	debugApp("This will be logged")
	debugAPI("This will also be logged")
	debugMetrics("This will NOT be logged")

	// Output is sent to stderr
}
