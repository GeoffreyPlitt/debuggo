# debuggo

[![Go Report Card](https://goreportcard.com/badge/github.com/GeoffreyPlitt/debuggo)](https://goreportcard.com/report/github.com/GeoffreyPlitt/debuggo)
[![GoDoc](https://godoc.org/github.com/GeoffreyPlitt/debuggo?status.svg)](https://godoc.org/github.com/GeoffreyPlitt/debuggo)
[![Build Status](https://github.com/GeoffreyPlitt/debuggo/workflows/Go/badge.svg)](https://github.com/GeoffreyPlitt/debuggo/actions)
[![codecov](https://codecov.io/gh/GeoffreyPlitt/debuggo/branch/main/graph/badge.svg)](https://codecov.io/gh/GeoffreyPlitt/debuggo)

A lightweight debugging utility for Go applications inspired by the Node.js [debug](https://www.npmjs.com/package/debug) package.

## Features

- Environment variable controlled debugging
- Namespace support with hierarchical modules (e.g., `app:http:server`)
- Minimal overhead when disabled
- Runtime reconfiguration
- Selective module enabling/disabling

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

// Create debug loggers for different modules
var (
    debugApp = debuggo.Debug("app")
    debugDb = debuggo.Debug("db")
)

func main() {
    // Use debug loggers throughout your code
    debugApp("Application starting")
    
    // Debug messages only appear when the corresponding namespace is enabled
    debugDb("Connected to database")
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
DEBUG=api:* go run main.go  # Enables all 'api:' modules

# Enable everything except specific modules
DEBUG=*,!verbose go run main.go
```

## Advanced Usage

### Check if Debugging is Enabled

```go
if debuggo.IsEnabled("expensive:computation") {
    // Only do expensive debug operations if enabled
    debugComputation("Detailed debug info: %v", computeExpensiveDebugData())
}
```

### Runtime Reconfiguration

```go
// Update debug settings without restarting your application
os.Setenv("DEBUG", "newmodule,*:important")
debuggo.ReloadDebugSettings()
```

### Using PrefixWriter for Command Output

```go
cmd := exec.Command("some-program")

// Prefix stdout and stderr with custom identifiers
cmd.Stdout = &debuggo.PrefixWriter{Prefix: "CMD-STDOUT"}
cmd.Stderr = &debuggo.PrefixWriter{
    Prefix: "CMD-STDERR",
    Ignores: []string{"known warning to ignore"},
}
```

## Examples

### Basic Example

See [examples/basic/main.go](examples/basic/main.go) for a simple usage example.

### Advanced Example

See [examples/advanced/main.go](examples/advanced/main.go) for an example with hierarchical namespaces and runtime reconfiguration.

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