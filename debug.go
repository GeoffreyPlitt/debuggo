// Package debuggo provides a lightweight debugging utility for Go applications
// inspired by the Node.js debug package. It allows for environment variable
// controlled debugging with namespace support.
//
// # Overview
//
// debuggo makes it easy to add configurable debug statements to your Go applications.
// Debug output is controlled via the DEBUG environment variable, allowing you to
// enable or disable specific debugging components without changing code or recompiling.
//
// # Key Features
//
//   - Namespace support with hierarchical components (e.g., "app:server:http")
//   - Selective enabling/disabling of debug components
//   - Wildcard support for enabling groups of related debug components
//   - Negation support to exclude specific components
//   - Runtime reconfiguration of debug settings
//
// # Basic Usage
//
//	var debug = debuggo.Debug("myapp:component")
//	debug("Processing item %s", item.ID)
//
// # Environment Variable Control
//
//	DEBUG=* # Enable all debug messages
//	DEBUG=myapp:* # Enable all myapp namespace messages
//	DEBUG=*,!verbose # Enable all except verbose namespace
//	DEBUG=app:*,!app:db # Enable all app components except database
//
// # Advanced Usage
//
// See the examples directory for more detailed usage examples.
package debuggo

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	debugNamespaces map[string]bool
	negatedModules  map[string]bool
	wildcardEnabled bool
	debugMu         sync.RWMutex
	isInitialized   bool
)

func init() {
	parseDebugEnv()
}

// parseDebugEnv parses the DEBUG environment variable to determine which modules to log.
// Format: DEBUG=namespace1,namespace2:*,!namespace3
// - Use comma to separate multiple namespaces
// - Use * as wildcard for all namespaces
// - Prefix with ! to negate a namespace
// - Use colon (:) for hierarchical namespaces
func parseDebugEnv() {
	debugMu.Lock()
	defer debugMu.Unlock()

	if isInitialized {
		return
	}

	// Reset state
	debugNamespaces = make(map[string]bool)
	negatedModules = make(map[string]bool)
	wildcardEnabled = false

	debugValue := os.Getenv("DEBUG")

	if debugValue == "" {
		// No DEBUG env var set
		isInitialized = true
		return
	}

	// Parse comma-separated namespaces
	namespaces := strings.Split(debugValue, ",")
	for _, ns := range namespaces {
		ns = strings.TrimSpace(ns)
		if ns == "" {
			continue
		}

		// Support negation with ! prefix
		if strings.HasPrefix(ns, "!") {
			trimmedNS := ns[1:]
			negatedModules[trimmedNS] = true
			debugNamespaces[trimmedNS] = false
		} else if ns == "*" {
			// Global wildcard
			wildcardEnabled = true
		} else {
			// Normal namespace
			debugNamespaces[ns] = true
		}
	}

	isInitialized = true
}

// Debug returns a function that logs debug messages for the specified module.
// The returned function mimics fmt.Printf, but only outputs when the module
// is enabled via the DEBUG environment variable.
//
// The debug function will:
//   - Check if the module is enabled based on the DEBUG environment variable
//   - Add a timestamp and module prefix to each message
//   - Output to stderr (for easy redirection)
//
// Debug messages are printed with the format:
//
//	15:04:05.000 module_name message
//
// Example:
//
//	debug := Debug("app:server")
//	debug("Server starting on port %d", port)
//
// Output:
//
//	12:34:56.789 app:server Server starting on port 8080
func Debug(module string) func(format string, args ...interface{}) {
	return func(format string, args ...interface{}) {
		// We need to ensure we're checking the same condition as IsEnabled
		debugMu.RLock()
		enabled := checkEnabled(module)
		debugMu.RUnlock()

		if !enabled {
			return
		}

		// Get timestamp
		timestamp := time.Now().Format("15:04:05.000")

		// Format message
		message := fmt.Sprintf(format, args...)

		// Print with timestamp and module name
		fmt.Fprintf(os.Stderr, "%s %s %s\n", timestamp, module, message)
	}
}

// IsEnabled checks if debugging is enabled for a module.
// This is useful for wrapping expensive debug operations that should only
// be executed when debugging is enabled.
//
// This function reads the current debug settings from memory, so it's efficient
// to call frequently. The debug settings are parsed from the DEBUG environment
// variable when the package is initialized or when ReloadDebugSettings is called.
//
// Example:
//
//	if IsEnabled("app:metrics") {
//	    // Compute expensive debug data
//	    metrics := calculateDetailedMetrics()
//	    debug("System metrics: %+v", metrics)
//	}
func IsEnabled(module string) bool {
	debugMu.RLock()
	defer debugMu.RUnlock()
	return checkEnabled(module)
}

// checkEnabled is the core function to check if a module is enabled
// This must be called with the lock held
func checkEnabled(module string) bool {
	// First check if module is explicitly negated
	if isNegated(module) {
		return false
	}

	// Then check if wildcard is enabled (enabling everything not explicitly negated)
	if wildcardEnabled {
		return true
	}

	// Check if this specific module is directly enabled
	if debugNamespaces[module] {
		return true
	}

	// Check for wildcard namespace match
	return isEnabledByWildcard(module)
}

// isNegated checks if a module is explicitly negated
// This must be called with the lock held
func isNegated(module string) bool {
	// Direct negation
	if negatedModules[module] {
		return true
	}

	// Check if parent namespace is negated with wildcard
	parts := strings.Split(module, ":")
	for i := 1; i <= len(parts); i++ {
		prefix := strings.Join(parts[:i], ":")
		if negatedModules[prefix] || negatedModules[prefix+"*"] || negatedModules[prefix+":*"] {
			return true
		}
	}

	return false
}

// isEnabledByWildcard checks if a module is enabled via wildcard namespace
// This must be called with the lock held
func isEnabledByWildcard(module string) bool {
	parts := strings.Split(module, ":")

	// Try increasingly specific namespace patterns
	for i := 1; i < len(parts); i++ {
		ns := strings.Join(parts[:i], ":")

		// Check for pattern like "app:*" that would enable "app:server"
		if debugNamespaces[ns+":*"] {
			return true
		}

		// Also check for pattern like "app*" (although less common)
		if debugNamespaces[ns+"*"] {
			return true
		}
	}

	return false
}

// ReloadDebugSettings allows reloading DEBUG environment variable at runtime.
// This is useful for changing debug settings without restarting the application.
//
// You typically call this after changing the DEBUG environment variable with os.Setenv().
// The new settings will take effect immediately for all subsequent debug calls.
//
// Example:
//
//	// Initially run with minimal debugging
//	os.Setenv("DEBUG", "app:errors")
//
//	// Later enable more verbose debugging dynamically
//	os.Setenv("DEBUG", "app:*,!app:metrics") // All app components except metrics
//	debuggo.ReloadDebugSettings()
func ReloadDebugSettings() {
	debugMu.Lock()
	isInitialized = false
	debugMu.Unlock()
	parseDebugEnv()
}

// PrefixWriter is a writer that adds a prefix to each line written.
// It can also be configured to ignore certain phrases.
// Implements io.Writer interface for integration with standard libraries.
//
// This can be useful for:
//   - Redirecting standard library log output to include a debug prefix
//   - Filtering out unwanted messages
//   - Integrating with libraries that expect an io.Writer
//
// Example:
//
//	// Redirect standard library log output
//	log.SetOutput(&debuggo.PrefixWriter{Prefix: "app:log"})
//
//	// Filter out health check logs
//	logger := &debuggo.PrefixWriter{
//	    Prefix: "app:api",
//	    Ignores: []string{"/health", "/ping"},
//	}
type PrefixWriter struct {
	// Prefix is added to the beginning of each line
	Prefix string
	// Ignores is a list of phrases that will cause the line to be skipped if found
	Ignores []string
}

// Write implements the io.Writer interface.
// It adds a prefix to each line and filters out lines containing ignored phrases.
//
// The method:
//   - Checks if the text contains any phrases listed in Ignores
//   - If not, prepends the Prefix to the text and writes to os.Stderr
//   - Always returns the original input length to satisfy the io.Writer contract
func (pw *PrefixWriter) Write(p []byte) (n int, err error) {
	text := string(p)

	// Check if text contains any ignored phrases
	if pw.Ignores != nil {
		for _, ignore := range pw.Ignores {
			if strings.Contains(text, ignore) {
				return len(p), nil
			}
		}
	}

	fmt.Fprint(os.Stderr, pw.Prefix+" "+text)
	return len(p), nil
}
