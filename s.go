package main

import (
	"embed"
	"github.com/go-chi/chi/v5"
	"github.com/yarlson/chistaticmiddleware/static"
	"time"
)

//go:embed path/to/static/files/*
var staticFiles embed.FS

func main() {
	r := chi.NewRouter()

	staticConfig := static.Config{
		Fs:            staticFiles,
		Root:          "path/to/static/files",
		FilePrefix:    "/static",
		CacheDuration: 24 * time.Hour, // Optional: Cache for 24 hours
	}

	r.Use(static.Handler(staticConfig))

	// setup other routes and start server...
}
