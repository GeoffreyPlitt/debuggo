package main

import (
	"fmt"
	"os"
	"time"

	"github.com/GeoffreyPlitt/debuggo"
)

var (
	debugServer     = debuggo.Debug("app:server")
	debugDatabase   = debuggo.Debug("app:database")
	debugHttpServer = debuggo.Debug("app:server:http")
	debugWsServer   = debuggo.Debug("app:server:websocket")
	debugSecurity   = debuggo.Debug("app:security")
)

func main() {
	fmt.Println("Starting application with DEBUG=" + os.Getenv("DEBUG"))
	fmt.Println("Try running with different DEBUG settings:")
	fmt.Println("  DEBUG=app:* ./advanced")
	fmt.Println("  DEBUG=app:server:* ./advanced")
	fmt.Println("  DEBUG=*,!app:server:websocket ./advanced")
	fmt.Println("")

	// Simulate application startup
	debugServer("Server initializing")
	debugHttpServer("HTTP server starting on port 8080")
	debugWsServer("WebSocket server starting on port 8081")
	debugDatabase("Connecting to database")

	// Simulate some application operations
	simulateHttpRequest("/api/users")
	simulateHttpRequest("/api/products")
	simulateWebSocketMessage("user-connected")

	// Demonstrate runtime reconfiguration
	fmt.Println("\n--- Changing debug configuration at runtime ---")

	// Change debug settings at runtime
	// For example, disable websocket debugging but enable database debugging
	if os.Getenv("DEBUG") != "" {
		fmt.Println("Changing DEBUG from '" + os.Getenv("DEBUG") + "' to '*,!app:server:websocket,app:database'")
		os.Setenv("DEBUG", "*,!app:server:websocket,app:database")
		debuggo.ReloadDebugSettings()
	} else {
		fmt.Println("DEBUG was not set. Setting to 'app:database'")
		os.Setenv("DEBUG", "app:database")
		debuggo.ReloadDebugSettings()
	}

	// Now debug messages are filtered according to new settings
	simulateHttpRequest("/api/settings")
	simulateWebSocketMessage("message-received") // This won't show if using the new settings

	// Database operations will now be visible even if they weren't before
	debugDatabase("Executing complex query")
	debugDatabase("Query completed in 25ms")

	debugSecurity("Security audit completed")
}

func simulateHttpRequest(path string) {
	debugHttpServer("Received HTTP request: %s", path)
	time.Sleep(50 * time.Millisecond)
	debugHttpServer("HTTP request completed: %s", path)
}

func simulateWebSocketMessage(message string) {
	debugWsServer("WebSocket message received: %s", message)
	time.Sleep(30 * time.Millisecond)
	debugWsServer("WebSocket message processed: %s", message)
}
