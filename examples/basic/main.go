// Basic example of using debuggo
//
// Run with:
//
//	DEBUG=* go run main.go         # Show all debug messages
//	DEBUG=db go run main.go        # Show only database messages
//	DEBUG=app,api go run main.go   # Show app and API messages
package main

import (
	"time"

	"github.com/GeoffreyPlitt/debuggo"
)

// Initialize debug loggers at package level for different components
var (
	debugApp = debuggo.Debug("app")
	debugDb  = debuggo.Debug("db")
	debugApi = debuggo.Debug("api")
)

func main() {
	// Application startup sequence with debug logs
	debugApp("Application starting")

	// Simulating database connection
	debugDb("Connecting to database")
	time.Sleep(100 * time.Millisecond)
	debugDb("Database connected")

	// Simulating API server startup
	debugApi("API server listening on port 8080")

	// Conditionally execute expensive debug operations
	if debuggo.IsEnabled("app") {
		// This code only runs when "app" debugging is enabled
		debugApp("Detailed startup information: %v", getDetailedInfo())
	}
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
