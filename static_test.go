package chistaticmiddleware

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"
)

// TestNewStaticMiddleware tests the initialization of the StaticMiddleware.
func TestNewStaticMiddleware(t *testing.T) {
	tests := []struct {
		name      string
		config    Config
		wantDebug bool
	}{
		{
			name: "With Debug and No Logger",
			config: Config{
				Debug: true,
			},
			wantDebug: true,
		},
		{
			name: "Without Debug and No Logger",
			config: Config{
				Debug: false,
			},
			wantDebug: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStaticMiddleware(tt.config)
			if got.config.Debug != tt.wantDebug {
				t.Errorf("NewStaticMiddleware().Debug = %v, want %v", got.config.Debug, tt.wantDebug)
			}

			if tt.wantDebug && got.config.Logger == nil {
				t.Errorf("Expected logger to be set when Debug is true")
			}
		})
	}
}

// TestHandler tests the handling of requests by the StaticMiddleware.
func TestHandler(t *testing.T) {
	// Create a mock file system using fstest.MapFS
	mockFS := fstest.MapFS{
		"static/testfile.css": &fstest.MapFile{
			Data:    []byte("body {}"),
			ModTime: time.Now(),
		},
	}

	r := chi.NewRouter()
	staticConfig := Config{
		StaticFS:         mockFS,
		StaticRoot:       "static",
		StaticFilePrefix: "/static",
	}
	staticMiddleware := NewStaticMiddleware(staticConfig)

	r.Use(staticMiddleware.Handler())

	// Next handler for non-static routes
	nextHandlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
	})

	r.Get("/*", nextHandler)

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Test serving a static file
	res, err := http.Get(ts.URL + "/static/testfile.css")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}

	// Check MIME type
	expectedMimeType := "text/css; charset=utf-8"
	if contentType := res.Header.Get("Content-Type"); contentType != expectedMimeType {
		t.Errorf("Expected Content-Type '%s', got '%s'", expectedMimeType, contentType)
	}

	// Test passing to next handler
	nextHandlerCalled = false
	_, err = http.Get(ts.URL + "/nonexistent")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	if !nextHandlerCalled {
		t.Errorf("Expected next handler to be called for non-matching path")
	}
}

// TestHandler404 tests the handling of requests by the StaticMiddleware when the requested file does not exist.
func TestHandler404(t *testing.T) {
	// Create a mock file system using fstest.MapFS
	mockFS := fstest.MapFS{
		"static/testfile.css": &fstest.MapFile{
			Data:    []byte("body {}"),
			ModTime: time.Now(),
		},
	}

	r := chi.NewRouter()
	staticConfig := Config{
		StaticFS:         mockFS,
		StaticRoot:       "static",
		StaticFilePrefix: "/static",
	}
	staticMiddleware := NewStaticMiddleware(staticConfig)

	r.Use(staticMiddleware.Handler())

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Expected static handler to be called for matching path")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Test serving a static file
	res, err := http.Get(ts.URL + "/static/testfile.txt")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNotFound {
		t.Errorf("Expected status code 404, got %d", res.StatusCode)
	}
}

// TestHandlerError tests the behavior of the Handler function when fs.Sub(m.config.StaticFS, m.config.StaticRoot) raises an error.
func TestHandlerError(t *testing.T) {
	// Create a mock file system using fstest.MapFS
	mockFS := fstest.MapFS{
		"static/testfile.css": &fstest.MapFile{
			Data:    []byte("body {}"),
			ModTime: time.Now(),
		},
	}

	r := chi.NewRouter()
	staticConfig := Config{
		StaticFS:         mockFS,
		StaticRoot:       "./static",
		StaticFilePrefix: "/static",
		Debug:            true,
	}
	staticMiddleware := NewStaticMiddleware(staticConfig)

	r.Use(staticMiddleware.Handler())

	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("Expected static handler to be called for matching path")
	})

	ts := httptest.NewServer(r)
	defer ts.Close()

	// Test serving a static file
	res, err := http.Get(ts.URL + "/static/testfile.css")
	if err != nil {
		t.Fatalf("Failed to make GET request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("Expected status code 200, got %d", res.StatusCode)
	}
}
