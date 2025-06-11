package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var logger = NewColoredLogger("", nil)

var ignoredPaths = []string{
	"/.well-known/",
	"/favicon.ico",
	"/robots.txt",
	"/sitemap.xml",
}

func main() {
	err := godotenv.Load()
	if err != nil {
		logger.Error("Error loading .env file: %v", err)
		// Don't exit in production, just log the error
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := &exactPathMux{
		mux: http.NewServeMux(),
	}

	router.mux.HandleFunc("/", handleIndex)
	router.mux.HandleFunc("/contact-request", handleContactRequest)
	router.mux.Handle("/public/", http.StripPrefix("/public/", serveStaticFiles()))

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if shouldIgnorePath(r.URL.Path) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		h, pattern := router.pathHandler(r)

		if pattern == "" || !router.isExactMatch(r.URL.Path, pattern) {
			logger.Warning("404 Not Found: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - Page Not Found"))
			return
		}

		h.ServeHTTP(w, r)
	})

	logger.Info("Server starting on :%s", port)
	logger.Error("%v", http.ListenAndServe(":"+port, handler))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	logger.Info("Serving index page to %s", r.RemoteAddr)
	// Use the embedded filesystem to serve index.html
	indexFile, err := publicFS.Open("public/html/index.html")
	if err != nil {
		logger.Error("Error opening index.html: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer indexFile.Close()
	http.ServeContent(w, r, "index.html", time.Now(), indexFile.(io.ReadSeeker))
}

func handleContactRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logger.Warning("Invalid method %s from %s", r.Method, r.RemoteAddr)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Company  string `json:"company"`
		Message  string `json:"message"`
		Industry string `json:"industry"`
		Interest string `json:"interest"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Error("Error reading request body from %s: %v", r.RemoteAddr, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &req)
	if err != nil {
		logger.Error("Error unmarshaling request from %s: %v", r.RemoteAddr, err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = SendContactNotificationEmail(req.Name, req.Email, req.Company, req.Message, req.Industry, req.Interest)
	if err != nil {
		logger.Error("Failed to send notification email: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

type exactPathMux struct {
	mux *http.ServeMux
}

func (m *exactPathMux) pathHandler(r *http.Request) (h http.Handler, pattern string) {
	return m.mux.Handler(r)
}

func (m *exactPathMux) isExactMatch(path, pattern string) bool {
	if pattern == "/public/" {
		return strings.HasPrefix(path, pattern)
	}
	return path == pattern
}

func shouldIgnorePath(path string) bool {
	for _, ignored := range ignoredPaths {
		if strings.HasPrefix(path, ignored) {
			return true
		}
	}
	return false
}
