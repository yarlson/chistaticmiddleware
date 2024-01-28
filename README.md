# Chi Static Middleware

Chi Static Middleware is a Go package designed to work with the Chi router for serving static files efficiently. It is versatile, supporting both physical and embedded file systems, making it suitable for a wide range of applications, including web servers, SPAs, and HTMX.

## Features

- Serve static files with configurable URL prefixes.
- Support for physical and embedded file systems.
- Debug logging capabilities.
- Customizable logging interface.
- Configurable cache duration for static files.

## Installation

To install Chi Static Middleware, use the following command:

```bash
go get github.com/yarlson/chistaticmiddleware
```

## Usage

### Basic Setup

First, import the package along with Chi:

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/yarlson/chistaticmiddleware"
)
```

### Configuring Cache Duration

Set the cache duration for your static files to control browser caching. This is particularly useful for optimizing load times and reducing server load.

```go
staticConfig := chistaticmiddleware.Config{
    // ... other config settings ...
    CacheDuration: 24 * time.Hour, // Cache static files for 24 hours
}
```

### Using Physical File System

To serve files from a physical file system, configure the middleware like so:

```go
package main

import (
    "github.com/go-chi/chi/v5"
    "github.com/yarlson/chistaticmiddleware"
    "os"
    "time"
)

func main() {
    r := chi.NewRouter()

    staticConfig := chistaticmiddleware.Config{
        StaticFS:         os.DirFS("path/to/static/files"),
        StaticRoot:       "", // use "" for the root
        StaticFilePrefix: "/static",
        CacheDuration:    24 * time.Hour, // Optional: Cache for 24 hours
    }

    staticMiddleware := chistaticmiddleware.NewStaticMiddleware(staticConfig)
    r.Use(staticMiddleware.Handler())

    // setup other routes and start server...
}
```

### Using Embedded File System

If you're using Go 1.16 or later, serve static files from an embedded file system:

```go
package main

import (
    "embed"
    "github.com/go-chi/chi/v5"
    "github.com/yarlson/chistaticmiddleware"
    "time"
)

//go:embed path/to/static/files/*
var staticFiles embed.FS

func main() {
    r := chi.NewRouter()

    staticConfig := chistaticmiddleware.Config{
        StaticFS:         staticFiles,
        StaticRoot:       "path/to/static/files",
        StaticFilePrefix: "/static",
        CacheDuration:    24 * time.Hour, // Optional: Cache for 24 hours
    }

    staticMiddleware := chistaticmiddleware.NewStaticMiddleware(staticConfig)
    r.Use(staticMiddleware.Handler())

    // setup other routes and start server...
}
```

## Debugging

Enable debugging by setting the `Debug` flag in the configuration:

```go
staticConfig := chistaticmiddleware.Config{
    // ... other config
    Debug: true,
}
```

## Custom Logging

Implement the `Logger` interface to integrate custom logging:

```go
type CustomLogger struct {
    // implementation of the Logger interface
}

staticConfig := chistaticmiddleware.Config{
    // ... other config
    Logger: &CustomLogger{},
}
```

## License

This project is licensed under the MIT License. See the [LICENSE.md](LICENSE.md) file for details.
