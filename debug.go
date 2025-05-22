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

// parseDebugEnv parses the DEBUG environment variable to determine which modules to log
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

// Debug returns a function that logs debug messages for the specified module
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

// IsEnabled checks if debugging is enabled for a module
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

// ReloadDebugSettings allows reloading DEBUG environment variable at runtime
func ReloadDebugSettings() {
	debugMu.Lock()
	isInitialized = false
	debugMu.Unlock()
	parseDebugEnv()
}

// PrefixWriter is a writer that adds a prefix to each line written
type PrefixWriter struct {
	Prefix  string
	Ignores []string
}

// Write implements the io.Writer interface
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
