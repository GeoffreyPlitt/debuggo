package main

import (
	"time"

	"github.com/GeoffreyPlitt/debuggo"
)

func main() {
	// Create debug loggers for different modules
	debugApp := debuggo.Debug("app")
	debugDb := debuggo.Debug("db")
	debugApi := debuggo.Debug("api")

	// Log some messages
	debugApp("Application starting")

	debugDb("Connecting to database")
	time.Sleep(100 * time.Millisecond)
	debugDb("Database connected")

	debugApi("API server listening on port 8080")

	// You can also check if debugging is enabled
	if debuggo.IsEnabled("app") {
		// Expensive debug operations can be wrapped in this check
		debugApp("Detailed startup information: %v", getDetailedInfo())
	}
}

func getDetailedInfo() map[string]string {
	return map[string]string{
		"version":     "1.0.0",
		"environment": "development",
		"buildDate":   time.Now().String(),
	}
}
