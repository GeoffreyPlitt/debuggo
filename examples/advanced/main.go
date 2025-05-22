// Advanced example of using debuggo
//
// This example demonstrates:
// - Hierarchical namespaces (app:server:http, app:database, etc.)
// - Runtime reconfiguration of debug settings
// - Selective enabling/disabling of debug modules
//
// Run with different settings to see the effect:
//
//	DEBUG=app:* ./advanced                     # All app components
//	DEBUG=app:server:* ./advanced              # Only server components
//	DEBUG=*,!app:server:websocket ./advanced   # Everything except websocket
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/GeoffreyPlitt/debuggo"
)

// Define debug loggers with hierarchical namespaces
var (
	debugServer     = debuggo.Debug("app:server")
	debugDatabase   = debuggo.Debug("app:database")
	debugHttpServer = debuggo.Debug("app:server:http")
	debugWsServer   = debuggo.Debug("app:server:websocket")
	debugSecurity   = debuggo.Debug("app:security")
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
	fmt.Println("  DEBUG=\"app:*\" ./advanced")
	fmt.Println("  DEBUG=\"app:server:*\" ./advanced")
	fmt.Println("  DEBUG=\"*,!app:server:websocket\" ./advanced")
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
	debuggo.ReloadDebugSettings()
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
