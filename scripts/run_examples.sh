#!/bin/bash
# ==========================================
# run_examples.sh - Run debuggo examples and capture their output
# This script executes the examples with different DEBUG settings
# and captures their output for documentation purposes.
# ==========================================
set -e  # Exit on error

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Banner
echo -e "${GREEN}==================================================${NC}"
echo -e "${GREEN}   debuggo Examples Runner${NC}"
echo -e "${GREEN}==================================================${NC}"

# Function to handle errors gracefully
handle_error() {
	echo -e "${YELLOW}Warning: An error occurred but continuing...${NC}"
}

# Ensure script continues even if some commands fail
trap handle_error ERR

# Clean up any previous output files
echo -e "${YELLOW}Cleaning up previous output files...${NC}"
rm -f ./*.stdout ./*.stderr ./examples/*.stdout ./examples/*.stderr 2>/dev/null || true
rm -f ./basic_*.out* ./advanced_*.out* 2>/dev/null || true

# Create temporary directory
TMPDIR=$(mktemp -d)
echo -e "${YELLOW}Using temporary directory: ${TMPDIR}${NC}"

# Clean up on exit
trap "echo 'Cleaning up temporary files...'; rm -rf $TMPDIR" EXIT

# Create a standalone version of the examples with self-contained debuggo package
mkdir -p $TMPDIR/basic
mkdir -p $TMPDIR/advanced

# Create go.mod files for both examples
cat > $TMPDIR/basic/go.mod << EOF
module basicexample

go 1.19
EOF

cat > $TMPDIR/advanced/go.mod << EOF
module advancedexample

go 1.19
EOF

# Create basic example with local imports
cat > $TMPDIR/basic/debug.go << EOF
package main

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

func IsEnabled(module string) bool {
	debugMu.RLock()
	defer debugMu.RUnlock()
	return checkEnabled(module)
}

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

func ReloadDebugSettings() {
	debugMu.Lock()
	isInitialized = false
	debugMu.Unlock()
	parseDebugEnv()
}

type PrefixWriter struct {
	// Prefix is added to the beginning of each line
	Prefix string
	// Ignores is a list of phrases that will cause the line to be skipped if found
	Ignores []string
}

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
EOF

cat > $TMPDIR/basic/main.go << EOF
// Basic example of using debuggo
//
// Run with:
//   DEBUG=* go run main.go         # Show all debug messages
//   DEBUG=db go run main.go        # Show only database messages
//   DEBUG=app,api go run main.go   # Show app and API messages
package main

import (
	"fmt"
	"time"
)

// Initialize debug loggers at package level for different components
var (
	debugApp = Debug("app")
	debugDb  = Debug("db")
	debugApi = Debug("api")
)

func main() {
	// Application startup sequence with debug logs
	fmt.Println("Starting application...")
	debugApp("Application starting")

	// Simulating database connection
	debugDb("Connecting to database")
	time.Sleep(100 * time.Millisecond)
	debugDb("Database connected")

	// Simulating API server startup
	debugApi("API server listening on port 8080")

	// Conditionally execute expensive debug operations
	if IsEnabled("app") {
		// This code only runs when "app" debugging is enabled
		debugApp("Detailed startup information: %v", getDetailedInfo())
	}
	
	fmt.Println("Application running. Debug messages were sent to stderr.")
}

// getDetailedInfo returns detailed information for debugging
// This simulates an expensive operation you might want to skip when debugging is disabled
func getDetailedInfo() map[string]string {
	return map[string]string{
		"version":     "1.0.0",
		"environment": "development",
		"buildDate":   time.Now().String(),
	}
}
EOF

# Create advanced example with local imports
cp $TMPDIR/basic/debug.go $TMPDIR/advanced/debug.go

cat > $TMPDIR/advanced/main.go << EOF
// Advanced example of using debuggo
//
// This example demonstrates:
// - Hierarchical namespaces (app:server:http, app:database, etc.)
// - Runtime reconfiguration of debug settings
// - Selective enabling/disabling of debug modules
//
// Run with different settings to see the effect:
//   DEBUG=app:* ./advanced                     # All app components
//   DEBUG=app:server:* ./advanced              # Only server components
//   DEBUG=*,!app:server:websocket ./advanced   # Everything except websocket
package main

import (
	"fmt"
	"os"
	"time"
)

// Define debug loggers with hierarchical namespaces
var (
	debugServer     = Debug("app:server")
	debugDatabase   = Debug("app:database")
	debugHttpServer = Debug("app:server:http")
	debugWsServer   = Debug("app:server:websocket")
	debugSecurity   = Debug("app:security")
)

func main() {
	// Display startup instructions
	printInstructions()

	// ===== PART 1: Application Startup =====
	debugServer("Server initializing")
	debugHttpServer("HTTP server starting on port 8080")
	debugWsServer("WebSocket server starting on port 8081")
	debugDatabase("Connecting to database")

	// ===== PART 2: Simulate Application Activity =====
	simulateHttpRequest("/api/users")
	simulateHttpRequest("/api/products")
	simulateWebSocketMessage("user-connected")

	// ===== PART 3: Runtime Reconfiguration =====
	fmt.Println("\n--- Changing debug configuration at runtime ---")
	reconfigureDebugSettings()

	// ===== PART 4: Post-Reconfiguration Activity =====
	// These operations will be logged or not based on new settings
	simulateHttpRequest("/api/settings")
	simulateWebSocketMessage("message-received") // May not show with new settings
	debugDatabase("Executing complex query")     // Should show with new settings
	debugDatabase("Query completed in 25ms")

	debugSecurity("Security audit completed")
}

// printInstructions displays guidance on how to run the example
func printInstructions() {
	fmt.Println("Starting application with DEBUG=" + os.Getenv("DEBUG"))
	fmt.Println("Try running with different DEBUG settings:")
	fmt.Println("  DEBUG=app:* ./advanced")
	fmt.Println("  DEBUG=app:server:* ./advanced")
	fmt.Println("  DEBUG=*,!app:server:websocket ./advanced")
	fmt.Println("")
}

// reconfigureDebugSettings demonstrates changing debug configuration at runtime
func reconfigureDebugSettings() {
	if os.Getenv("DEBUG") != "" {
		fmt.Println("Changing DEBUG from '" + os.Getenv("DEBUG") + "' to '*,!app:server:websocket,app:database'")
		os.Setenv("DEBUG", "*,!app:server:websocket,app:database")
	} else {
		fmt.Println("DEBUG was not set. Setting to 'app:database'")
		os.Setenv("DEBUG", "app:database")
	}
	
	// Apply the new settings
	ReloadDebugSettings()
}

// simulateHttpRequest simulates an HTTP request with debug logging
func simulateHttpRequest(path string) {
	debugHttpServer("Received HTTP request: %s", path)
	time.Sleep(50 * time.Millisecond)
	debugHttpServer("HTTP request completed: %s", path)
}

// simulateWebSocketMessage simulates a WebSocket message with debug logging
func simulateWebSocketMessage(message string) {
	debugWsServer("WebSocket message received: %s", message)
	time.Sleep(30 * time.Millisecond)
	debugWsServer("WebSocket message processed: %s", message)
}
EOF

# Output directory - ensure it exists with absolute path
OUTDIR=$(pwd)/example_outputs
mkdir -p "$OUTDIR"

echo -e "${YELLOW}Output will be saved to ${OUTDIR}${NC}"

# Run examples and capture output
echo -e "${YELLOW}=== Running basic example with DEBUG=* ===${NC}"
(cd $TMPDIR/basic && DEBUG="*" go run . > $OUTDIR/basic_all_output.stdout 2> $OUTDIR/basic_all_output.stderr)

echo ""
echo -e "${YELLOW}=== Running basic example with DEBUG=db ===${NC}"
(cd $TMPDIR/basic && DEBUG="db" go run . > $OUTDIR/basic_db_output.stdout 2> $OUTDIR/basic_db_output.stderr)

echo ""
echo -e "${YELLOW}=== Running advanced example with DEBUG=app:server:* ===${NC}"
(cd $TMPDIR/advanced && DEBUG="app:server:*" go run . > $OUTDIR/advanced_server_output.stdout 2> $OUTDIR/advanced_server_output.stderr)

echo ""
echo -e "${YELLOW}=== Running advanced example with DEBUG=* ===${NC}"
(cd $TMPDIR/advanced && DEBUG="*" go run . > $OUTDIR/advanced_all_output.stdout 2> $OUTDIR/advanced_all_output.stderr)

echo ""
echo -e "${GREEN}All outputs captured in ${OUTDIR}/${NC}"
echo -e "${GREEN}You can use these files to update the README.md or documentation.${NC}"

# Print a summary of what was captured
echo ""
echo -e "${GREEN}=== Basic Example (DEBUG=*) ===${NC}"
echo -e "${YELLOW}STDOUT:${NC}"
cat $OUTDIR/basic_all_output.stdout 2>/dev/null || echo "No stdout captured"
echo -e "${YELLOW}STDERR:${NC}"
cat $OUTDIR/basic_all_output.stderr 2>/dev/null || echo "No stderr captured"

echo ""
echo -e "${GREEN}=== Basic Example (DEBUG=db) ===${NC}"
echo -e "${YELLOW}STDOUT:${NC}"
cat $OUTDIR/basic_db_output.stdout 2>/dev/null || echo "No stdout captured"
echo -e "${YELLOW}STDERR:${NC}"
cat $OUTDIR/basic_db_output.stderr 2>/dev/null || echo "No stderr captured"

echo ""
echo -e "${GREEN}=== Advanced Example (DEBUG=app:server:*) ===${NC}"
echo -e "${YELLOW}STDOUT:${NC}"
cat $OUTDIR/advanced_server_output.stdout 2>/dev/null || echo "No stdout captured"
echo -e "${YELLOW}STDERR:${NC}"
cat $OUTDIR/advanced_server_output.stderr 2>/dev/null || echo "No stderr captured"

echo ""
echo -e "${GREEN}=== Advanced Example (DEBUG=*) ===${NC}"
echo -e "${YELLOW}STDOUT:${NC}"
cat $OUTDIR/advanced_all_output.stdout 2>/dev/null || echo "No stdout captured"
echo -e "${YELLOW}STDERR:${NC}"
cat $OUTDIR/advanced_all_output.stderr 2>/dev/null || echo "No stderr captured"

echo ""
echo -e "${GREEN}==================================================${NC}"
echo -e "${GREEN}   Script completed successfully${NC}"
echo -e "${GREEN}==================================================${NC}"

# Instructions for next steps
echo ""
echo -e "${YELLOW}To use these outputs in the README.md:${NC}"
echo "1. Copy the content from the output files"
echo "2. Update the README.md with the actual output examples"
echo ""
echo -e "${YELLOW}To clean up all generated files:${NC}"
echo "  rm -rf $OUTDIR"

# Always exit with success
exit 0 