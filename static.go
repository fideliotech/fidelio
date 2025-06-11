//go:build !dev
// +build !dev

package main

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed public
var publicFS embed.FS

// getPublicFS returns a http.FileSystem that serves the embedded public directory
func getPublicFS() http.FileSystem {
	// Create a sub-filesystem for the public directory
	subFS, err := fs.Sub(publicFS, "public")
	if err != nil {
		panic(err)
	}
	return http.FS(subFS)
}

// serveStaticFiles returns a handler that serves static files from the embedded filesystem
func serveStaticFiles() http.Handler {
	return http.FileServer(getPublicFS())
}
