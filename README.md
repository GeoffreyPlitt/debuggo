# debuggo

[![Go Report Card](https://goreportcard.com/badge/github.com/GeoffreyPlitt/debuggo)](https://goreportcard.com/report/github.com/GeoffreyPlitt/debuggo)
[![GoDoc](https://godoc.org/github.com/GeoffreyPlitt/debuggo?status.svg)](https://godoc.org/github.com/GeoffreyPlitt/debuggo)
[![Build Status](https://github.com/GeoffreyPlitt/debuggo/workflows/Go/badge.svg)](https://github.com/GeoffreyPlitt/debuggo/actions)
[![codecov](https://codecov.io/gh/GeoffreyPlitt/debuggo/branch/main/graph/badge.svg)](https://codecov.io/gh/GeoffreyPlitt/debuggo)

A lightweight debugging utility for Go applications inspired by the Node.js [debug](https://www.npmjs.com/package/debug) package.

## Features

- **Environment variable control** - Enable/disable debug output via the `DEBUG` env var
- **Hierarchical namespaces** - Organize debug messages in namespaces like `app:server:http`
- **Zero overhead when disabled** - Debug statements incur virtually no cost when disabled
- **Runtime reconfiguration** - Change debug settings without restarting your application
- **Conditional debugging** - Skip expensive debug operations when not needed

## Installation

```bash
go get -u github.com/GeoffreyPlitt/debuggo
```

## Quick Start

```go
package main

import (
    "github.com/GeoffreyPlitt/debuggo"
)

// Create debug loggers for different components
var debug = debuggo.Debug("myapp")
var dbDebug = debuggo.Debug("myapp:database")

func main() {
    debug("Application starting up")
    
    if dbDebug("Connecting to database..."); dbDebug("Connected!") {
        // Debug messages are only printed when enabled
    }
    
    // Check if debugging is enabled before doing expensive operations
    if debuggo.IsEnabled("myapp:metrics") {
        // Generate expensive debug data only when needed
        debug("Memory usage: %v", collectMemoryMetrics())
    }
}
```

## Environment Variable Control

The `DEBUG` environment variable controls which debug messages are displayed:

```bash
# Enable all debug output
DEBUG=* go run main.go

# Enable specific modules
DEBUG=http,db go run main.go

# Enable hierarchical modules
DEBUG=myapp:* go run main.go  # Enables all 'myapp:' modules

# Enable everything except specific modules
DEBUG=*,!verbose go run main.go
```

## Advanced Usage

### Hierarchical Namespaces

Organize your debug loggers with namespaces to enable selective debugging:

```go
var (
    debugHttp = debuggo.Debug("app:http")       // HTTP server component
    debugWs   = debuggo.Debug("app:websocket")  // WebSocket component
    debugDb   = debuggo.Debug("app:database")   // Database component
)

// Enable only HTTP logs with: DEBUG=app:http
// Enable all components with: DEBUG=app:*
```

### Conditional Debugging

Skip expensive debug operations when debugging is disabled:

```go
if debuggo.IsEnabled("app:metrics") {
    debugMetrics("System stats: %v", collectDetailedMetrics())
}
```

### Runtime Reconfiguration

Change debug settings without restarting your application:

```go
// Change debug settings dynamically
os.Setenv("DEBUG", "newmodule,*:important")
debuggo.ReloadDebugSettings()
```

### Using PrefixWriter

Capture and prefix command output for easier debugging:

```go
cmd := exec.Command("some-program")

// Prefix stdout and stderr with custom identifiers
cmd.Stdout = &debuggo.PrefixWriter{Prefix: "CMD-OUT>"}
cmd.Stderr = &debuggo.PrefixWriter{
    Prefix:  "CMD-ERR>",
    Ignores: []string{"known warning to ignore"},
}
```

## Examples

The repository contains example applications demonstrating various features:

### Basic Example

See [examples/basic/main.go](examples/basic/main.go) for a simple usage example:

```go
// Basic example of using debuggo
package main

import (
    "fmt"
    "time"
    "github.com/GeoffreyPlitt/debuggo"
)

// Initialize debug loggers at package level
var (
    debugApp = debuggo.Debug("app")
    debugDb  = debuggo.Debug("db")
    debugApi = debuggo.Debug("api")
)

func main() {
    // Application startup
    fmt.Println("Starting application...")
    debugApp("Application starting")
    
    debugDb("Connecting to database")
    time.Sleep(100 * time.Millisecond)
    debugDb("Database connected")
    
    debugApi("API server listening on port 8080")
    
    if debuggo.IsEnabled("app") {
        // Only runs when "app" debugging is enabled
        debugApp("Detailed startup information: %v", getDetailedInfo())
    }
    
    fmt.Println("Application running.")
}
```

#### Output with DEBUG=*

```
$ DEBUG=* go run examples/basic/main.go
Starting application...
20:16:27.045 app Application starting
20:16:27.045 db Connecting to database
20:16:27.146 db Database connected
20:16:27.146 api API server listening on port 8080
20:16:27.146 app Detailed startup information: map[buildDate:2025-05-21 20:16:27.146 environment:development version:1.0.0]
Application running. Debug messages were sent to stderr.
```

#### Output with DEBUG=db

```
$ DEBUG=db go run examples/basic/main.go
Starting application...
20:16:27.331 db Connecting to database
20:16:27.432 db Database connected
Application running. Debug messages were sent to stderr.
```

### Advanced Example

See [examples/advanced/main.go](examples/advanced/main.go) for:
- Hierarchical namespaces
- Runtime reconfiguration
- Selective enabling/disabling

#### Output with DEBUG=app:server:*

```
$ DEBUG=app:server:* go run examples/advanced/main.go
Starting application with DEBUG=app:server:*
Try running with different DEBUG settings:
  DEBUG=app:* ./advanced
  DEBUG=app:server:* ./advanced
  DEBUG=*,!app:server:websocket ./advanced

20:16:27.648 app:server:http HTTP server starting on port 8080
20:16:27.648 app:server:websocket WebSocket server starting on port 8081
20:16:27.648 app:server:http Received HTTP request: /api/users
20:16:27.699 app:server:http HTTP request completed: /api/users
20:16:27.699 app:server:http Received HTTP request: /api/products
20:16:27.750 app:server:http HTTP request completed: /api/products
20:16:27.750 app:server:websocket WebSocket message received: user-connected
20:16:27.781 app:server:websocket WebSocket message processed: user-connected

--- Changing debug configuration at runtime ---
Changing DEBUG from 'app:server:*' to '*,!app:server:websocket,app:database'
20:16:27.781 app:server:http Received HTTP request: /api/settings
20:16:27.832 app:server:http HTTP request completed: /api/settings
20:16:27.863 app:database Executing complex query
20:16:27.864 app:database Query completed in 25ms
20:16:27.864 app:security Security audit completed
```

#### Output with DEBUG=*

```
$ DEBUG=* go run examples/advanced/main.go
Starting application with DEBUG=*
Try running with different DEBUG settings:
  DEBUG=app:* ./advanced
  DEBUG=app:server:* ./advanced
  DEBUG=*,!app:server:websocket ./advanced

20:16:28.037 app:server Server initializing
20:16:28.037 app:server:http HTTP server starting on port 8080
20:16:28.037 app:server:websocket WebSocket server starting on port 8081
20:16:28.037 app:database Connecting to database
20:16:28.037 app:server:http Received HTTP request: /api/users
20:16:28.088 app:server:http HTTP request completed: /api/users
20:16:28.088 app:server:http Received HTTP request: /api/products
20:16:28.139 app:server:http HTTP request completed: /api/products
20:16:28.139 app:server:websocket WebSocket message received: user-connected
20:16:28.170 app:server:websocket WebSocket message processed: user-connected

--- Changing debug configuration at runtime ---
Changing DEBUG from '*' to '*,!app:server:websocket,app:database'
20:16:28.170 app:server:http Received HTTP request: /api/settings
20:16:28.221 app:server:http HTTP request completed: /api/settings
20:16:28.252 app:database Executing complex query
20:16:28.252 app:database Query completed in 25ms
20:16:28.252 app:security Security audit completed
```

## Development

### Running Tests

The repository includes a script to run tests and generate code coverage reports:

```bash
# Run tests and generate coverage report in coverage.txt
./scripts/run_tests.sh
```

### Example Runner

The repository includes a script to run examples and capture their output for documentation:

```bash
# Run examples with different DEBUG settings and capture output
./scripts/run_examples.sh
```

This script:
- Creates a standalone environment for examples 
- Runs examples with different DEBUG settings
- Captures standard output and error streams
- Saves output files for documentation purposes
- Provides comprehensive cleanup steps

To clean up after running the script:

```bash
# Remove generated output files
rm -rf ./example_outputs
```

## License

[MIT](LICENSE) 