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

### Basic Example

See [examples/basic/main.go](examples/basic/main.go) for a simple usage example.

Run the basic example with:

```bash
# Show all debug messages
DEBUG=* go run examples/basic/main.go

# Show only database-related messages
DEBUG=db go run examples/basic/main.go
```

### Advanced Example

See [examples/advanced/main.go](examples/advanced/main.go) for:
- Hierarchical namespaces
- Runtime reconfiguration
- Selective enabling/disabling

Run the advanced example with:

```bash
# Show all debug messages
DEBUG=* go run examples/advanced/main.go

# Show only server-related messages
DEBUG=app:server:* go run examples/advanced/main.go

# Show all except websocket messages
DEBUG=*,!app:server:websocket go run examples/advanced/main.go
```

## License

[MIT](LICENSE) 