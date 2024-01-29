// Package static provides middleware for the Chi router to serve static files.
// This package allows detailed configuration including the static file prefix and the root
// directory for static files. It supports both real file systems and embedded file systems.
package static

import (
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// Logger interface defines the logging mechanism. It can be implemented by any logging library
// that provides a Printf method. This flexibility allows users to integrate their preferred logging
// solution.
type Logger interface {
	Printf(format string, v ...interface{})
}

// Config struct defines the configuration for the static file serving middleware.
// Fs refers to the file system (which can be embedded) containing the static files.
// Root specifies the root directory within the file system for the static files.
// FilePrefix is the URL prefix used to serve static files.
// CacheDuration is the duration for which the static files are cached.
//
// The Debug flag enables additional logging for troubleshooting, and Logger is an interface
// for a custom logging mechanism. If Logger is nil and Debug is true, a default logger is used.
type Config struct {
	Fs            fs.FS
	Root          string
	FilePrefix    string
	CacheDuration time.Duration

	Debug  bool
	Logger Logger
}

// middleware struct holds the configuration for a middleware instance.
type middleware struct {
	config Config
}

func Handler(config Config) func(next http.Handler) http.Handler {
	if config.Debug && config.Logger == nil {
		config.Logger = log.New(os.Stdout, "DEBUG: ", log.LstdFlags)
	}
	c := &middleware{config: config}

	return c.handler
}

// handler sets up the HTTP middleware handler. It serves static files based on the URL path
// matching the configured FilePrefix. If the path does not match, it passes the request
// to the next handler in the middleware chain.
func (m *middleware) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, m.config.FilePrefix) {
			if m.config.Debug {
				m.config.Logger.Printf("Serving static file: %s", r.URL.Path)
			}
			m.serveStaticFiles(w, r)
		} else {
			if m.config.Debug {
				m.config.Logger.Printf("Passing request to next handler: %s", r.URL.Path)
			}
			next.ServeHTTP(w, r)
		}
	})
}

// serveStaticFiles is responsible for serving the static files. It creates a sub-filesystem
// from the configured static root directory and serves the files using the standard library's
// file server.
func (m *middleware) serveStaticFiles(w http.ResponseWriter, r *http.Request) {
	staticFS, err := fs.Sub(m.config.Fs, m.config.Root)
	if err != nil {
		if m.config.Debug {
			m.config.Logger.Printf("Error creating sub-filesystem: %s", err)
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set Cache-Control header if CacheDuration is specified
	if m.config.CacheDuration > 0 {
		cacheControlValue := fmt.Sprintf("public, max-age=%d", int(m.config.CacheDuration.Seconds()))
		w.Header().Set("Cache-Control", cacheControlValue)
	}

	fileServer := http.FileServer(http.FS(staticFS))
	http.StripPrefix(m.config.FilePrefix, fileServer).ServeHTTP(w, r)
}
