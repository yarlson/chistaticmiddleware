package static

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"
)

// TestHandler tests the handling of requests by the middleware.
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
		Fs:            mockFS,
		Root:          "static",
		FilePrefix:    "/static",
		Debug:         true,
		CacheDuration: 365 * 24 * time.Hour,
	}

	r.Use(Handler(staticConfig))

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

// TestHandler404 tests the handling of requests by the middleware when the requested file does not exist.
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
		Fs:         mockFS,
		Root:       "static",
		FilePrefix: "/static",
	}

	r.Use(Handler(staticConfig))

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

// TestHandlerError tests the behavior of the Handler function when fs.Sub(m.config.Fs, m.config.Root) raises an error.
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
		Fs:         mockFS,
		Root:       "./static",
		FilePrefix: "/static",
		Debug:      true,
	}

	r.Use(Handler(staticConfig))

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
